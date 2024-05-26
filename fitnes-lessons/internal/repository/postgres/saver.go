package postgres

import (
	"context"
	models "fitnes-lessons/internal/models"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"time"
)

// Saver
// ==============================================

func (s *Storage) EditLessonDb(ctx context.Context, lesson *models.Lesson) (int64, error) {
	const op = "pq.Edit"
	stmt, err := s.db.Prepare("UPDATE lessons SET name = $1,available_seats = $2, description = $3, difficulty = $4, starttime = $5, day_of_week = $6 WHERE id = $7")
	if err != nil {
		return -1, fmt.Errorf("%s, %w", op, err)
	}
	timeFormatTime := remakeStrToTimeType(lesson.Time)
	_, err = stmt.ExecContext(ctx, lesson.Title, lesson.AvailableSeats, lesson.Description, lesson.Difficult, timeFormatTime, lesson.DayOfWeek, lesson.LessonId)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return lesson.LessonId, nil
}

func (s *Storage) DeleteLessonDb(ctx context.Context, lessonId int64) (bool, error) {
	const op = "pq.DeleteLesson"
	stmt, err := s.db.Prepare(
		"DELETE FROM lessons WHERE id = $1")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, lessonId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	stmt, err = s.db.Prepare(
		"DELETE FROM student_lessons where lesson_id = $1")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.ExecContext(ctx, lessonId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return true, nil
}

func (s *Storage) CreateLessonDb(ctx context.Context, lesson *models.Lesson) (int64, error) {
	const op = "pq.CreateLesson"
	stmt, err := s.db.Prepare(
		"INSERT INTO lessons (id, name, trainer_id, available_seats,description,difficulty, startTime, day_of_week) values ($1,$2,$3,$4,$5,$6,$7,$8)")
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	timeFormatTime := remakeStrToTimeType(lesson.Time)
	_, err = stmt.ExecContext(ctx, lesson.LessonId, lesson.Title, lesson.TrainerId, lesson.AvailableSeats, lesson.Description, lesson.Difficult, timeFormatTime, lesson.DayOfWeek)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return lesson.LessonId, nil
}

func remakeStrToTimeType(timeStr string) time.Time {
	timeT, err := time.Parse(timeStr, timeStr)
	if err != nil {
		slog.Error(err.Error())
	}
	return timeT
}
func (s *Storage) ClearLessonSub(ctx context.Context, lessonId int64) (bool, error) {
	const op = "pq.ClearLessonSub"
	query := "DELETE FROM student_lessons WHERE lesson_id = $1"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, lessonId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return true, nil
}
