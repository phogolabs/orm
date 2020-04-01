package orm

import (
	"context"

	"github.com/phogolabs/orm/dialect"
	"github.com/phogolabs/orm/dialect/sql"
	"github.com/phogolabs/prana/sqlexec"
)

// RunTxFunc is a transaction function
type RunTxFunc func(gatekeeper *TxGateway) error

// TxGateway represents a gateway in transaction
type TxGateway struct {
	tx       dialect.Tx
	provider *sqlexec.Provider
	dialect  string
}

// All executes the query and returns a list of entities.
func (g *TxGateway) All(ctx context.Context, q sql.Querier, v interface{}) error {
	return g.exec().All(ctx, q, v)
}

// Only returns the only entity in the query, returns an error if not
// exactly one entity was returned.
func (g *TxGateway) Only(ctx context.Context, q sql.Querier, v interface{}) error {
	return g.exec().Only(ctx, q, v)
}

// First returns the first entity in the query. Returns *NotFoundError
// when no user was found.
func (g *TxGateway) First(ctx context.Context, q sql.Querier, v interface{}) error {
	return g.exec().First(ctx, q, v)
}

// Query executes a query that returns rows, typically a SELECT in SQL.
// It scans the result into the pointer v. In SQL, you it's usually *sql.Rows.
func (g *TxGateway) Query(ctx context.Context, q sql.Querier) (*sql.Rows, error) {
	return g.exec().Query(ctx, q)
}

// Exec executes a query that doesn't return rows. For example, in SQL, INSERT
// or UPDATE.  It scans the result into the pointer v. In SQL, you it's usually
// sql.Result.
func (g *TxGateway) Exec(ctx context.Context, q sql.Querier) (sql.Result, error) {
	return g.exec().Exec(ctx, q)
}

// Commit commits the transaction
func (g *TxGateway) Commit() error {
	return g.tx.Commit()
}

// Rollback rollbacks the transaction
func (g *TxGateway) Rollback() error {
	return g.tx.Rollback()
}

func (g *TxGateway) exec() *ExecGateway {
	execer := &ExecGateway{
		driver:   g.tx,
		dialect:  g.dialect,
		provider: g.provider,
	}

	return execer
}
