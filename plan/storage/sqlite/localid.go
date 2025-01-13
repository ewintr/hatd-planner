package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"go-mod.ewintr.nl/planner/plan/storage"
)

type LocalID struct {
	tx *storage.Tx
}

func (l *LocalID) FindOne(lid int) (string, error) {
	var id string
	err := l.tx.QueryRow(`
SELECT id
FROM localids
WHERE local_id = ?
`, lid).Scan(&id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return "", storage.ErrNotFound
	case err != nil:
		return "", fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	return id, nil
}

func (l *LocalID) FindAll() (map[string]int, error) {
	rows, err := l.tx.Query(`
SELECT id, local_id
FROM localids
`)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	result := make(map[string]int)
	defer rows.Close()
	for rows.Next() {
		var id string
		var localID int
		if err := rows.Scan(&id, &localID); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
		result[id] = localID
	}

	return result, nil
}

func (l *LocalID) FindOrNext(id string) (int, error) {
	return 0, nil
}

func (l *LocalID) Next() (int, error) {
	idMap, err := l.FindAll()
	if err != nil {
		return 0, err
	}
	cur := make([]int, 0, len(idMap))
	for _, localID := range idMap {
		cur = append(cur, localID)
	}

	return storage.NextLocalID(cur), nil
}

func (l *LocalID) Store(id string, localID int) error {
	if _, err := l.tx.Exec(`
INSERT INTO localids
(id, local_id)
VALUES
(? ,?)`, id, localID); err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	return nil
}

func (l *LocalID) Delete(id string) error {
	result, err := l.tx.Exec(`
DELETE FROM localids
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
