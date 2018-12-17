package orm

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/phogolabs/prana/sqlexec"
)

// Tx is an sqlx wrapper around sqlx.Tx with extra functionality
type Tx struct {
	provider *sqlexec.Provider
	tx       *sqlx.Tx
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
func (tx *Tx) Select(dest Entity, query NamedQuery) error {
	return namedSelectMany(context.Background(), tx.tx, tx.provider, dest, query)
}

// SelectContext executes a given query and maps the result to the provided slice of entities.
func (tx *Tx) SelectContext(ctx context.Context, dest Entity, query NamedQuery) error {
	return namedSelectMany(ctx, tx.tx, tx.provider, dest, query)
}

// SelectOne executes a given query and maps a single result to the provided entity.
func (tx *Tx) SelectOne(dest Entity, query NamedQuery) error {
	return namedSelectOne(context.Background(), tx.tx, tx.provider, dest, query)
}

// SelectOneContext executes a given query and maps a single result to the provided entity.
func (tx *Tx) SelectOneContext(ctx context.Context, dest Entity, query NamedQuery) error {
	return namedSelectOne(context.Background(), tx.tx, tx.provider, dest, query)
}

// Query executes a given query and returns an instance of rows cursor.
func (tx *Tx) Query(query NamedQuery) (*Rows, error) {
	return namedQueryRows(context.Background(), tx.tx, tx.provider, query)
}

// QueryContext executes a given query and returns an instance of rows cursor.
func (tx *Tx) QueryContext(ctx context.Context, query NamedQuery) (*Rows, error) {
	return namedQueryRows(ctx, tx.tx, tx.provider, query)
}

// QueryRow executes a given query and returns an instance of row.
func (tx *Tx) QueryRow(query NamedQuery) (*Row, error) {
	return namedQueryRow(context.Background(), tx.tx, tx.provider, query)
}

// QueryRowContext executes a given query and returns an instance of row.
func (tx *Tx) QueryRowContext(ctx context.Context, query NamedQuery) (*Row, error) {
	return namedQueryRow(ctx, tx.tx, tx.provider, query)
}

// Exec executes a given query. It returns a result that provides information
// about the affected rows.
func (tx *Tx) Exec(query NamedQuery) (Result, error) {
	return namedExec(context.Background(), tx.tx, tx.provider, query)
}

// ExecContext executes a given query. It returns a result that provides information
// about the affected rows.
func (tx *Tx) ExecContext(ctx context.Context, query NamedQuery) (Result, error) {
	return namedExec(ctx, tx.tx, tx.provider, query)
}
