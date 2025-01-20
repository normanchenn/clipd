package storage

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/normanchenn/clipd/daemon/internal/config"
)

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(config config.Config) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", config.DbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite db: %w", err)
	}

	query := `
CREATE TABLE IF NOT EXISTS clipboard (
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP PRIMARY KEY,
    id TEXT UNIQUE NOT NULL,
    data TEXT NOT NULL
);
    `
	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	return &SQLiteStorage{
		db: db,
	}, nil
}

// created timestamps are done by the database
func (s *SQLiteStorage) Save(entry Entry) (*Entry, error) {
	query := `INSERT INTO clipboard (id, data) VALUES (?, ?)`
	_, err := s.db.Exec(query, entry.ID, entry.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to save entry: %w", err)
	}
	return &entry, nil
}

func (s *SQLiteStorage) GetByIndex(index int) (Entry, error) {
	query := `SELECT id, created, data from clipboard ORDER BY created DESC LIMIT 1 OFFSET ?`
	entry, err := s.fetchSingleEntry(query, index)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Entry{}, ErrIndexOutOfRange
		}
		return Entry{}, err
	}
	return entry, nil
}

func (s *SQLiteStorage) GetByID(id string) (Entry, error) {
	query := `SELECT id, created, data from clipboard WHERE id = ?`
	entry, err := s.fetchSingleEntry(query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Entry{}, ErrIdDoesntExist
		}
		return Entry{}, err
	}
	return entry, nil
}

func (s *SQLiteStorage) GetRange(start int, end int) ([]Entry, error) {
	if start >= end {
		return nil, ErrInvalidRange
	}

	query := `SELECT id, created, data FROM clipboard ORDER BY created DESC LIMIT ? OFFSET ?`
	limit := end - start

	entries, err := s.fetchMultipleEntries(query, limit, start)
	if err != nil {
		return nil, err
	}

	if len(entries) != limit {
		return entries, ErrIncompleteRange

	}
	return entries, nil
}

func (s *SQLiteStorage) GetAll() ([]Entry, error) {
	query := `SELECT id, created, data FROM clipboard ORDER BY created DESC`
	return s.fetchMultipleEntries(query)
}

func (s *SQLiteStorage) Size() (int, error) {
	query := `SELECT COUNT(*) FROM clipboard`
	var count int
	err := s.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to scan entry: %w", err)
	}
	return count, nil
}

func (s *SQLiteStorage) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close SQLite database: %w", err)
	}
	return nil
}

func (s *SQLiteStorage) fetchSingleEntry(query string, args ...interface{}) (Entry, error) {
	var entry Entry

	row := s.db.QueryRow(query, args...)
	err := row.Scan(&entry.ID, &entry.Created, &entry.Data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Entry{}, fmt.Errorf("no entry found when querying single row (query: %q, args: %v): %w", query, args, err)
		}
		return Entry{}, fmt.Errorf("failed to fetch entry (query: %q, args: %v): %w", query, args, err)
	}
	return entry, nil

}

func (s *SQLiteStorage) fetchMultipleEntries(query string, args ...interface{}) ([]Entry, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query (query: %q, args: %v): %w", query, args, err)
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var entry Entry
		// data must be returned in this order
		err := rows.Scan(&entry.ID, &entry.Created, &entry.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row (query: %q, args: %v): %w", query, args, err)
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
