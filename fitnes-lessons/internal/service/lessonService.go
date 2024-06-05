package service

import (
	"context"
	"fitnes-lessons/internal/models"
	"fitnes-lessons/internal/service/buffer"
	"fmt"
	"log/slog"
	"time"
)

type LessonService struct {
	lessonSaver      LessonSaver
	lessonProvider   LessonProvider
	lessonSubscriber LessonSubscriber
	buffer           *buffer.Buffer
}

type LessonSaver interface {
	CreateLessonDb(ctx context.Context, lesson *models.Lesson) (int64, error)

	EditLessonDb(ctx context.Context, lesson *models.Lesson) (int64, error)

	DeleteLessonDb(ctx context.Context, lessonId int64) (bool, error)
	ClearLessonSub(ctx context.Context, lessonId int64) (bool, error)
}
type LessonProvider interface {
	GetAllLessonsFromDb(ctx context.Context) ([]*models.Lesson, error)
	GetAllUserLessons(ctx context.Context) (map[int64][]int64, error)
}
type LessonSubscriber interface {
	SignForLesson(ctx context.Context, lessonId int64, userId int64) error
	SignOutFromLesson(ctx context.Context, lessonId int64, userId int64) error
}

func NewLessonService(ctx context.Context, lessonSaver LessonSaver, provider LessonProvider, subs LessonSubscriber) (*LessonService, error) {
	lessons, err := provider.GetAllLessonsFromDb(ctx)
	if err != nil {
		return nil, err
	}
	userLessons, err := provider.GetAllUserLessons(ctx)
	if err != nil {
		return nil, err
	}
	// Todo запуск воркера
	buff := buffer.NewBuffer(lessons, userLessons)
	service := &LessonService{
		lessonSaver:      lessonSaver,
		lessonProvider:   provider,
		lessonSubscriber: subs,
		buffer:           buff,
	}
	go service.runWorker(ctx)
	return service, nil
}

func (l *LessonService) runWorker(ctx context.Context) {
	// Работает таймер
	// Как только отрабатывает мы получаем сегодняшний день недели в числе
	// А так же получаем текщуее время
	// Отправляем время в буфер
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			return
			// Завершаем работу
		case <-ticker.C:
			// Получение текущего времени
			now := time.Now()
			// Получение дня недели в числовом виде (0 - воскресенье, 6 - суббота)
			weekday := int32(now.Weekday())
			formattedTime := now.Format("15:04")
			clearRes := l.buffer.ClearSubs(weekday, formattedTime)
			for _, lessonId := range clearRes {
				_, err := l.lessonSaver.ClearLessonSub(ctx, lessonId)
				if err != nil {
					slog.Error(err.Error())
				}
			}
		}
	}
}

// Походы в репозиторий совершаются только тогда, когда:
// 1 создаем занятие
// 2 Подписываем на занятие
// 3 Отписываем от занятия
// 4 Удаляем
// Изменяем ? не пользуемся
// Все остальное, берём из кеша

func (l *LessonService) CreateLesson(ctx context.Context, lesson *models.Lesson) (int64, error) {
	const op = "Lessons.CreatNew"
	// 0 Генерация id для занятия

	if id, ok := l.buffer.IsExist(lesson); ok {
		return id, nil
	}
	lesson.LessonId = l.buffer.GenerateLessonId()
	// 1 создание в бд
	lessonId, err := l.lessonSaver.CreateLessonDb(ctx, lesson)
	if err != nil {
		slog.Error("failed to create lesson", err)

		return -1, fmt.Errorf("%s: %w", op, err)
	}
	// 2 запись в буфер
	l.buffer.AddNewLesson(lesson)
	return lessonId, nil

}

func (l *LessonService) DeleteLesson(ctx context.Context, lessonId int64) (bool, error) {
	const op = "Lesson.Delete"
	// Удаляем из бд
	_, err := l.lessonSaver.DeleteLessonDb(ctx, lessonId)
	if err != nil {
		slog.Error("failed to delete lesson", err)

		return false, fmt.Errorf("%s: %w", op, err)
	}
	// Удаляем из буфера
	delteRes := l.buffer.DeleteLesson(lessonId)
	return delteRes, nil
}

func (l *LessonService) GetLesson(ctx context.Context, lessonId int64) (*models.Lesson, error) {
	const op = "Lesson.Get.ByLId"
	// Пытаемся получить из буфера
	lesson := l.buffer.GetLessonById(lessonId)
	if lesson == nil {
		slog.Error("failed to get lesson by id", lessonId)
	}
	return lesson, nil
}

func (l *LessonService) GetLessonsByTrainerId(ctx context.Context, trainerId int64, page int32, limit int32) ([]*models.Lesson, error) {
	const op = "Lesson.Get.ByTrainerId"
	// Смотрим в буфере
	lessons, _ := l.buffer.GetByTrainer(trainerId, limit, page)
	if lessons == nil {
		return nil, fmt.Errorf("%s: %s", op, "failed to get trainer lessons")
	}
	return lessons, nil
}

func (l *LessonService) GetLessonsByUserId(ctx context.Context, userId int64, page int32, limit int32) ([]*models.Lesson, error) {
	const op = "Lesson.Get.ByUserId"
	// Идем в буфер и смотрим там
	userLesson, _ := l.buffer.GetByUser(userId, limit, page)
	if userLesson == nil {
		return nil, fmt.Errorf("%s: %s", op, "failed to get user lessons")
	}
	return userLesson, nil
}

func (l *LessonService) GetLessonsByWeekDay(ctx context.Context, weekDay int32) ([]*models.Lesson, error) {
	const op = "Lesson.Get.All"

	res, _ := l.buffer.GetByWeekDay(weekDay)
	if res == nil {
		return nil, fmt.Errorf("%s: %s", op, "failed to get lessons")
	}
	return res, nil
}

func (l *LessonService) SignUpForLessonOrCancel(ctx context.Context, lessonId int64, userId int64) (string, error) {
	const op = "Lesson.SignUpOrCancel"
	// Если пользователь попытается записаться на занятие сегодняшнего дня, но которое уже прошло, то выдаем ошибку
	// 1. Смотрим через кеш, записан ли пользователь уже на это занятие
	isSub := l.buffer.IsUserSigned(userId, lessonId)
	var msg string
	if isSub {
		// Да - Метод отписки
		err := l.lessonSubscriber.SignOutFromLesson(ctx, lessonId, userId)
		if err != nil {
			slog.Error("failed to sign out lesson", err)
			return "Произошла ошибка при отписке с занятия, попробуйте позже", err
		}
		err = l.buffer.SignOutLesson(userId, lessonId)
		if err != nil {
			return err.Error(), err
		}
		msg = "Вы отписаны от занятия"
	} else {
		// Нет - вызываем методы подписки
		// Количество свободных мест == 0 => сообщение о невозомонжости операции
		// Количество свободных мест >= 0 => делаем операцию
		seatsCount := l.buffer.LessonFreeSeats(lessonId)
		if seatsCount == 0 {
			return "Запись невозможна, мест нет", nil
		}
		err := l.lessonSubscriber.SignForLesson(ctx, lessonId, userId)
		if err != nil {
			slog.Error("failed to sign in lesson", err)
			return "Произошла ошибка при записи на занятие, попробуйте позже", err
		}
		err = l.buffer.SignOnLesson(userId, lessonId)
		if err != nil {
			return err.Error(), err
		}
		msg = "Вы записаны на занятие"
	}
	return msg, nil

}

// ===============================
// Корзина

//	func (l *LessonService) CloseLesson(ctx context.Context, lessonId int64, trainerId int64) (bool, error) {
//		const op = "Lesson.CloseLesson"
//		// Метод не используется
//		res, err := l.lessonSubscriber.CloseLessonDb(ctx, lessonId)
//		if err != nil {
//			slog.Error("failed to close lesson", err)
//
//			return false, fmt.Errorf("%s: %w", op, err)
//		}
//		return res, nil
//	}
func (l *LessonService) EditLesson(ctx context.Context, lesson *models.Lesson) (int64, error) {

	const op = "Lesson.Edit"
	// Изменяем в бд
	lessonPrev := l.buffer.GetLessonById(lesson.LessonId)

	// Переменная содержит количесство свободных мест и если вдруг у нас количество мест изменилось, то количество свободных надо пересчитать
	if lessonPrev.AvailableSeats != lesson.AvailableSeats {
		userCount := lessonPrev.AvailableSeats - lessonPrev.CurUsersCount
		if userCount > lesson.AvailableSeats {
			return -1, fmt.Errorf("Нельзя установить такое количество мест")
		}
		lesson.CurUsersCount = lesson.AvailableSeats - userCount
	}

	res, err := l.lessonSaver.EditLessonDb(ctx, lesson)
	if err != nil {
		slog.Error("failed to edit lesson", err)

		return -1, fmt.Errorf("%s: %w", op, err)
	}
	// Если изменения прошли успешно, то меняем и в буфере
	l.buffer.EditLesson(lesson)
	return res, nil
}
