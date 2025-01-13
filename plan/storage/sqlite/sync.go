package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage"
)

type Sync struct {
	tx *storage.Tx
}

func (s *Sync) FindAll() ([]item.Item, error) {
	rows, err := s.tx.Query("SELECT id, kind, updated, deleted, date, recurrer, recur_next, body FROM items")
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

func (s *Sync) Store(i item.Item) error {
	if i.Updated.IsZero() {
		i.Updated = time.Now()
	}
	var recurStr string
	if i.Recurrer != nil {
		recurStr = i.Recurrer.String()
	}

	_, err := s.tx.Exec(
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

func (s *Sync) DeleteAll() error {
	_, err := s.tx.Exec("DELETE FROM items")
	if err != nil {
		return fmt.Errorf("%w: failed to delete all items: %v", ErrSqliteFailure, err)
	}
	return nil
}

func (s *Sync) SetLastUpdate(ts time.Time) error {
	if _, err := s.tx.Exec(`UPDATE syncupdate SET timestamp = ?`, ts.Format(time.RFC3339)); err != nil {
		return fmt.Errorf("%w: could not store timestamp: %v", ErrSqliteFailure, err)
	}
	return nil
}

func (s *Sync) LastUpdate() (time.Time, error) {
	var tsStr string
	if err := s.tx.QueryRow("SELECT timestamp FROM syncupdate").Scan(&tsStr); err != nil {
		return time.Time{}, fmt.Errorf("%w: failed to get last update: %v", ErrSqliteFailure, err)
	}
	ts, err := time.Parse(time.RFC3339, tsStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: could not convert db timstamp into time.Time: %v", ErrSqliteFailure, err)
	}

	return ts, nil
}
