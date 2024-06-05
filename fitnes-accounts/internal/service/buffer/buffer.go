package buffer

import (
	"fitnes-account/internal/models"
	"sync"
)

type Buffer struct {
	// =======
	// ========================
	adminLevels      map[int64]int
	adminLevelsMutex sync.Mutex
	// ========================
	userIds    []int64
	users      map[int64]*models.User
	usersMutex sync.Mutex
	// ========================
	loginPassword      map[string]string
	loginPasswordMutex sync.Mutex
	// ========================
}

func NewBuffer(adminLevels map[int64]int, users map[int64]*models.User, loginPassword map[string]string) *Buffer {
	return &Buffer{
		adminLevels:        adminLevels,
		adminLevelsMutex:   sync.Mutex{},
		users:              users,
		usersMutex:         sync.Mutex{},
		loginPassword:      loginPassword,
		loginPasswordMutex: sync.Mutex{},
	}
}
