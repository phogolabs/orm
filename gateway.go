package orm

import (
	"context"

	"github.com/phogolabs/log"
	"github.com/phogolabs/orm/dialect"
	"github.com/phogolabs/orm/dialect/sql"
	"github.com/phogolabs/prana"
)

var _ Querier = &Gateway{}

// Gateway is connected to a database and can executes SQL queries against it.
type Gateway struct {
	engine *engine
}

// Connect creates a new gateway connecto to the provided URL.
func Connect(url string, opts ...Option) (*Gateway, error) {
	driver, source, err := prana.ParseURL(url)
	if err != nil {
		return nil, err
	}

	gateway, err := Open(driver, source, opts...)
	if err != nil {
		return nil, err
	}

	if err = gateway.Ping(context.TODO()); err != nil {
		return nil, err
	}

	return gateway, nil
}

// Open creates a new gateway connected to the provided source.
func Open(name, source string, opts ...Option) (*Gateway, error) {
	driver, err := sql.Open(name, source)
	if err != nil {
		return nil, err
	}

	dialect := driver.Dialect()
	// setup the provider
	provider := &sql.Provider{}
	provider.SetDialect(dialect)

	gateway := &Gateway{
		engine: &engine{
			querier:  driver,
			dialect:  dialect,
			provider: provider,
		},
	}

	for _, opt := range opts {
		if err := opt.Apply(gateway); err != nil {
			return nil, err
		}
	}

	return gateway, nil
}

// Ping pins the underlying database
func (g *Gateway) Ping(ctx context.Context) error {
	driver := g.engine.querier.(dialect.Driver)
	// make a ping request
	return driver.Ping(ctx)
}

// Close closes the connection to the  database.
func (g *Gateway) Close() error {
	driver := g.engine.querier.(dialect.Driver)
	// close the connection
	return driver.Close()
}

// Dialect returns the driver's dialect
func (g *Gateway) Dialect() string {
	return g.engine.dialect
}

// Migrate runs all pending migration
func (g *Gateway) Migrate(storage FileSystem) error {
	driver := g.engine.querier.(dialect.Driver)
	// run the migration
	return driver.Migrate(storage)
}

// Begin begins a transaction and returns an *Tx
func (g *Gateway) Begin(ctx context.Context) (*GatewayTx, error) {
	driver := g.engine.querier.(dialect.Driver)

	tx, err := driver.Tx(ctx)
	if err != nil {
		return nil, err
	}

	gtx := &GatewayTx{
		engine: &engine{
			querier:  tx,
			dialect:  g.engine.dialect,
			provider: g.engine.provider,
		},
	}

	return gtx, nil
}

// RunInTx runs a callback function within a transaction. It commits the
// transaction if succeeds, otherwise rollbacks.
func (g *Gateway) RunInTx(ctx context.Context, fn RunTxFunc) error {
	gtx, err := g.Begin(ctx)
	if err != nil {
		return err
	}

	if err := fn(gtx); err != nil {
		if err := gtx.Rollback(); err != nil {
			log.WithError(err).Error("cannot rollback")
		}
		return err
	}

	return gtx.Commit()
}

// All executes the query and returns a list of entities.
func (g *Gateway) All(ctx context.Context, q sql.Querier, v interface{}) error {
	return g.engine.All(ctx, q, v)
}

// Only returns the only entity in the query, returns an error if not
// exactly one entity was returned.
func (g *Gateway) Only(ctx context.Context, q sql.Querier, v interface{}) error {
	return g.engine.Only(ctx, q, v)
}

// First returns the first entity in the query. Returns *NotFoundError
// when no records were found.
func (g *Gateway) First(ctx context.Context, q sql.Querier, v interface{}) error {
	return g.engine.First(ctx, q, v)
}

// Query executes a query that returns rows, typically a SELECT in SQL.
// It scans the result into the pointer v. In SQL, you it's usually *sql.Rows.
func (g *Gateway) Query(ctx context.Context, q sql.Querier) (*sql.Rows, error) {
	return g.engine.Query(ctx, q)
}

// Exec executes a query that doesn't return rows. For example, in SQL, INSERT
// or UPDATE.  It scans the result into the pointer v. In SQL, you it's usually
// sql.Result.
func (g *Gateway) Exec(ctx context.Context, q sql.Querier) (sql.Result, error) {
	return g.engine.Exec(ctx, q)
}
