package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage"
)

type SqliteTask struct {
	db *sql.DB
}

func (t *SqliteTask) Store(tsk item.Task) error {
	var recurStr string
	if tsk.Recurrer != nil {
		recurStr = tsk.Recurrer.String()
	}
	if _, err := t.db.Exec(`
INSERT INTO tasks
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
		tsk.ID, tsk.Title, tsk.Date.String(), tsk.Time.String(), tsk.Duration.String(), recurStr,
		tsk.Title, tsk.Date.String(), tsk.Time.String(), tsk.Duration.String(), recurStr); err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	return nil
}

func (t *SqliteTask) Find(id string) (item.Task, error) {
	var tsk item.Task
	var dateStr, timeStr, recurStr, durStr string
	err := t.db.QueryRow(`
SELECT id, title, date, time, duration, recurrer
FROM tasks
WHERE id = ?`, id).Scan(&tsk.ID, &tsk.Title, &dateStr, &timeStr, &durStr, &recurStr)
	switch {
	case err == sql.ErrNoRows:
		return item.Task{}, fmt.Errorf("event not found: %w", err)
	case err != nil:
		return item.Task{}, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	tsk.Date = item.NewDateFromString(dateStr)
	tsk.Time = item.NewTimeFromString(timeStr)
	dur, err := time.ParseDuration(durStr)
	if err != nil {
		return item.Task{}, fmt.Errorf("could not unmarshal recurrer: %v", err)
	}
	tsk.Duration = dur
	tsk.Recurrer = item.NewRecurrer(recurStr)

	return tsk, nil
}

func (t *SqliteTask) FindAll() ([]item.Task, error) {
	rows, err := t.db.Query(`
SELECT id, title, date, time, duration, recurrer
FROM tasks`)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	result := make([]item.Task, 0)
	defer rows.Close()
	for rows.Next() {
		var tsk item.Task
		var dateStr, timeStr, recurStr, durStr string
		if err := rows.Scan(&tsk.ID, &tsk.Title, &dateStr, &timeStr, &durStr, &recurStr); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
		dur, err := time.ParseDuration(durStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
		tsk.Date = item.NewDateFromString(dateStr)
		tsk.Time = item.NewTimeFromString(timeStr)
		tsk.Duration = dur
		tsk.Recurrer = item.NewRecurrer(recurStr)

		result = append(result, tsk)
	}

	return result, nil
}

func (s *SqliteTask) Delete(id string) error {
	result, err := s.db.Exec(`
DELETE FROM tasks
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
