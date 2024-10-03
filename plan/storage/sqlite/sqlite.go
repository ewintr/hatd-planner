package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

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
	`CREATE TABLE localids ("id" TEXT UNIQUE, "local_id" INTEGER)`,
}

var (
	ErrInvalidConfiguration     = errors.New("invalid configuration")
	ErrIncompatibleSQLMigration = errors.New("incompatible migration")
	ErrNotEnoughSQLMigrations   = errors.New("already more migrations than wanted")
	ErrSqliteFailure            = errors.New("sqlite returned an error")
)

func NewSqlites(dbPath string) (*LocalID, *SqliteEvent, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", ErrInvalidConfiguration, err)
	}

	sl := &LocalID{
		db: db,
	}
	se := &SqliteEvent{
		db: db,
	}

	if err := migrate(db, migrations); err != nil {
		return nil, nil, err
	}

	return sl, se, nil
}

func migrate(db *sql.DB, wanted []string) error {
	// admin table
	if _, err := db.Exec(`
CREATE TABLE IF NOT EXISTS migration
("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "query" TEXT)
`); err != nil {
		return err
	}

	// find existing
	rows, err := db.Query(`SELECT query FROM migration ORDER BY id`)
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
		if _, err := db.Exec(string(query)); err != nil {
			return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}

		// register
		if _, err := db.Exec(`
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
