package storage

import "database/sql"

// Tx wraps sql.Tx so transactions can be skipped for in-memory repositories
type Tx struct {
	tx *sql.Tx
}

func NewTx(tx *sql.Tx) *Tx {
	return &Tx{tx: tx}
}

func (tx *Tx) Rollback() error {
	if tx.tx == nil {
		return nil
	}

	return tx.tx.Rollback()
}

func (tx *Tx) Commit() error {
	if tx.tx == nil {
		return nil
	}

	return tx.tx.Commit()
}

func (tx *Tx) QueryRow(query string, args ...any) *sql.Row {
	return tx.tx.QueryRow(query, args...)
}

func (tx *Tx) Query(query string, args ...any) (*sql.Rows, error) {
	return tx.tx.Query(query, args...)
}

func (tx *Tx) Exec(query string, args ...any) (sql.Result, error) {
	return tx.tx.Exec(query, args...)
}
