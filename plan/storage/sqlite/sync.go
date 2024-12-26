package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"go-mod.ewintr.nl/planner/item"
)

type SqliteSync struct {
	db *sql.DB
}

func NewSqliteSync(db *sql.DB) *SqliteSync {
	return &SqliteSync{db: db}
}

func (s *SqliteSync) FindAll() ([]item.Item, error) {
	rows, err := s.db.Query("SELECT id, kind, updated, deleted, date, recurrer, recur_next, body FROM items")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to query items: %v", ErrSqliteFailure, err)
	}
	defer rows.Close()

	var items []item.Item
	for rows.Next() {
		var i item.Item
		var updatedStr, dateStr, recurStr, recurNextStr string
		err := rows.Scan(&i.ID, &i.Kind, &updatedStr, &i.Deleted, &dateStr, &recurStr, &recurNextStr, &i.Body)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to scan item: %v", ErrSqliteFailure, err)
		}
		i.Updated, err = time.Parse(time.RFC3339, updatedStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse updated time: %v", err)
		}
		i.Date = item.NewDateFromString(dateStr)
		i.Recurrer = item.NewRecurrer(recurStr)
		i.RecurNext = item.NewDateFromString(recurNextStr)

		items = append(items, i)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return items, nil
}

func (s *SqliteSync) Store(i item.Item) error {
	if i.Updated.IsZero() {
		i.Updated = time.Now()
	}
	var recurStr string
	if i.Recurrer != nil {
		recurStr = i.Recurrer.String()
	}

	_, err := s.db.Exec(
		`INSERT OR REPLACE INTO items (id, kind, updated, deleted, date, recurrer, recur_next, body)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		i.ID,
		i.Kind,
		i.Updated.UTC().Format(time.RFC3339),
		i.Deleted,
		i.Date.String(),
		recurStr,
		i.RecurNext.String(),
		sql.NullString{String: i.Body, Valid: i.Body != ""}, // This allows empty string but not NULL
	)
	if err != nil {
		return fmt.Errorf("%w: failed to store item: %v", ErrSqliteFailure, err)
	}
	return nil
}

func (s *SqliteSync) DeleteAll() error {
	_, err := s.db.Exec("DELETE FROM items")
	if err != nil {
		return fmt.Errorf("%w: failed to delete all items: %v", ErrSqliteFailure, err)
	}
	return nil
}

func (s *SqliteSync) LastUpdate() (time.Time, error) {
	var updatedStr sql.NullString
	err := s.db.QueryRow("SELECT MAX(updated) FROM items").Scan(&updatedStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: failed to get last update: %v", ErrSqliteFailure, err)
	}

	if !updatedStr.Valid {
		return time.Time{}, nil // Return zero time if NULL or no rows
	}

	lastUpdate, err := time.Parse(time.RFC3339, updatedStr.String)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: failed to parse last update time: %v", ErrSqliteFailure, err)
	}
	return lastUpdate, nil
}
