package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

const (
	timestampFormat = "2006-01-02 15:04:05"
)

var migrations = []string{
	`CREATE TABLE items ("id" TEXT UNIQUE, "kind" TEXT, "updated" TIMESTAMP, "deleted" INTEGER, "body" TEXT)`,
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

func (s *Sqlite) Update(item Item) error {
	if _, err := s.db.Exec(`
INSERT INTO items
(id, kind, updated, deleted, body)
VALUES
(?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE
SET
kind=excluded.kind,
updated=excluded.updated,
deleted=excluded.deleted,
body=excluded.body`,
		item.ID, item.Kind, item.Updated.Format(timestampFormat), item.Deleted, item.Body); err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	return nil
}

func (s *Sqlite) Updated(ks []Kind, t time.Time) ([]Item, error) {
	query := `
SELECT id, kind, updated, deleted, body
FROM items
WHERE updated > ?`
	var rows *sql.Rows
	var err error
	if len(ks) == 0 {
		rows, err = s.db.Query(query, t.Format(timestampFormat))
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
	} else {
		args := []any{t.Format(timestampFormat)}
		ph := make([]string, 0, len(ks))
		for _, k := range ks {
			args = append(args, string(k))
			ph = append(ph, "?")
		}
		query = fmt.Sprintf("%s AND kind in (%s)", query, strings.Join(ph, ","))
		rows, err = s.db.Query(query, args...)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
	}

	result := make([]Item, 0)
	defer rows.Close()
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Kind, &item.Updated, &item.Deleted, &item.Body); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
		result = append(result, item)
	}

	return result, nil
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
