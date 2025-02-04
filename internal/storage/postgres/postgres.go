package postgres

import (
	"awesomeProject/internal/error-handler"
	"awesomeProject/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(storagePath string) (*Storage, error) {
	const op = "storage.postgres.NewStorage"
	db, err := sql.Open("postgres", storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS url (
    	id INTEGER PRIMARY KEY,
    	alias TEXT UNIQUE NOT NULL,
		url TEXT NOT NULL);`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_alias ON url (alias);`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) error {

	const op = "storage.postgres.SaveURL"
	stmt, err := s.db.Prepare(`INSERT INTO url (alias, url) VALUES ($1, $2)`)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()
	_, err = stmt.Exec(alias, urlToSave)
	if err != nil {
		_ = errorHandler.OnHandleError(err, op)
	}
	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"

	stmt, err := s.db.Prepare(`SELECT url FROM url WHERE alias = $1`)
	if err != nil {
		return "", fmt.Errorf("Error when prepare sql request  %s: %w", op, err)
	}

	defer stmt.Close()
	var url string
	err = stmt.QueryRow(alias).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("Row not found  %s: %w", op, storage.ErrNotFound)
		}
		_ = errorHandler.OnHandleError(err, op)
	}
	return url, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.postgres.DeleteURL"
	_, err := s.db.Exec(`DELETE FROM url WHERE alias = $1`, alias)
	if err != nil {
		_ = errorHandler.OnHandleError(err, op)
	}
	return nil

}
