package sqlite

import (
	"context"
	models "fitnes-lessons/internal/models"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"time"
)

var (
	elOnPageCount = 10
)

// Provider
// ==============================================

func (s Storage) GetAllLessonsFromDb(ctx context.Context, page int32, limit int32) ([]models.Lesson, error) {
	const op = "SQLite.GetAllLessons"
	//queryRow := "SELECT id, title,date_and_time,trainer_id,available_seats,description,difficult FROM lessons WHERE strftime('%m-%d %H:%M', date_and_time) >= strftime('%m-%d %H:%M', 'now') ORDER BY date_and_time LIMIT ?, ?"
	queryRow := "SELECT id, title,date_and_time,trainer_id,available_seats,description,difficult FROM lessons ORDER BY date_and_time LIMIT ?, ?"
	rows, err := s.db.Query(queryRow, page, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	list := make([]models.Lesson, 0)
	for rows.Next() {
		var timeAndDate string
		lesson := models.Lesson{}
		err = rows.Scan(&lesson.LessonId,
			&lesson.Title,
			&timeAndDate,
			&lesson.TrainerId,
			&lesson.AvailableSeats,
			&lesson.Description,
			&lesson.Difficult)
		timeStr, dateStr := returnStringTimeAndDate(timeAndDate)
		lesson.Time = timeStr
		lesson.Date = dateStr
		if err != nil {
			return nil, fmt.Errorf("%s, %w", op, err)
		}
		list = append(list, lesson)
	}
	return list, nil
}
func (s Storage) GetLessonsByTrainerIdFromDb(ctx context.Context, trainerId int64, page int32, limit int32) ([]models.Lesson, error) {
	const op = "SQLite.GetAllLessons"
	queryRow := "SELECT id, title,date_and_time,trainer_id,available_seats,description,difficult FROM lessons WHERE strftime('%m-%d %H:%M', date_and_time) >= strftime('%m-%d %H:%M', 'now') AND trainer_id = ? ORDER BY date_and_time LIMIT ?, ?"
	rows, err := s.db.Query(queryRow, trainerId, page, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	list := make([]models.Lesson, 0)
	for rows.Next() {
		var timeAndDate string
		lesson := models.Lesson{}
		err = rows.Scan(&lesson.LessonId,
			&lesson.Title,
			&timeAndDate,
			&lesson.TrainerId,
			&lesson.AvailableSeats,
			&lesson.Description,
			&lesson.Difficult)
		timeStr, dateStr := returnStringTimeAndDate(timeAndDate)
		lesson.Time = timeStr
		lesson.Date = dateStr
		if err != nil {
			return nil, fmt.Errorf("%s, %w", op, err)
		}
		list = append(list, lesson)
	}
	return list, nil
}

func (s Storage) GetLessonsByUserIdFromDb(ctx context.Context, userId int64, page int32, limit int32) ([]models.Lesson, error) {
	const op = "SQLite.GetAllLessons"
	queryRowGetUserLessons := "SELECT lesson_id FROM student_lessons where user_id = ?"
	lessonIdList := make([]int64, 0)
	rows, err := s.db.Query(queryRowGetUserLessons, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	for rows.Next() {
		var lessonId int64
		err = rows.Scan(&lessonId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		lessonIdList = append(lessonIdList, lessonId)
	}

	rows.Close()
	list := make([]models.Lesson, 0)
	for _, lessonId := range lessonIdList {
		queryRow := "SELECT id, title,date_and_time,trainer_id,available_seats,description,difficult FROM lessons WHERE id = ? ORDER BY date_and_time DESC LIMIT ?,?"
		rows, err := s.db.Query(queryRow, lessonId, page, limit)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		defer rows.Close()

		for rows.Next() {
			var timeAndDate string
			lesson := models.Lesson{}
			err = rows.Scan(&lesson.LessonId,
				&lesson.Title,
				&timeAndDate,
				&lesson.TrainerId,
				&lesson.AvailableSeats,
				&lesson.Description,
				&lesson.Difficult)
			timeStr, dateStr := returnStringTimeAndDate(timeAndDate)
			lesson.Time = timeStr
			lesson.Date = dateStr
			if err != nil {
				return nil, fmt.Errorf("%s, %w", op, err)
			}
			list = append(list, lesson)
		}
	}
	return list, nil
}

func (s Storage) GetLessonFromDb(ctx context.Context, lessonId int64) (*models.Lesson, error) {
	const op = "SQLite.GetLesson"
	queryRow := "SELECT id, title,date_and_time,trainer_id,available_seats,description,difficult FROM lessons WHERE id = ?"
	lesson := &models.Lesson{}
	var timeAndDate string
	err := s.db.QueryRow(queryRow, lessonId).
		Scan(&lesson.LessonId,
			&lesson.Title,
			&timeAndDate,
			&lesson.TrainerId,
			&lesson.AvailableSeats,
			&lesson.Description,
			&lesson.Difficult)
	timeStr, dateStr := returnStringTimeAndDate(timeAndDate)
	lesson.Time = timeStr
	lesson.Date = dateStr
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return lesson, nil
}
func returnStringTimeAndDate(dateAndTime string) (string, string) {
	// Преобразование строки времени и даты в тип time.Time
	timeAndDate, _ := time.Parse("2006-01-02 15:04:05-07:00", dateAndTime)
	// Получение времени в формате "hh:mm"
	timeStr := timeAndDate.Format("15:04")

	// Получение даты в формате "день.месяц"
	dateStr := timeAndDate.Format("02.01")

	return timeStr, dateStr
}
