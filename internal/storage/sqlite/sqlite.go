package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"url_shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {
	const op = "internal.storage.sqlite.New"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
        CREATE TABLE IF NOT EXISTS url(
            id INTEGER PRIMARY KEY,
            alias TEXT NOT NULL UNIQUE,
            url TEXT NOT NULL
        );
        CREATE INDEX IF NOT EXISTS idx_alias ON url(alias)
    `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUrl(urlToSave string, alias string) (int64, error) {
	const op = "internal.storage.sqlite.SaveUrl"
	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	exec, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrAliasExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := exec.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get inserted id: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "internal.storage.sqlite.GetUrl"

	stmt, err := s.db.Prepare("SELECT url FROM url where alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	var url string
	if err = stmt.QueryRow(alias).Scan(&url); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrAliasNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return url, nil
}

// TODO impl func DeleteAlias
