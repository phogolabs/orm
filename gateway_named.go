package oak

import (
	"context"
)

// NamedRoutine returns a SQL statement for given name and map the parameters as
// named.
func (g *Gateway) NamedRoutine(name string, params ...Param) (NamedQuery, error) {
	query, err := g.provider.Query(name)
	if err != nil {
		return nil, err
	}

	return NamedSQL(query, params...), nil
}

// NamedSelect executes a given query and maps the result to the provided slice of entities.
func (g *Gateway) NamedSelect(dest Entity, query NamedQuery) error {
	return namedSelectMany(context.Background(), g.db, dest, query)
}

// NamedSelectContext executes a given query and maps the result to the provided slice of entities.
func (g *Gateway) NamedSelectContext(ctx context.Context, dest Entity, query NamedQuery) error {
	return namedSelectMany(ctx, g.db, dest, query)
}

// NamedSelectOne executes a given query and maps a single result to the provided entity.
func (g *Gateway) NamedSelectOne(dest Entity, query NamedQuery) error {
	return namedSelectOne(context.Background(), g.db, dest, query)
}

// NamedSelectOneContext executes a given query and maps a single result to the provided entity.
func (g *Gateway) NamedSelectOneContext(ctx context.Context, dest Entity, query NamedQuery) error {
	return namedSelectOne(ctx, g.db, dest, query)
}

// NamedQuery executes a given query and returns an instance of rows cursor.
func (g *Gateway) NamedQuery(query NamedQuery) (*Rows, error) {
	return namedQueryRows(context.Background(), g.db, query)
}

// NamedQueryContext executes a given query and returns an instance of rows cursor.
func (g *Gateway) NamedQueryContext(ctx context.Context, query NamedQuery) (*Rows, error) {
	return namedQueryRows(ctx, g.db, query)
}

// NamedQueryRow executes a given query and returns an instance of row.
func (g *Gateway) NamedQueryRow(query NamedQuery) (*Row, error) {
	return namedQueryRow(context.Background(), g.db, query)
}

// NamedQueryRowContext executes a given query and returns an instance of row.
func (g *Gateway) NamedQueryRowContext(ctx context.Context, query NamedQuery) (*Row, error) {
	return namedQueryRow(ctx, g.db, query)
}

// NamedExec executes a given query. It returns a result that provides information
// about the affected rows.
func (g *Gateway) NamedExec(query NamedQuery) (Result, error) {
	return namedExec(context.Background(), g.db, query)
}

// NamedExecContext executes a given query. It returns a result that provides information
// about the affected rows.
func (g *Gateway) NamedExecContext(ctx context.Context, query NamedQuery) (Result, error) {
	return namedExec(ctx, g.db, query)
}
