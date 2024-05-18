package service

import (
	"context"
	"fitnes-lessons/internal/models"
	"fmt"
	"log/slog"
)

type LessonService struct {
	logger           *slog.Logger
	lessonSaver      LessonSaver
	lessonProvider   LessonProvider
	lessonSubscriber LessonSubscriber
}

type LessonSaver interface {
	CreateLessonDb(ctx context.Context, lesson *models.Lesson) (int64, error)
	EditLessonDb(ctx context.Context, lesson *models.Lesson) (int64, error)
	DeleteLessonDb(ctx context.Context, lessonId int64) (bool, error)
}
type LessonProvider interface {
	GetLessonFromDb(ctx context.Context, lessonId int64) (*models.Lesson, error)
	GetLessonsByTrainerIdFromDb(ctx context.Context, trainerId int64, page int32, limit int32) ([]models.Lesson, error)
	GetLessonsByUserIdFromDb(ctx context.Context, userId int64, page int32, limit int32) ([]models.Lesson, error)
	GetAllLessonsFromDb(ctx context.Context, page int32, limit int32) ([]models.Lesson, error)
}
type LessonSubscriber interface {
	SignUpForLessonOrCancelDb(ctx context.Context, lessonId int64, userId int64) (string, error)
	CloseLessonDb(ctx context.Context, lessonId int64) (bool, error)
}

func NewLessonService(log *slog.Logger, lessonSaver LessonSaver, provider LessonProvider, subs LessonSubscriber) *LessonService {
	return &LessonService{
		logger:           log,
		lessonSaver:      lessonSaver,
		lessonProvider:   provider,
		lessonSubscriber: subs,
	}
}

func (l LessonService) CreateLesson(ctx context.Context, lesson *models.Lesson) (int64, error) {
	const op = "Lessons.CreatNew"
	lessonId, err := l.lessonSaver.CreateLessonDb(ctx, lesson)
	if err != nil {
		l.logger.Error("failed to create lesson", err)

		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return lessonId, nil

}

func (l LessonService) GetLesson(ctx context.Context, lessonId int64) (*models.Lesson, error) {
	const op = "Lesson.Get.ByLId"
	lesson, err := l.lessonProvider.GetLessonFromDb(ctx, lessonId)
	if err != nil {
		l.logger.Error("failed to get lesson", err)

		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return lesson, nil
}

func (l LessonService) DeleteLesson(ctx context.Context, lessonId int64) (bool, error) {
	const op = "Lesson.Delete"
	res, err := l.lessonSaver.DeleteLessonDb(ctx, lessonId)
	if err != nil {
		l.logger.Error("failed to delete lesson", err)

		return false, fmt.Errorf("%s: %w", op, err)
	}
	return res, nil
}

func (l LessonService) EditLesson(ctx context.Context, lesson *models.Lesson) (int64, error) {
	const op = "Lesson.Edit"
	res, err := l.lessonSaver.EditLessonDb(ctx, lesson)
	if err != nil {
		l.logger.Error("failed to edit lesson", err)

		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return res, nil
}

func (l LessonService) GetLessonsByTrainerId(ctx context.Context, trainerId int64, page int32, limit int32) ([]models.Lesson, error) {
	const op = "Lesson.Get.ByTrainerId"
	list, err := l.lessonProvider.GetLessonsByTrainerIdFromDb(ctx, trainerId, page, limit)
	if err != nil {
		l.logger.Error("failed to get lessons by trainer id", err)

		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return list, nil

}

func (l LessonService) GetLessonsByUserId(ctx context.Context, userId int64, page int32, limit int32) ([]models.Lesson, error) {
	const op = "Lesson.Get.ByUserId"
	list, err := l.lessonProvider.GetLessonsByUserIdFromDb(ctx, userId, page, limit)
	if err != nil {
		l.logger.Error("failed to get lessons by user id", err)

		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return list, nil
}

func (l LessonService) GetAllLessons(ctx context.Context, page int32, limit int32) ([]models.Lesson, error) {
	const op = "Lesson.Get.All"
	res, err := l.lessonProvider.GetAllLessonsFromDb(ctx, page, limit)
	if err != nil {
		l.logger.Error("failed to get lessons", err)

		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return res, nil
}

func (l LessonService) CloseLesson(ctx context.Context, lessonId int64, trainerId int64) (bool, error) {
	const op = "Lesson.CloseLesson"
	res, err := l.lessonSubscriber.CloseLessonDb(ctx, lessonId)
	if err != nil {
		l.logger.Error("failed to close lesson", err)

		return false, fmt.Errorf("%s: %w", op, err)
	}
	return res, nil
}

func (l LessonService) SignUpForLessonOrCancel(ctx context.Context, lessonId int64, userId int64) (string, error) {
	// TODO придумать что будет возвращать
	const op = "Lesson.SignUpOrCancel"
	res, err := l.lessonSubscriber.SignUpForLessonOrCancelDb(ctx, lessonId, userId)
	if err != nil {
		l.logger.Error("failed to delete lesson", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}
	return res, nil
}
