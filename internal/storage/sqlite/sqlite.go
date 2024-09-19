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

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	res, err := stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	_ = res

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) error {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare(`INSERT INTO url(url, alias) VALUES (?, ?)`)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s : %w", op, storage.ErrorURLExists)
		}
		return fmt.Errorf("%s : %w", op, err)
	}
	_ = res
	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare(`SELECT url FROM url WHERE alias = ?`)
	if err != nil {
		return "", fmt.Errorf("%s : %w", op, err)
	}

	var url string
	err = stmt.QueryRow(alias).Scan(&url)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no url with alias '%s'", alias)
		} else {
			return "", fmt.Errorf("%s : %w", op, err)
		}
	}

	return url, nil
}

func (s *Storage) DeleteURL(alias string) error {
	op := "storage.sqlite.DeleteURL"

	stmt, err := s.db.Prepare(`DELETE FROM url WHERE alias = ?`)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}

	res, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	_ = res

	return nil
}

func (s *Storage) IsAliasUnique(alias string) (bool, error) {
	const op = "storage.sqlite.IsAliasUnique"

	var exists bool
	stmt := `SELECT exists(SELECT 1 FROM url WHERE alias=?)`
	err := s.db.QueryRow(stmt, alias).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s : %w", op, err)
	}

	return !exists, nil
}
