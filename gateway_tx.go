package oak

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Tx is an sqlx wrapper around sqlx.Tx with extra functionality
type Tx struct {
	tx *sqlx.Tx
}

// Commit commits the transaction.
func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}

// Rollback aborts the transaction.
func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}

// Select executes a given query and maps the result to the provided slice of entities.
func (tx *Tx) Select(dest Entity, query Query) error {
	return selectMany(context.Background(), tx.tx, dest, query)
}

// SelectContext executes a given query and maps the result to the provided slice of entities.
func (tx *Tx) SelectContext(ctx context.Context, dest Entity, query Query) error {
	return selectMany(ctx, tx.tx, dest, query)
}

// SelectOne executes a given query and maps a single result to the provided entity.
func (tx *Tx) SelectOne(dest Entity, query Query) error {
	return selectOne(context.Background(), tx.tx, dest, query)
}

// SelectOneContext executes a given query and maps a single result to the provided entity.
func (tx *Tx) SelectOneContext(ctx context.Context, dest Entity, query Query) error {
	return selectOne(context.Background(), tx.tx, dest, query)
}

// Query executes a given query and returns an instance of rows cursor.
func (tx *Tx) Query(query Query) (*Rows, error) {
	return queryRows(context.Background(), tx.tx, query)
}

// QueryContext executes a given query and returns an instance of rows cursor.
func (tx *Tx) QueryContext(ctx context.Context, query Query) (*Rows, error) {
	return queryRows(ctx, tx.tx, query)
}

// QueryRow executes a given query and returns an instance of row.
func (tx *Tx) QueryRow(query Query) (*Row, error) {
	return queryRow(context.Background(), tx.tx, query)
}

// QueryRowContext executes a given query and returns an instance of row.
func (tx *Tx) QueryRowContext(ctx context.Context, query Query) (*Row, error) {
	return queryRow(ctx, tx.tx, query)
}

// Exec executes a given query. It returns a result that provides information
// about the affected rows.
func (tx *Tx) Exec(query Query) (Result, error) {
	return exec(context.Background(), tx.tx, query)
}

// ExecContext executes a given query. It returns a result that provides information
// about the affected rows.
func (tx *Tx) ExecContext(ctx context.Context, query Query) (Result, error) {
	return exec(ctx, tx.tx, query)
}
