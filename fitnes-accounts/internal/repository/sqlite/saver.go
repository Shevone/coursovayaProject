package sqlite

import (
	"context"
	"fitnes-account/internal/models"
	"fmt"
)

const (
	roleAdmin = "Admin"
)

// SaveUser записывает нового пользователя в бд
func (s *Storage) SaveUser(ctx context.Context, user models.User) (int64, error) {
	const op = "storage.sqlite.SaveUser"

	// Простой запрос на добавление пользователя
	stmt, err := s.db.Prepare(
		"INSERT INTO users (email,pass_hash,name,surname,patronymic,phone_number,role) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// Выполняем запрос, передав параметры
	res, err := stmt.Exec(user.Email, user.PassHash, user.Name, user.Surname, user.Patronymic, user.PhoneNumber, user.Role)
	if err != nil {
		/*var sqliteErr sqlite3.Error

		// Небольшое кунг-фу для выявления ошибки ErrConstraintUnique
		// (см. подробности ниже)
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", op, repository.ErrUserExists)
		}
		*/
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем ID созданной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// EditUser редактирует данные пользователя в бд
func (s *Storage) EditUser(ctx context.Context, editedProfile models.User) error {
	const op = "storage.sqlite.SaveUser"

	// Простой запрос на добавление пользователя
	stmt, err := s.db.Prepare(
		"UPDATE users SET name = ?, surname = ?, patronymic =?  WHERE id = ?")
	if err != nil {

		return fmt.Errorf("%s: %w", op, err)
	}

	// Выполняем запрос, передав параметры
	_, err = stmt.ExecContext(ctx, editedProfile.Name, editedProfile.Surname, editedProfile.Patronymic, editedProfile.ID)
	if err != nil {
		/*var sqliteErr sqlite3.Error

		// Небольшое кунг-фу для выявления ошибки ErrConstraintUnique
		// (см. подробности ниже)
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return fmt.Errorf("%s: %w", op, repository.ErrUserExists)
		}
		*/
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) EditUserRole(ctx context.Context, userId int64, newRole string) (string, error) {
	const op = "storage.sqlite.EditUserRole"

	// Для начала смотрим, какая роль у пользователя сейчас

	stmt, err := s.db.Prepare("SELECT role FROM users WHERE id = ?")

	if err != nil {

		return "err", fmt.Errorf("%s: %w", op, err)
	}
	row := stmt.QueryRowContext(ctx, userId)
	var prevUserRole string
	err = row.Scan(&prevUserRole)
	if err != nil {
		return "err", fmt.Errorf("%s: %w", op, err)
	}
	if prevUserRole == newRole {
		return "Пользователь уже имеет эту роль", nil
	}
	if prevUserRole == roleAdmin {
		// TODO запрос из смежной таблицы(уровень админа - изменяемого, уровень админа - того, который изменяет

		// TODO если уровень того, который изменяет, выше уровня изменяемого то выдаем ошибку
	}
	// Если же роль пользователя отличается, то мы заменяем её

	stmt, err = s.db.Prepare("UPDATE users SET role = ? WHERE id = ?")

	if err != nil {

		return "err", fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.ExecContext(ctx, newRole, userId)
	if err != nil {
		return "err", fmt.Errorf("%s: %w", op, err)
	}
	return "Роль успешно обновлена", nil

}
