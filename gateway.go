package orm

import (
	"context"
	"fmt"
	"io"

	"github.com/jmoiron/sqlx"
	"github.com/phogolabs/orm/dialect"
	"github.com/phogolabs/orm/dialect/sql"
	"github.com/phogolabs/prana"
	"github.com/phogolabs/prana/sqlexec"
	"github.com/phogolabs/prana/sqlmigr"
)

var _ GatewayQuerier = &Gateway{}

// Gateway is connected to a database and can executes SQL queries against it.
type Gateway struct {
	db       *sql.DB
	provider *sqlexec.Provider
	driver   dialect.Driver
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
	var (
		driver *sql.Driver
		err    error
	)

	switch name {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		if driver, err = sql.Open(name, source); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("orm: unsupported driver: %q", name)
	}

	gateway := &Gateway{
		provider: &sqlexec.Provider{DriverName: name},
		driver:   driver,
		db:       driver.DB(),
	}

	for _, opt := range opts {
		opt.Apply(gateway)
	}

	return gateway, nil
}

// Ping pins the underlying database
func (g *Gateway) Ping(ctx context.Context) error {
	return g.db.PingContext(ctx)
}

// Close closes the connection to the  database.
func (g *Gateway) Close() error {
	return g.driver.Close()
}

// Dialect returns the driver's dialect
func (g *Gateway) Dialect() string {
	return g.driver.Dialect()
}

// Migrate runs all pending migration
func (g *Gateway) Migrate(fileSystem FileSystem) error {
	db := sqlx.NewDb(g.db, g.driver.Dialect())
	return sqlmigr.RunAll(db, fileSystem)
}

// ReadDir loads all script commands from a given directory. Note that all
// scripts should have .sql extension and support the database driver.
func (g *Gateway) ReadDir(fs FileSystem) error {
	return g.provider.ReadDir(fs)
}

// ReadFrom loads all script commands from a given directory. Note that all
// scripts should have .sql extension.
func (g *Gateway) ReadFrom(r io.Reader) (int64, error) {
	return g.provider.ReadFrom(r)
}

// Begin begins a transaction and returns an *Tx
func (g *Gateway) Begin(ctx context.Context) (*TxGateway, error) {
	tx, err := g.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}

	return g.begin(tx), nil
}

// RunInTx runs a callback function within a transaction. It commits the
// transaction if succeeds, otherwise rollbacks.
func (g *Gateway) RunInTx(ctx context.Context, fn RunTxFunc) error {
	tx, err := g.driver.Tx(ctx)
	if err != nil {
		return err
	}

	txg := g.begin(tx)

	if err := fn(txg); err != nil {
		txg.Rollback()
		return err
	}

	return txg.Commit()
}

// All executes the query and returns a list of entities.
func (g *Gateway) All(ctx context.Context, q sql.Querier, v interface{}) error {
	return g.exec().All(ctx, q, v)
}

// Only returns the only entity in the query, returns an error if not
// exactly one entity was returned.
func (g *Gateway) Only(ctx context.Context, q sql.Querier, v interface{}) error {
	return g.exec().Only(ctx, q, v)
}

// First returns the first entity in the query. Returns *NotFoundError
// when no records were found.
func (g *Gateway) First(ctx context.Context, q sql.Querier, v interface{}) error {
	return g.exec().First(ctx, q, v)
}

// Query executes a query that returns rows, typically a SELECT in SQL.
// It scans the result into the pointer v. In SQL, you it's usually *sql.Rows.
func (g *Gateway) Query(ctx context.Context, q sql.Querier) (*sql.Rows, error) {
	return g.exec().Query(ctx, q)
}

// Exec executes a query that doesn't return rows. For example, in SQL, INSERT
// or UPDATE.  It scans the result into the pointer v. In SQL, you it's usually
// sql.Result.
func (g *Gateway) Exec(ctx context.Context, q sql.Querier) (sql.Result, error) {
	return g.exec().Exec(ctx, q)
}

func (g *Gateway) exec() *ExecGateway {
	execer := &ExecGateway{
		driver:   g.driver,
		dialect:  g.driver.Dialect(),
		provider: g.provider,
	}

	return execer
}

func (g *Gateway) begin(tx dialect.Tx) *TxGateway {
	execer := &TxGateway{
		tx:       tx,
		provider: g.provider,
		dialect:  g.driver.Dialect(),
	}

	return execer
}
