package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage"
)

type SqliteEvent struct {
	db *sql.DB
}

func (s *SqliteEvent) Store(event item.Event) error {
	var recurStr string
	if event.Recurrer != nil {
		recurStr = event.Recurrer.String()
	}
	if _, err := s.db.Exec(`
INSERT INTO events
(id, title, date, time, duration, recurrer)
VALUES
(?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE
SET
title=?,
date=?,
time=?,
duration=?,
recurrer=?
`,
		event.ID, event.Title, event.Date.String(), event.Time.String(), event.Duration.String(), recurStr,
		event.Title, event.Date.String(), event.Time.String(), event.Duration.String(), recurStr); err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	return nil
}

func (s *SqliteEvent) Find(id string) (item.Event, error) {
	var event item.Event
	var dateStr, timeStr, recurStr, durStr string
	err := s.db.QueryRow(`
SELECT id, title, date, time, duration, recurrer
FROM events
WHERE id = ?`, id).Scan(&event.ID, &event.Title, &dateStr, &timeStr, &durStr, &recurStr)
	switch {
	case err == sql.ErrNoRows:
		return item.Event{}, fmt.Errorf("event not found: %w", err)
	case err != nil:
		return item.Event{}, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	event.Date = item.NewDateFromString(dateStr)
	event.Time = item.NewTimeFromString(timeStr)
	dur, err := time.ParseDuration(durStr)
	if err != nil {
		return item.Event{}, fmt.Errorf("could not unmarshal recurrer: %v", err)
	}
	event.Duration = dur
	event.Recurrer = item.NewRecurrer(recurStr)

	return event, nil
}

func (s *SqliteEvent) FindAll() ([]item.Event, error) {
	rows, err := s.db.Query(`
SELECT id, title, date, time, duration, recurrer
FROM events`)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	result := make([]item.Event, 0)
	defer rows.Close()
	for rows.Next() {
		var event item.Event
		var dateStr, timeStr, recurStr, durStr string
		if err := rows.Scan(&event.ID, &event.Title, &dateStr, &timeStr, &durStr, &recurStr); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
		dur, err := time.ParseDuration(durStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
		event.Date = item.NewDateFromString(dateStr)
		event.Time = item.NewTimeFromString(timeStr)
		event.Duration = dur
		event.Recurrer = item.NewRecurrer(recurStr)

		result = append(result, event)
	}

	return result, nil
}

func (s *SqliteEvent) Delete(id string) error {
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
		return storage.ErrNotFound
	}

	return nil
}
