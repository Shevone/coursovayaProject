package buffer

import (
	"fitnes-lessons/internal/models"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Buffer struct {
	mutex sync.Mutex
	// Общее расписание кешируется тут
	// ключ - номер дня недели
	schedule map[int32][]*models.Lesson

	// ===============================
	// Просто все существующие занятия, храним их для удобства изменения данных
	// ключ - id занятия
	lessons map[int64]*models.Lesson

	// ===============================
	// Map - занятия их тренера, которые их проводят
	// ключ - trainerId
	trainerLessons map[int64][]*models.Lesson

	// ===============================
	// Map - записи пользователей на занятие
	// ключ - trainerId
	userLessons map[int64][]*models.Lesson
}

// ======================================================================================================

func NewBuffer(lessons []*models.Lesson, userLessons map[int64][]int64) *Buffer {
	// На вход нам поступает
	// 1. Список занятий полный
	// 2. Map[userId] []int
	buffer := Buffer{
		schedule: map[int32][]*models.Lesson{
			0: make([]*models.Lesson, 0),
			1: make([]*models.Lesson, 0),
			2: make([]*models.Lesson, 0),
			3: make([]*models.Lesson, 0),
			4: make([]*models.Lesson, 0),
			5: make([]*models.Lesson, 0),
			6: make([]*models.Lesson, 0),
		},
		lessons:        make(map[int64]*models.Lesson),
		trainerLessons: make(map[int64][]*models.Lesson),
		userLessons:    make(map[int64][]*models.Lesson),
	}
	// Добавление в общие мапы
	buffer.AddLessons(lessons)
	// мапа для пользователей
	buffer.AddUserLessons(userLessons)

	return &buffer
}
func (b *Buffer) IsExist(targetLesson *models.Lesson) (int64, bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	// Проверям, существует
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for _, lesson := range b.lessons {
		// Сравнение по заданным полям
		if lesson.Title == targetLesson.Title &&
			lesson.Time == targetLesson.Time &&
			lesson.DayOfWeek == targetLesson.DayOfWeek &&
			lesson.TrainerId == targetLesson.TrainerId &&
			lesson.AvailableSeats == targetLesson.AvailableSeats &&
			lesson.Description == targetLesson.Description &&
			lesson.Difficult == targetLesson.Difficult {
			return lesson.LessonId, true
		}
	}
	return -1, false
}
func (b *Buffer) AddLessons(lessons []*models.Lesson) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for _, lesson := range lessons {
		lesson.Users = make([]int64, 0)
		b.lessons[lesson.LessonId] = lesson
		b.schedule[lesson.DayOfWeek] = append(b.schedule[lesson.DayOfWeek], lesson)
		b.trainerLessons[lesson.TrainerId] = append(b.trainerLessons[lesson.TrainerId], lesson)
	}
}
func (b *Buffer) AddUserLessons(userLessons map[int64][]int64) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for userId, curUserLessonsIds := range userLessons {
		b.userLessons[userId] = make([]*models.Lesson, 0)
		for _, lessonId := range curUserLessonsIds {
			b.userLessons[userId] = append(b.userLessons[userId], b.lessons[lessonId])
		}
	}
}

// ======================================================================================================

func (b *Buffer) GenerateLessonId() int64 {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	rand.Seed(time.Now().UnixNano())
	var id int64
	for {
		id = rand.Int63n(65000) // Генерация в диапазоне от 0 до максимального int64
		if _, ok := b.lessons[id]; !ok {
			break
		}
	}

	return id
}

// AddNewLesson Добавление новОго занятия
func (b *Buffer) AddNewLesson(lesson *models.Lesson) {
	// Добавить в расписание
	b.mutex.Lock()
	defer b.mutex.Unlock()
	lesson.Users = make([]int64, 0)
	b.schedule[lesson.DayOfWeek] = append(b.schedule[lesson.DayOfWeek], lesson)

	// Добавить в мапу id тренеров

	b.trainerLessons[lesson.TrainerId] = append(b.trainerLessons[lesson.TrainerId], lesson)

	// Добавить в map занятий по id

	b.lessons[lesson.LessonId] = lesson
}

// DeleteLesson AddLesson Удалить
func (b *Buffer) DeleteLesson(lessonId int64) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	lesson, ok := b.lessons[lessonId]
	if !ok {
		return false
	}
	// Удалить занятие из расписания
	b.schedule[lesson.DayOfWeek] = deleteFromSlice(b.schedule[lesson.DayOfWeek], lesson.LessonId)

	// Удалить из списка у тренера

	b.trainerLessons[lesson.TrainerId] = deleteFromSlice(b.trainerLessons[lesson.TrainerId], lesson.LessonId)

	// Удалить из общего списка
	delete(b.lessons, lesson.LessonId)
	for userId, userLessons := range b.userLessons {
		b.userLessons[userId] = deleteFromSlice(userLessons, lessonId)
	}
	return true
}

// Удаление из списка элемента с определенным id
func deleteFromSlice(slice []*models.Lesson, deleteId int64) []*models.Lesson {
	result := make([]*models.Lesson, 0, len(slice))
	for _, item := range slice {
		if item.LessonId != deleteId {
			result = append(result, item)
		}
	}
	return result
}

// EditLesson Edit Изменить
// Todo как будто бы не нужен
func (b *Buffer) EditLesson(lesson *models.Lesson) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	prevLesson := b.lessons[lesson.LessonId]
	lesson.Users = prevLesson.Users
	// Делаем подмену во всех мап-ах

	// В общем списке
	b.lessons[lesson.LessonId] = lesson
	// В расписании
	// получить, какой по счету и заменить
	if lesson.DayOfWeek == prevLesson.DayOfWeek && lesson.Time == prevLesson.Time {

		numberInShedule := b.getLessonIdFrom(b.schedule[lesson.DayOfWeek], lesson.LessonId)
		b.schedule[lesson.DayOfWeek][numberInShedule] = lesson

	} else {
		// Есл день недели изменился, то удаляем из прошлого дня и добавляем в новый
		b.schedule[prevLesson.DayOfWeek] = deleteFromSlice(b.schedule[prevLesson.DayOfWeek], lesson.LessonId)
		b.schedule[lesson.DayOfWeek] = append(b.schedule[lesson.DayOfWeek], lesson)
	}

	// В списке у тренеров
	// Получаем номер в списке у тренера
	numberInTrainerList := b.getLessonIdFrom(b.trainerLessons[lesson.TrainerId], lesson.LessonId)
	b.trainerLessons[lesson.TrainerId][numberInTrainerList] = lesson
	// В списке у пользователей
	// Тут мы проходимся по списку пользователей, у которые записаны в на занятие
	for _, userId := range lesson.Users {
		// у пользователя, который записан на занятие мы должны поменять ссылку
		numberInUserList := b.getLessonIdFrom(b.userLessons[userId], lesson.LessonId)
		b.userLessons[userId][numberInUserList] = lesson
	}
}
func (b *Buffer) getLessonIdFrom(lessonList []*models.Lesson, lessonId int64) int64 {
	for i, lesson := range lessonList {
		if lesson.LessonId == lessonId {
			return int64(i)
		}
	}
	return -1
}

func (b *Buffer) GetLessonById(lessonId int64) *models.Lesson {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.lessons[lessonId]
}

// =========================================================

// GetByTrainer Получить по тренеру
func (b *Buffer) GetByTrainer(trainerId int64, limit, offset int32) ([]*models.Lesson, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	trainerLessons := b.trainerLessons[trainerId]
	result := getElementsFromEnd(trainerLessons, int(limit), int(offset))
	return result, nil
}

// GetByUser Получить по пользователю
func (b *Buffer) GetByUser(userId int64, limit, offset int32) ([]*models.Lesson, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	userLessons := b.userLessons[userId]
	if userLessons == nil {
		return []*models.Lesson{}, nil
	}
	result := getElementsFromEnd(userLessons, int(limit), int(offset))
	return result, nil
}
func getElementsFromEnd(slice []*models.Lesson, n int, offset int) []*models.Lesson {
	// Проверяем, что offset не выходит за границы слайса

	startEl := len(slice) - (n * offset)
	if startEl < 0 {
		return nil
	}
	endEl := startEl - n
	if endEl < 0 {
		endEl = 0
	}

	return slice[endEl:startEl]
}

// =========================================================

// GetByWeekDay Получить по дню недели
func (b *Buffer) GetByWeekDay(weekDay int32) ([]*models.Lesson, error) {
	b.mutex.Lock()
	b.mutex.Unlock()
	return b.schedule[weekDay], nil
}

// IsUserSigned Записан ли пользователь на занятие
func (b *Buffer) IsUserSigned(userId int64, lessonId int64) bool {
	// Необходимо проверить если ли занятие с lessonId в списке у пользователя
	b.mutex.Lock()
	defer b.mutex.Unlock()
	userLessons := b.userLessons[userId]
	for _, lesson := range userLessons {
		if lesson.LessonId == lessonId {
			return true
		}
	}
	return false
}

// SignOnLesson Записаться на занятие
func (b *Buffer) SignOnLesson(userId int64, lessonId int64) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	lesson := b.lessons[lessonId]
	if b.userLessons[userId] == nil {
		b.userLessons[userId] = make([]*models.Lesson, 0)
	}
	lesson.Users = append(lesson.Users, userId)
	b.userLessons[userId] = append(b.userLessons[userId], lesson)
	if lesson.CurUsersCount == 0 {
		return fmt.Errorf("все места уже заняты")
	}
	lesson.CurUsersCount -= 1
	return nil
}

// SignOutLesson Отписаться
func (b *Buffer) SignOutLesson(userId int64, lessonId int64) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	// Удалем из общего
	lesson := b.lessons[lessonId]
	lesson.Users = delteteFromSlice(lesson.Users, userId)
	// Удаляем из мапы
	b.userLessons[userId] = deleteFromSlice(b.userLessons[userId], lesson.LessonId)
	if lesson.CurUsersCount == lesson.AvailableSeats {
		return fmt.Errorf("какая то ерунда")
	}
	lesson.CurUsersCount += 1
	return nil
}
func delteteFromSlice(slice []int64, deleteId int64) []int64 {
	result := make([]int64, 0, len(slice))
	for _, item := range slice {
		if item != deleteId {
			result = append(result, item)
		}
	}
	return result
}

// ClearSubs Очистить подписки по времени.
func (b *Buffer) ClearSubs(weekDay int32, inputTime string) []int64 {
	// необходимо удалить все подписки на занятие этого дня
	// в диапазоне врмение от текущего - полчаса
	b.mutex.Lock()
	defer b.mutex.Unlock()
	inputHour, inputMinute, err := parseTime(inputTime)
	if err != nil {
		return nil // Обработка ошибки
	}
	// Создаем время на 30 минут раньше
	startTime := time.Date(0, 0, 0, inputHour, inputMinute-30, 0, 0, time.UTC)

	// Создаем время окончания
	endTime := time.Date(0, 0, 0, inputHour, inputMinute, 0, 0, time.UTC)

	resultArray := make([]int64, 0)
	lessonsOnDay := b.schedule[weekDay]
	for _, lesson := range lessonsOnDay {
		// Парсим время события
		eventHour, eventMinute, err := parseTime(lesson.Time)
		if err != nil {
			continue // Пропускаем событие с ошибкой
		}
		// Создаем время события
		eventTime := time.Date(0, 0, 0, eventHour, eventMinute, 0, 0, time.UTC)

		// Проверяем, находится ли время события в диапазоне
		if eventTime.After(startTime) && eventTime.Before(endTime) {
			// если время в диапазоне, то чистим
			// Удалям из буфера
			b.removeSubs(lesson)
			// И записываем id в результирующий список
			resultArray = append(resultArray, lesson.LessonId)
		}
	}
	return resultArray

}
func (b *Buffer) removeSubs(lesson *models.Lesson) {
	// Убираем у каждого пользователя
	for _, userId := range lesson.Users {
		b.userLessons[userId] = deleteFromSlice(b.userLessons[userId], lesson.LessonId)
	}
	// Очищаем список у самого занятия
	lesson.Users = make([]int64, 0)
	lesson.CurUsersCount = 0

}

func (b *Buffer) LessonFreeSeats(lessonId int64) int32 {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	lesson := b.lessons[lessonId]
	freSeat := lesson.CurUsersCount
	return freSeat
}

func parseTime(timeStr string) (int, int, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат времени: %s", timeStr)
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("неверный формат часа: %s", parts[0])
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("неверный формат минуты: %s", parts[1])
	}

	return hour, minute, nil
}
