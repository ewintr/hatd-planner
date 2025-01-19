package sqlite

import (
	"fmt"

	"go-mod.ewintr.nl/planner/item"
	"go-mod.ewintr.nl/planner/plan/storage"
)

type SqliteSchedule struct {
	tx *storage.Tx
}

func (ss *SqliteSchedule) Store(sched item.Schedule) error {
	var recurStr string
	if sched.Recurrer != nil {
		recurStr = sched.Recurrer.String()
	}
	if _, err := ss.tx.Exec(`
INSERT INTO schedules
(id, title, date, recur)
VALUES
(?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE
SET
title=?,
date=?,
recur=?
`,
		sched.ID, sched.Title, sched.Date.String(), recurStr,
		sched.Title, sched.Date.String(), recurStr); err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	return nil
}

func (ss *SqliteSchedule) Find(start, end item.Date) ([]item.Schedule, error) {
	rows, err := ss.tx.Query(`SELECT
id, title, date, recur
FROM schedules
WHERE date >= ? AND date <= ?`, start.String(), end.String())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	defer rows.Close()
	scheds := make([]item.Schedule, 0)
	for rows.Next() {
		var sched item.Schedule
		var dateStr, recurStr string
		if err := rows.Scan(&sched.ID, &sched.Title, &dateStr, &recurStr); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
		sched.Date = item.NewDateFromString(dateStr)
		sched.Recurrer = item.NewRecurrer(recurStr)

		scheds = append(scheds, sched)
	}

	return scheds, nil
}

func (ss *SqliteSchedule) Delete(id string) error {

	result, err := ss.tx.Exec(`
DELETE FROM schedules
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
