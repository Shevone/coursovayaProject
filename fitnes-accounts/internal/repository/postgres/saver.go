package postgres

import (
	"context"
	"fitnes-account/internal/models"
	"fmt"
)

const (
	roleAdmin   = "Admin"
	userRoleCtx = "userRole"
	userIdCtx   = "userId"
)

func (s *Storage) UpdateUserPassword(ctx context.Context, userId int64, newPassword []byte, updaterId int64) (string, error) {
	const op = "storage.pq.UpdatePassword"

	// Проверям, не является ли изменяемый пользователь старшим админом
	if s.checkUpdatedUserIsSeniorAdmin(updaterId, userId) {
		return "err", fmt.Errorf("У вас не достаточно прав для этого")
	}
	// В ином случаем делаем запрос в бд на обновление данных
	stmt, err := s.db.Prepare(
		"UPDATE users SET pass_hash = $1 WHERE id = $2",
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", err.Error(), err)
	}
	_, err = stmt.ExecContext(ctx, newPassword, userId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", err.Error(), err)
	}
	return "Пароль успешно изменен", nil
}

// SaveUser записывает нового пользователя в бд
func (s *Storage) SaveUser(ctx context.Context, user *models.User) (int64, error) {
	const op = "storage.pq.SaveUser"
	// вызов метода генерации id
	userId := s.generateUserID()
	user.ID = userId
	// Простой запрос на добавление пользователя
	stmt, err := s.db.Prepare(
		"INSERT INTO users (id,email,pass_hash,name,surname,patronymic,phone_number,role) VALUES($8,$1,$2,$3,$4,$5,$6,$7)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// Выполняем запрос, передав параметры
	_, err = stmt.ExecContext(ctx, user.Email, user.PassHash, user.Name, user.Surname, user.Patronymic, user.PhoneNumber, user.Role, user.ID)
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
	// Записываем нового пользователя в map
	s.addUserToMap(user)

	return user.ID, nil
}

// EditUser редактирует данные пользователя в бд
func (s *Storage) EditUser(ctx context.Context, editedProfile *models.User, updaterId int64) error {
	const op = "storage.pq.SaveUser"

	// Если пользователь которому мы меняем роль - админ
	if s.checkUpdatedUserIsSeniorAdmin(updaterId, editedProfile.ID) {
		return fmt.Errorf("У вас нет на это прав")
	}

	// Простой запрос на добавление пользователя
	stmt, err := s.db.Prepare(
		"UPDATE users SET name = $1, surname = $2, patronymic = $3 WHERE id = $4")
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
	// Обновляем пользователя в map
	s.updateUser(editedProfile.ID, editedProfile.Name, editedProfile.Surname, editedProfile.Patronymic)
	return nil
}

func (s *Storage) EditUserRole(ctx context.Context, userId int64, newRole string, updaterId int64) (string, error) {
	const op = "storage.pq.EditUserRole"

	// Явяляется ли пользователь, которму мы изменяем роль админом высшего ранга?
	// Если да, то возвращаем ошибку
	// Если нет, то делаем

	if s.checkUpdatedUserIsSeniorAdmin(updaterId, userId) {
		return "err", fmt.Errorf("У вас нет прав для этого действия")
	} else if _, ok := s.adminLevels[userId]; ok && newRole != roleAdmin {
		// Если мы меняем роль админу
		defer s.removeFromAdminMap(userId)
	}

	// Если же роль пользователя отличается, то мы заменяем её

	stmt, err := s.db.Prepare("UPDATE users SET role = $1 WHERE id = $2")

	if err != nil {

		return "err", fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.ExecContext(ctx, newRole, userId)
	if err != nil {
		return "err", fmt.Errorf("%s: %w", op, err)
	}
	if newRole == roleAdmin {
		// Сохраняем в таблицу с уровнем админов
		updaterLevel := s.getUpdaterAdminLevel(updaterId) + 1
		s.saveNewAdmin(ctx, userId, updaterLevel)
		// Сохраняем в мапу с уровнем админства
	}
	// Сохранеяем изменения в map
	s.updateUserRole(userId, newRole)
	return "Роль успешно обновлена", nil
}
