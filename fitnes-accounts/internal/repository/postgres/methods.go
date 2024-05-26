package postgres

import (
	"context"
	"fitnes-account/internal/models"
	"fmt"
	"math/rand"
	"time"
)

func (s *Storage) loadAdminLevels() error {
	const op = "storage.pq.loadAdminLevels"
	s.adminLevelsMutex.Lock()
	defer s.adminLevelsMutex.Unlock()
	// сделать запрос, на выборку из таблицы с уровнем админства
	// занести это все в map
	// поместить мапу в поле структуры
	stmt, err := s.db.
		Prepare(
			"SELECT admin_id, admin_level FROM admin_level",
		)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	row, err := stmt.Query()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer row.Close()
	for row.Next() {
		var adminID int64
		var adminLevel int
		err = row.Scan(&adminID, &adminLevel)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		s.adminLevels[adminID] = adminLevel
	}
	if err = row.Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// загружаем пользователей в map
func (s *Storage) loadUsers() error {
	s.usersMutex.Lock()
	defer s.usersMutex.Unlock()
	// Взять из бд всех пользователей и поместить их сюда
	const op = "storage.pq.loadUsers"
	users, err := s.getUsers(context.Background())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	for _, user := range users {
		s.userIds = append(s.userIds, user.ID)
		s.users[user.ID] = &user
	}
	return nil
}

func (s *Storage) generateUserID() int64 {
	s.usersMutex.Lock()
	defer s.usersMutex.Unlock()

	rand.Seed(time.Now().UnixNano())
	var id int64
	for {
		id = rand.Int63n(65000) // Генерация в диапазоне от 0 до максимального int64
		if _, ok := s.users[id]; !ok {
			break
		}
	}

	return id
}
func (s *Storage) addUserToMap(user *models.User) {
	s.usersMutex.Lock()
	defer s.usersMutex.Unlock()
	s.userIds = append(s.userIds, user.ID)
	s.users[user.ID] = user
}

func (s *Storage) updateUser(userId int64, userName string, userSurname string, userPatron string) {
	s.usersMutex.Lock()
	defer s.usersMutex.Unlock()
	s.users[userId].Name = userName
	s.users[userId].Surname = userSurname
	s.users[userId].Patronymic = userPatron
}

func (s *Storage) updateUserRole(userId int64, userRole string) {
	s.usersMutex.Lock()
	defer s.usersMutex.Unlock()
	s.users[userId].Role = userRole
}

func (s *Storage) getUpdaterAdminLevel(updaterId int64) int {

	return s.adminLevels[updaterId]
}

func (s *Storage) checkUpdatedUserIsSeniorAdmin(updaterId int64, userId int64) bool {
	s.adminLevelsMutex.Lock()
	defer s.adminLevelsMutex.Unlock()
	if updaterId == userId {
		return false
	}
	adminLevel, ok := s.adminLevels[userId]
	if ok {

		if adminLevel <= s.adminLevels[updaterId] {
			return true
		}
	}
	return false
}

func (s *Storage) saveNewAdmin(ctx context.Context, newAdminId int64, newAdminLevel int) error {
	// Сохраняем в бд
	s.adminLevelsMutex.Lock()
	defer s.adminLevelsMutex.Unlock()
	stmt, err := s.db.Prepare("INSERT INTO admin_level(admin_id, admin_level) VALUES ($1,$2)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, newAdminId, newAdminLevel)
	if err != nil {
		return err
	}
	// Сохраняем в мапу
	s.adminLevels[newAdminId] = newAdminLevel
	return nil
}
func (s *Storage) removeFromAdminMap(removeId int64) {
	s.adminLevelsMutex.Lock()
	defer s.adminLevelsMutex.Unlock()
	delete(s.adminLevels, removeId)
}
