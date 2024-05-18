package sqlite

import (
	"context"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

// Subscriber
// ==============================================

func (s Storage) CloseLessonDb(ctx context.Context, lessonId int64) (bool, error) {
	const op = "SQLite.CloseLesson"
	stmt, err := s.db.Prepare("UPDATE lessons SET is_complete = false WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s : %w", op, err)
	}
	_, err = stmt.ExecContext(ctx, lessonId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return true, nil

}
func (s Storage) SignUpForLessonOrCancelDb(ctx context.Context, lessonId int64, userId int64) (string, error) {
	const op = "SQLite.SignForLessonOrCancel"
	// Получить is_complete у данного урока
	queryRow := "SELECT lessons.is_complete, lessons.available_seats from lessons WHERE id = ?"
	var isComplete bool
	var avalieableSeats int
	err := s.db.QueryRow(queryRow, lessonId).Scan(&isComplete, &avalieableSeats)
	if err != nil {
		return "", fmt.Errorf("%s:%w", op, err)
	}

	// Смотрим подписан ли пользователь
	queryRow = "SELECT count(*) from student_lessons where user_id = ? AND lesson_id = ?"
	var count int
	err = s.db.QueryRow(queryRow, userId, lessonId).Scan(&count)
	if err != nil {
		return "", fmt.Errorf("%s:%w", op, err)
	}

	// Если подписан, то пытаемся отписать
	if count > 0 {
		return s.unSubFromLesson(lessonId, userId, avalieableSeats)
	}
	// Если не подписаны
	// Смотрим, если закрыт
	if isComplete {
		// если закрыта то взовращаем осообзение
		return "Запись на заняите уже закрыта", nil
	}
	// Подписываеммся
	return s.subForLesson(lessonId, userId, avalieableSeats)
}
func (s Storage) subForLesson(lessonId int64, userId int64, count int) (string, error) {
	const op = "SQLite.SubFromLesson"
	stmt, err := s.db.Prepare("INSERT INTO student_lessons (user_id, lesson_id) values (?,?)")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(userId, lessonId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	stmt, err = s.db.Prepare("UPDATE lessons SET available_seats = ? WHERE lessons.id = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(count-1, lessonId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return "Вы подписаны на это занятия", nil
}
func (s Storage) unSubFromLesson(lessonId int64, userId int64, count int) (string, error) {
	const op = "SQLite.UnSubFromLesson"
	stmt, err := s.db.Prepare("DELETE FROM student_lessons where lesson_id = ? and user_id = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(lessonId, userId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	stmt, err = s.db.Prepare("UPDATE lessons SET available_seats = ? WHERE lessons.id = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(count+1, lessonId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return "Вы отписаны от этого занятия", nil
}
