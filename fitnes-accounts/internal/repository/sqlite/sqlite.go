package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	DbPath string `default:"./storage/database.db"`
}
type Storage struct {
	db *sql.DB
}

func NewSqliteRepository(cfg *Config) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", cfg.DbPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}
