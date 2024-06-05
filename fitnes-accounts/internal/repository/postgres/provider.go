package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fitnes-account/internal/models"
	"fitnes-account/internal/repository"
	"fmt"
)

// User получаем пользователя из бд по email-у
func (s *Storage) User(ctx context.Context, email string) (*models.User, error) {
	const op = "storage.pq.User"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash, name, surname, patronymic, phone_number, role FROM users WHERE email = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	row := stmt.QueryRowContext(ctx, email)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash, &user.Name, &user.Surname, &user.Patronymic, &user.PhoneNumber, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrAppNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (s *Storage) GetUserDataById(ctx context.Context, userid int64) (*models.User, error) {
	const op = "storage.pq.UserById"

	stmt, err := s.db.Prepare(
		"SELECT email, name,surname,patronymic,phone_number, role FROM users WHERE id = $1",
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	row := stmt.QueryRowContext(ctx, userid)

	var user models.User
	err = row.Scan(&user.Email, &user.Name, &user.Surname, &user.Patronymic, &user.PhoneNumber, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return nil, fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

// GetUsers получаем из бд список юзеров для админа
func (s *Storage) GetUsers(ctx context.Context, page int64, limit int64) ([]*models.User, error) {
	const op = "storage.pq.GetUsers"
	s.usersMutex.Lock()
	defer s.usersMutex.Unlock()
	result := make([]*models.User, 0)
	startPosition := int(page * limit)
	endPosition := startPosition + int(limit)
	// Если у нас элементов меньше чем стартовая позиция
	if startPosition > len(s.userIds) {
		return result, nil
	}
	// Если у нас конечная позоция больше чем длина
	if endPosition > len(s.userIds) {
		endPosition = len(s.userIds)
	}
	// Копируем не все(включаем стартовый элемент и не включаем последний)
	for i := startPosition; i < endPosition; i++ {
		userId := s.userIds[i]
		user := s.users[userId]
		result = append(
			result,
			user,
		)
	}
	return result, nil
}

func (s *Storage) GetUsersOld(ctx context.Context, page int64, limit int64) ([]models.User, error) {
	const op = "storage.pq.GetUsers"

	stmt, err := s.db.Prepare(
		"SELECT id, email, name, surname, patronymic,role, phone_number FROM users ORDER BY id offset $1 limit $2")
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

// GetUsers получаем из бд список юзеров для админа
func (s *Storage) getUsers(ctx context.Context) ([]models.User, error) {
	const op = "storage.pq.GetUsers"

	stmt, err := s.db.Prepare(
		"SELECT id, email, name, surname, patronymic,role, phone_number FROM users ORDER BY id")
	if err != nil {

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	row, err := stmt.QueryContext(ctx)
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
func (s *Storage) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	s.adminLevelsMutex.Lock()
	defer s.adminLevelsMutex.Unlock()

	_, ok := s.adminLevels[userId]
	return ok, nil
}
