package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go-mod.ewintr.nl/planner/item"
	_ "modernc.org/sqlite"
)

const (
	timestampFormat = "2006-01-02 15:04:05"
)

var migrations = []string{
	`CREATE TABLE events ("id" TEXT UNIQUE, "title" TEXT, "start" TIMESTAMP, "duration" TEXT)`,
	`PRAGMA journal_mode=WAL`,
	`PRAGMA synchronous=NORMAL`,
	`PRAGMA cache_size=2000`,
}

var (
	ErrInvalidConfiguration     = errors.New("invalid configuration")
	ErrIncompatibleSQLMigration = errors.New("incompatible migration")
	ErrNotEnoughSQLMigrations   = errors.New("already more migrations than wanted")
	ErrSqliteFailure            = errors.New("sqlite returned an error")
)

type Sqlite struct {
	db *sql.DB
}

func NewSqlite(dbPath string) (*Sqlite, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return &Sqlite{}, fmt.Errorf("%w: %v", ErrInvalidConfiguration, err)
	}

	s := &Sqlite{
		db: db,
	}

	if err := s.migrate(migrations); err != nil {
		return &Sqlite{}, err
	}

	return s, nil
}

func (s *Sqlite) Store(event item.Event) error {
	if _, err := s.db.Exec(`
INSERT INTO events
(id, title, start, duration)
VALUES
(?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE
SET
title=?,
start=?,
duration=?`,
		event.ID, event.Title, event.Start.Format(timestampFormat), event.Duration.String(),
		event.Title, event.Start.Format(timestampFormat), event.Duration.String()); err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	return nil
}

func (s *Sqlite) Find(id string) (item.Event, error) {
	var event item.Event
	var durStr string
	err := s.db.QueryRow(`
SELECT id, title, start, duration
FROM events
WHERE id = ?`, id).Scan(&event.ID, &event.Title, &event.Start, &durStr)
	switch {
	case err == sql.ErrNoRows:
		return item.Event{}, fmt.Errorf("event not found: %w", err)
	case err != nil:
		return item.Event{}, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	dur, err := time.ParseDuration(durStr)
	if err != nil {
		return item.Event{}, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	event.Duration = dur

	return event, nil
}

func (s *Sqlite) FindAll() ([]item.Event, error) {
	rows, err := s.db.Query(`
SELECT id, title, start, duration
FROM events`)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	result := make([]item.Event, 0)
	defer rows.Close()
	for rows.Next() {
		var event item.Event
		var durStr string
		if err := rows.Scan(&event.ID, &event.Title, &event.Start, &durStr); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
		dur, err := time.ParseDuration(durStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
		event.Duration = dur
		result = append(result, event)
	}

	return result, nil
}

func (s *Sqlite) Delete(id string) error {
	result, err := s.db.Exec(`
DELETE FROM events
WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("event not found: %s", id)
	}

	return nil
}

func (s *Sqlite) migrate(wanted []string) error {
	// admin table
	if _, err := s.db.Exec(`
CREATE TABLE IF NOT EXISTS migration
("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "query" TEXT)
`); err != nil {
		return err
	}

	// find existing
	rows, err := s.db.Query(`SELECT query FROM migration ORDER BY id`)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	existing := []string{}
	for rows.Next() {
		var query string
		if err := rows.Scan(&query); err != nil {
			return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
		existing = append(existing, string(query))
	}
	rows.Close()

	// compare
	missing, err := compareMigrations(wanted, existing)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	// execute missing
	for _, query := range missing {
		if _, err := s.db.Exec(string(query)); err != nil {
			return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}

		// register
		if _, err := s.db.Exec(`
INSERT INTO migration
(query) VALUES (?)
`, query); err != nil {
			return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
	}

	return nil
}

func compareMigrations(wanted, existing []string) ([]string, error) {
	needed := []string{}
	if len(wanted) < len(existing) {
		return []string{}, ErrNotEnoughSQLMigrations
	}

	for i, want := range wanted {
		switch {
		case i >= len(existing):
			needed = append(needed, want)
		case want == existing[i]:
			// do nothing
		case want != existing[i]:
			return []string{}, fmt.Errorf("%w: %v", ErrIncompatibleSQLMigration, want)
		}
	}

	return needed, nil
}
