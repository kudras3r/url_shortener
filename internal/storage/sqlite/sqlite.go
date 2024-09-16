package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/kudras3r/url_shortener/internal/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	qry, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	res, err := qry.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	_ = res

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) error {
	const op = "storage.sqlite.SaveURL"

	qry, err := s.db.Prepare(`INSERT INTO url(url, alias) VALUES (?, ?)`)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}

	res, err := qry.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s : %w", op, storage.ErrorURLExists)
		}
		return fmt.Errorf("%s : %w", op, err)
	}
	_ = res
	return nil
}
