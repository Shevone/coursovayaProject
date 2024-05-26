package postgres

import (
	"database/sql"
	"fitnes-account/internal/models"
	"fmt"
	_ "github.com/lib/pq"
	"sync"
)

type Config struct {
	Host     string `default:"localhost"`
	Port     int    `default:"5432"`
	User     string `default:""`
	Password string `default:""`
	Database string `default:""`
}
type Storage struct {
	db *sql.DB
	// ========================
	adminLevels      map[int64]int
	adminLevelsMutex sync.Mutex
	// ========================
	userIds    []int64
	users      map[int64]*models.User
	usersMutex sync.Mutex
}

func NewPostgresRepository(cfg *Config) (*Storage, error) {
	const op = "storage.pq.New"
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	repo := &Storage{db: db}
	repo.adminLevels = make(map[int64]int)
	repo.userIds = make([]int64, 0)
	repo.users = make(map[int64]*models.User)
	err = repo.loadAdminLevels()
	if err != nil {
		return nil, err
	}
	err = repo.loadUsers()
	if err != nil {
		return nil, err
	}
	return repo, nil
}
