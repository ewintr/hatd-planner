package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"go-mod.ewintr.nl/planner/item"
)

const (
	timestampFormat = "2006-01-02 15:04:05"
)

var migrations = []string{
	`CREATE TABLE items (id TEXT PRIMARY KEY, kind TEXT, updated TIMESTAMP, deleted BOOLEAN, body TEXT)`,
	`CREATE INDEX idx_items_updated ON items(updated)`,
	`CREATE INDEX idx_items_kind ON items(kind)`,
	`ALTER TABLE items ADD COLUMN recurrer JSONB, ADD COLUMN recur_next TIMESTAMP`,
}

var (
	ErrInvalidConfiguration     = errors.New("invalid configuration")
	ErrIncompatibleSQLMigration = errors.New("incompatible migration")
	ErrNotEnoughSQLMigrations   = errors.New("already more migrations than wanted")
	ErrPostgresFailure          = errors.New("postgres returned an error")
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(host, port, dbname, user, password string) (*Postgres, error) {
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", host, port, dbname, user, password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidConfiguration, err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidConfiguration, err)
	}

	p := &Postgres{
		db: db,
	}

	if err := p.migrate(migrations); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Postgres) Update(item item.Item, ts time.Time) error {
	_, err := p.db.Exec(`
		INSERT INTO items (id, kind, updated, deleted, body, recurrer, recur_next)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE
		SET kind = EXCLUDED.kind,
			updated = EXCLUDED.updated,
			deleted = EXCLUDED.deleted,
			body = EXCLUDED.body,
			recurrer = EXCLUDED.recurrer,
			recur_next = EXCLUDED.recur_next`,
		item.ID, item.Kind, ts, item.Deleted, item.Body, item.Recurrer, item.RecurNext)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrPostgresFailure, err)
	}
	return nil
}

func (p *Postgres) Updated(ks []item.Kind, t time.Time) ([]item.Item, error) {
	query := `
		SELECT id, kind, updated, deleted, body, recurrer, recur_next
		FROM items
		WHERE updated > $1`
	args := []interface{}{t}

	if len(ks) > 0 {
		placeholder := make([]string, len(ks))
		for i := range ks {
			placeholder[i] = fmt.Sprintf("$%d", i+2)
			args = append(args, string(ks[i]))
		}
		query += fmt.Sprintf(" AND kind = ANY(ARRAY[%s])", strings.Join(placeholder, ","))
	}

	rows, err := p.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPostgresFailure, err)
	}
	defer rows.Close()

	result := make([]item.Item, 0)
	for rows.Next() {
		var item item.Item
		var recurNext sql.NullTime
		if err := rows.Scan(&item.ID, &item.Kind, &item.Updated, &item.Deleted, &item.Body, &item.Recurrer, &recurNext); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrPostgresFailure, err)
		}
		if recurNext.Valid {
			item.RecurNext = recurNext.Time
		}
		result = append(result, item)
	}

	return result, nil
}

func (p *Postgres) RecursBefore(date time.Time) ([]item.Item, error) {
	query := `
		SELECT id, kind, updated, deleted, body, recurrer, recur_next
		FROM items
		WHERE recur_next <= $1 AND recurrer IS NOT NULL`

	rows, err := p.db.Query(query, date)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPostgresFailure, err)
	}
	defer rows.Close()

	result := make([]item.Item, 0)
	for rows.Next() {
		var item item.Item
		var recurNext sql.NullTime
		if err := rows.Scan(&item.ID, &item.Kind, &item.Updated, &item.Deleted, &item.Body, &item.Recurrer, &recurNext); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrPostgresFailure, err)
		}
		if recurNext.Valid {
			item.RecurNext = recurNext.Time
		}
		result = append(result, item)
	}

	return result, nil
}

func (p *Postgres) RecursNext(id string, date time.Time, ts time.Time) error {
	var recurrer *item.Recur
	err := p.db.QueryRow(`
		SELECT recurrer
		FROM items
		WHERE id = $1`, id).Scan(&recurrer)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return fmt.Errorf("%w: %v", ErrPostgresFailure, err)
	}

	if recurrer == nil {
		return ErrNotARecurrer
	}

	// Verify that the new date is actually a valid recurrence
	if !recurrer.On(date) {
		return fmt.Errorf("%w: date %v is not a valid recurrence", ErrPostgresFailure, date)
	}

	_, err = p.db.Exec(`
		UPDATE items 
		SET recur_next = $1,
		  updated = $2  
		WHERE id = $3`, date, ts, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrPostgresFailure, err)
	}

	return nil
}

func (p *Postgres) migrate(wanted []string) error {
	// Create migration table if not exists
	_, err := p.db.Exec(`
		CREATE TABLE IF NOT EXISTS migration
		(id SERIAL PRIMARY KEY, query TEXT)
	`)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrPostgresFailure, err)
	}

	// Find existing migrations
	rows, err := p.db.Query(`SELECT query FROM migration ORDER BY id`)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrPostgresFailure, err)
	}
	defer rows.Close()

	var existing []string
	for rows.Next() {
		var query string
		if err := rows.Scan(&query); err != nil {
			return fmt.Errorf("%w: %v", ErrPostgresFailure, err)
		}
		existing = append(existing, query)
	}

	// Compare and execute missing migrations
	missing, err := compareMigrations(wanted, existing)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrPostgresFailure, err)
	}

	for _, query := range missing {
		if _, err := p.db.Exec(query); err != nil {
			return fmt.Errorf("%w: %v", ErrPostgresFailure, err)
		}

		// Register migration
		if _, err := p.db.Exec(`
			INSERT INTO migration (query) VALUES ($1)
		`, query); err != nil {
			return fmt.Errorf("%w: %v", ErrPostgresFailure, err)
		}
	}

	return nil
}

func compareMigrations(wanted, existing []string) ([]string, error) {
	var needed []string
	if len(wanted) < len(existing) {
		return nil, ErrNotEnoughSQLMigrations
	}

	for i, want := range wanted {
		switch {
		case i >= len(existing):
			needed = append(needed, want)
		case want == existing[i]:
			// do nothing
		case want != existing[i]:
			return nil, fmt.Errorf("%w: %v", ErrIncompatibleSQLMigration, want)
		}
	}

	return needed, nil
}
