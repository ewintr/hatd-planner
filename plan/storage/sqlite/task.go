package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage"
)

type SqliteTask struct {
	tx *storage.Tx
}

func (t *SqliteTask) Store(tsk item.Task) error {
	var recurStr string
	if tsk.Recurrer != nil {
		recurStr = tsk.Recurrer.String()
	}
	if _, err := t.tx.Exec(`
INSERT INTO tasks
(id, title, project, date, time, duration, recurrer)
VALUES
(?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE
SET
title=?,
project=?,
date=?,
time=?,
duration=?,
recurrer=?
`,
		tsk.ID, tsk.Title, tsk.Project, tsk.Date.String(), tsk.Time.String(), tsk.Duration.String(), recurStr,
		tsk.Title, tsk.Project, tsk.Date.String(), tsk.Time.String(), tsk.Duration.String(), recurStr); err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	return nil
}

func (t *SqliteTask) FindOne(id string) (item.Task, error) {
	var tsk item.Task
	var dateStr, timeStr, recurStr, durStr string
	err := t.tx.QueryRow(`
SELECT id, title, project, date, time, duration, recurrer
FROM tasks
WHERE id = ?`, id).Scan(&tsk.ID, &tsk.Title, &tsk.Project, &dateStr, &timeStr, &durStr, &recurStr)
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

func (t *SqliteTask) FindMany(params storage.TaskListParams) ([]item.Task, error) {
	query := `SELECT id, title, project, date, time, duration, recurrer FROM tasks`
	args := []interface{}{}

	where := make([]string, 0)
	var dateNonEmpty bool
	if params.HasRecurrer {
		where = append(where, `recurrer != ''`)
	}
	if !params.From.IsZero() {
		where = append(where, `date >= ?`)
		args = append(args, params.From.String())
		dateNonEmpty = true
	}
	if !params.To.IsZero() {
		where = append(where, `date <= ?`)
		args = append(args, params.To.String())
		dateNonEmpty = true
	}
	if params.Project != "" {
		where = append(where, `project = ?`)
		args = append(args, params.Project)
	}
	if dateNonEmpty {
		where = append(where, `date != ""`)
	}

	if len(where) > 0 {
		query += ` WHERE ` + where[0]
		for _, w := range where[1:] {
			query += ` AND ` + w
		}
	}

	rows, err := t.tx.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	tasks := make([]item.Task, 0)
	defer rows.Close()
	for rows.Next() {
		var tsk item.Task
		var dateStr, timeStr, recurStr, durStr string
		if err := rows.Scan(&tsk.ID, &tsk.Title, &tsk.Project, &dateStr, &timeStr, &durStr, &recurStr); err != nil {
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

		tasks = append(tasks, tsk)
	}

	return tasks, nil
}

func (t *SqliteTask) Delete(id string) error {
	result, err := t.tx.Exec(`
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

func (t *SqliteTask) Projects() (map[string]int, error) {
	rows, err := t.tx.Query(`SELECT project, count(*) FROM tasks GROUP BY project`)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	result := make(map[string]int)
	defer rows.Close()
	for rows.Next() {
		var pr string
		var count int
		if err := rows.Scan(&pr, &count); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
		result[pr] = count
	}

	return result, nil
}
