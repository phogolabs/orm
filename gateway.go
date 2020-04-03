package orm

import (
	"context"
	"fmt"
	"io"
	"time"

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
	driver   *sql.Driver
	provider *sqlexec.Provider
}

// Connect creates a new gateway connecto to the provided URL.
func Connect(url string) (*Gateway, error) {
	driver, source, err := prana.ParseURL(url)
	if err != nil {
		return nil, err
	}

	gateway, err := Open(driver, source)
	if err != nil {
		return nil, err
	}

	if err = gateway.Ping(); err != nil {
		return nil, err
	}

	return gateway, nil
}

// Open creates a new gateway connected to the provided source.
func Open(driver, source string) (*Gateway, error) {
	var (
		drv *sql.Driver
		err error
	)

	switch driver {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		if drv, err = sql.Open(driver, source); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("orm: unsupported driver: %q", driver)
	}

	drv.DB().SetMaxIdleConns(32)
	drv.DB().SetMaxOpenConns(32)

	return &Gateway{
		provider: &sqlexec.Provider{DriverName: driver},
		driver:   drv,
	}, nil
}

// Ping pins the underlying database
func (g *Gateway) Ping() error {
	return g.driver.DB().Ping()
}

// Close closes the connection to the  database.
func (g *Gateway) Close() error {
	return g.driver.Close()
}

// Dialect returns the driver's dialect
func (g *Gateway) Dialect() string {
	return g.driver.Dialect()
}

// Debug sets the debug logging
func (g *Gateway) Debug() *ExecGateway {
	execer := &ExecGateway{
		driver:   dialect.Debug(g.driver),
		dialect:  g.driver.Dialect(),
		provider: g.provider,
	}

	return execer
}

// SetMaxIdleConns sets the maximum number of connections in the idle
// connection pool.
//
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns,
// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit.
//
// If n <= 0, no idle connections are retained.
//
// The default max idle connections is currently 2. This may change in
// a future release.
func (g *Gateway) SetMaxIdleConns(value int) {
	g.driver.DB().SetMaxIdleConns(value)
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
//
// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
// MaxIdleConns, then MaxIdleConns will be reduced to match the new
// MaxOpenConns limit.
//
// If n <= 0, then there is no limit on the number of open connections.
// The default is 0 (unlimited).
func (g *Gateway) SetMaxOpenConns(value int) {
	g.driver.DB().SetMaxOpenConns(value)
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
//
// Expired connections may be closed lazily before reuse.
//
// If d <= 0, connections are reused forever.
func (g *Gateway) SetConnMaxLifetime(duration time.Duration) {
	g.driver.DB().SetConnMaxLifetime(duration)
}

// Migrate runs all pending migration
func (g *Gateway) Migrate(fileSystem FileSystem) error {
	db := sqlx.NewDb(g.driver.DB(), g.driver.Dialect())
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

// BeginTx begins a transaction and returns an *Tx
func (g *Gateway) BeginTx(ctx context.Context, opts *sql.TxOptions) (*TxGateway, error) {
	tx, err := g.driver.BeginTx(ctx, opts)
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
