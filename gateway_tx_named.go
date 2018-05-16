package oak

import (
	"context"
)

// NamedSelect executes a given query and maps the result to the provided slice of entities.
func (tx *Tx) NamedSelect(dest Entity, query NamedQuery) error {
	return namedSelectMany(context.Background(), tx.tx, dest, query)
}

// NamedSelectContext executes a given query and maps the result to the provided slice of entities.
func (tx *Tx) NamedSelectContext(ctx context.Context, dest Entity, query NamedQuery) error {
	return namedSelectMany(ctx, tx.tx, dest, query)
}

// NamedSelectOne executes a given query and maps a single result to the provided entity.
func (tx *Tx) NamedSelectOne(dest Entity, query NamedQuery) error {
	return namedSelectOne(context.Background(), tx.tx, dest, query)
}

// NamedSelectOneContext executes a given query and maps a single result to the provided entity.
func (tx *Tx) NamedSelectOneContext(ctx context.Context, dest Entity, query NamedQuery) error {
	return namedSelectOne(context.Background(), tx.tx, dest, query)
}

// NamedQuery executes a given query and returns an instance of rows cursor.
func (tx *Tx) NamedQuery(query NamedQuery) (*Rows, error) {
	return namedQueryRows(context.Background(), tx.tx, query)
}

// NamedQueryContext executes a given query and returns an instance of rows cursor.
func (tx *Tx) NamedQueryContext(ctx context.Context, query NamedQuery) (*Rows, error) {
	return namedQueryRows(ctx, tx.tx, query)
}

// NamedQueryRow executes a given query and returns an instance of row.
func (tx *Tx) NamedQueryRow(query NamedQuery) (*Row, error) {
	return namedQueryRow(context.Background(), tx.tx, query)
}

// NamedQueryRowContext executes a given query and returns an instance of row.
func (tx *Tx) NamedQueryRowContext(ctx context.Context, query NamedQuery) (*Row, error) {
	return namedQueryRow(ctx, tx.tx, query)
}

// NamedExec executes a given query. It returns a result that provides information
// about the affected rows.
func (tx *Tx) NamedExec(query NamedQuery) (Result, error) {
	return namedExec(context.Background(), tx.tx, query)
}

// NamedExecContext executes a given query. It returns a result that provides information
// about the affected rows.
func (tx *Tx) NamedExecContext(ctx context.Context, query NamedQuery) (Result, error) {
	return namedExec(ctx, tx.tx, query)
}
