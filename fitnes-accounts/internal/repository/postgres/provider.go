package repository

import (
	"context"
	"database/sql"
	"errors"
	"fitnes-account/internal/models"
	"fmt"
)

// User получаем пользователя из бд по email-у
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.sqlite.User"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash, name, surname,patronymic,phone_number,role FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash, &user.Name, &user.Surname, &user.Patronymic, &user.PhoneNumber, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) GetUserDataById(ctx context.Context, userid int64) (models.User, error) {
	const op = "storage.sqlite.UserById"

	stmt, err := s.db.Prepare(
		"SELECT email, name,surname,patronymic,phone_number, role FROM users WHERE id = ?",
	)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	row := stmt.QueryRowContext(ctx, userid)

	var user models.User
	err = row.Scan(&user.Email, &user.Name, &user.Surname, &user.Patronymic, &user.PhoneNumber, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

// GetUsers получаем из бд список юзеров для админа
func (s *Storage) GetUsers(ctx context.Context, page int64, limit int64) ([]models.User, error) {
	const op = "storage.sqlite.GetUsers"

	stmt, err := s.db.Prepare(
		"SELECT id, email, name, surname, patronymic,role, phone_number FROM users ORDER BY id LIMIT ?,?")
	if err != nil {

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	row, err := stmt.QueryContext(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	usersList := make([]models.User, 0)
	for row.Next() {
		user := models.User{}
		err = row.Scan(&user.ID,
			&user.Email,
			&user.Name,
			&user.Surname,
			&user.Patronymic,
			&user.Role,
			&user.PhoneNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		usersList = append(usersList, user)
	}
	return usersList, nil
}
