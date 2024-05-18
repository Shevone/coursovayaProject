package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
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
	return &Storage{db: db}, nil
}
