package postgres

import (
	"context"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

// Subscriber
// ==============================================

func (s *Storage) SignForLesson(ctx context.Context, lessonId int64, userId int64) error {
	const op = "postgres.SubFromLesson"
	stmt, err := s.db.Prepare("INSERT INTO student_lessons (user_id, lesson_id) values ($1,$2)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.ExecContext(ctx, userId, lessonId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) SignOutFromLesson(ctx context.Context, lessonId int64, userId int64) error {
	const op = "postgres.UnSubFromLesson"
	stmt, err := s.db.Prepare("DELETE FROM student_lessons where lesson_id = $1 and user_id = $2")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.ExecContext(ctx, lessonId, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
