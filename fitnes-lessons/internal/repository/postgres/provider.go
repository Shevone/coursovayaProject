package postgres

import (
	"context"
	models "fitnes-lessons/internal/models"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	elOnPageCount = 10
)

// Provider
// ==============================================

// GetAllLessonsFromDb Получаем Общий список занятий
func (s *Storage) GetAllLessonsFromDb(ctx context.Context) ([]*models.Lesson, error) {
	const op = "postgresql.GetAllLessonsFromDb"
	queryRow := "SELECT id, name,trainer_id,available_seats,description,difficulty,startTime,day_of_week FROM lessons"
	stmt, err := s.db.Prepare(queryRow)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	modelsList := make([]*models.Lesson, 0)
	for rows.Next() {
		lesson := &models.Lesson{}
		err = rows.Scan(
			&lesson.LessonId,
			&lesson.Title,
			&lesson.TrainerId,
			&lesson.AvailableSeats,
			&lesson.Description,
			&lesson.Difficult,
			&lesson.Time,
			&lesson.DayOfWeek,
		)
		lesson.CurUsersCount = lesson.AvailableSeats
		if err != nil {
			return nil, fmt.Errorf(op+": %w", err)
		}
		modelsList = append(modelsList, lesson)
	}
	return modelsList, nil
}

// GetAllUserLessons Получаем Общий список подписок
func (s *Storage) GetAllUserLessons(ctx context.Context) (map[int64][]int64, error) {
	const op = "postgresql.GetAllUserLessons"
	// Сделать запрос на соседнюю таблицу
	query := "SELECT user_id, lesson_id FROM student_lessons"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	hashMap := make(map[int64][]int64)
	for rows.Next() {
		var userId int64
		var lessonId int64
		err = rows.Scan(&userId, &lessonId)
		if err != nil {
			return nil, fmt.Errorf(op+": %w", err)
		}
		_, ok := hashMap[userId]
		if !ok {
			hashMap[userId] = make([]int64, 0)
		}
		hashMap[userId] = append(hashMap[userId], lessonId)
	}
	return hashMap, nil
}
