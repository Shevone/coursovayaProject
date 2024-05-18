package sqlite

import (
	"context"
	models "fitnes-lessons/internal/models"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

// Saver
// ==============================================

func (s Storage) EditLessonDb(ctx context.Context, lesson *models.Lesson) (int64, error) {
	const op = "SQLite.Edit"
	stmt, err := s.db.Prepare("UPDATE lessons SET title = ?,available_seats = ?, description = ?, difficult = ?, date_and_time = ? WHERE id = ?")
	if err != nil {
		return -1, fmt.Errorf("%s, %w", op, err)
	}
	timeAndDate := remakeStrToTimeType(lesson.Time, lesson.Date)
	_, err = stmt.ExecContext(ctx, lesson.Title, lesson.AvailableSeats, lesson.Description, lesson.Difficult, timeAndDate, lesson.LessonId)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return lesson.LessonId, nil
}

func (s Storage) DeleteLessonDb(ctx context.Context, lessonId int64) (bool, error) {
	const op = "SQLite.DeleteLesson"
	stmt, err := s.db.Prepare(
		"DELETE FROM lessons WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, lessonId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	stmt, err = s.db.Prepare(
		"DELETE FROM student_lessons where lesson_id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.ExecContext(ctx, lessonId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return true, nil
}

func (s Storage) CreateLessonDb(ctx context.Context, lesson *models.Lesson) (int64, error) {
	const op = "SQLite.CreateLesson"
	stmt, err := s.db.Prepare(
		"INSERT INTO main.lessons (title, date_and_time ,trainer_id, available_seats,description,difficult) values (?,?,?,?,?,?)")
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	timeAndDate := remakeStrToTimeType(lesson.Time, lesson.Date)
	res, err := stmt.ExecContext(ctx, lesson.Title, timeAndDate, lesson.TrainerId, lesson.AvailableSeats, lesson.Description, lesson.Difficult)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func remakeStrToTimeType(timeStr string, dateStr string) time.Time {
	timeLayout := "15:04 02.01" // Формат времени: часы:минуты:секунды день.месяц
	timeAndDateStr := timeStr + " " + dateStr
	timeAndDate, _ := time.Parse(timeLayout, timeAndDateStr)
	return timeAndDate
}
