package orm

import (
	"context"

	"github.com/phogolabs/orm/dialect"
	"github.com/phogolabs/orm/dialect/sql"
)

var _ Querier = &GatewayTx{}

// RunTxFunc is a transaction function
type RunTxFunc func(gateway *GatewayTx) error

// GatewayTx represents a gateway in transaction
type GatewayTx struct {
	engine *engine
}

// All executes the query and returns a list of entities.
func (g *GatewayTx) All(ctx context.Context, q sql.Querier, v interface{}) error {
	return g.engine.All(ctx, q, v)
}

// Only returns the only entity in the query, returns an error if not
// exactly one entity was returned.
func (g *GatewayTx) Only(ctx context.Context, q sql.Querier, v interface{}) error {
	return g.engine.Only(ctx, q, v)
}

// First returns the first entity in the query. Returns *NotFoundError
// when no user was found.
func (g *GatewayTx) First(ctx context.Context, q sql.Querier, v interface{}) error {
	return g.engine.First(ctx, q, v)
}

// Query executes a query that returns rows, typically a SELECT in SQL.
// It scans the result into the pointer v. In SQL, you it's usually *sql.Rows.
func (g *GatewayTx) Query(ctx context.Context, q sql.Querier) (*sql.Rows, error) {
	return g.engine.Query(ctx, q)
}

// Exec executes a query that doesn't return rows. For example, in SQL, INSERT
// or UPDATE.  It scans the result into the pointer v. In SQL, you it's usually
// sql.Result.
func (g *GatewayTx) Exec(ctx context.Context, q sql.Querier) (sql.Result, error) {
	return g.engine.Exec(ctx, q)
}

// Commit commits the transaction
func (g *GatewayTx) Commit() error {
	tx := g.engine.querier.(dialect.Tx)
	return tx.Commit()
}

// Rollback rollbacks the transaction
func (g *GatewayTx) Rollback() error {
	tx := g.engine.querier.(dialect.Tx)
	return tx.Rollback()
}
