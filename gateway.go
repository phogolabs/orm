package orm

import (
	"context"
	"io"

	"github.com/jmoiron/sqlx"
	"github.com/phogolabs/prana"
	"github.com/phogolabs/prana/sqlexec"
	"github.com/phogolabs/prana/sqlmigr"
)

// Gateway is connected to a database and can executes SQL queries against it.
type Gateway struct {
	db       *sqlx.DB
	provider *sqlexec.Provider
}

// Connect creates a new gateway connecto to the provided URL.
func Connect(url string) (*Gateway, error) {
	driver, source, err := prana.ParseURL(url)
	if err != nil {
		return nil, err
	}

	return Open(driver, source)
}

// Open creates a new gateway connected to the provided source.
func Open(driver, source string) (*Gateway, error) {
	db, err := sqlx.Open(driver, source)
	if err != nil {
		return nil, err
	}

	return &Gateway{
		provider: &sqlexec.Provider{DriverName: driver},
		db:       db,
	}, nil
}

// Close closes the connection to underlying database.
func (g *Gateway) Close() error {
	return g.db.Close()
}

// DriverName returns the driverName passed to the Open function for this DB.
func (g *Gateway) DriverName() string {
	return g.db.DriverName()
}

// Ping pins the underlying database
func (g *Gateway) Ping() error {
	return g.db.Ping()
}

// Migrate runs all pending migration
func (g *Gateway) Migrate(fileSystem FileSystem) error {
	return sqlmigr.RunAll(g.db, fileSystem)
}

// ReadDir loads all script commands from a given directory. Note that all
// scripts should have .sql extension and support the database driver.
func (g *Gateway) ReadDir(fileSystem FileSystem) error {
	return g.provider.ReadDir(fileSystem)
}

// ReadFrom loads all script commands from a given directory. Note that all
// scripts should have .sql extension.
func (g *Gateway) ReadFrom(reader io.Reader) (int64, error) {
	return g.provider.ReadFrom(reader)
}

// Begin begins a transaction and returns an *Tx
func (g *Gateway) Begin() (*Tx, error) {
	tx, err := g.db.Beginx()
	if err != nil {
		return nil, err
	}

	return &Tx{
		tx:       tx,
		provider: g.provider,
	}, nil
}

// Transaction starts a new transaction. It commits the transaction if
// succeeds, otherwise rollbacks.
func (g *Gateway) Transaction(fn TxFunc) error {
	fnCtx := func(ctx context.Context, tx *Tx) error {
		return fn(tx)
	}
	return g.TransactionContext(context.Background(), fnCtx)
}

// TransactionContext starts a new transaction. It commits the transaction if
// succeeds, otherwise rollbacks.
func (g *Gateway) TransactionContext(ctx context.Context, fn TxContextFunc) error {
	tx, err := g.Begin()
	if err != nil {
		return err
	}

	if fErr := fn(ctx, tx); fErr != nil {
		if tErr := tx.Rollback(); tErr != nil {
			return tErr
		}

		return fErr
	}

	return tx.Commit()
}

// Select executes a given query and maps the result to the provided slice of entities.
func (g *Gateway) Select(dest Entity, query NamedQuery) error {
	return namedSelectMany(context.Background(), g.db, g.provider, dest, query)
}

// SelectContext executes a given query and maps the result to the provided slice of entities.
func (g *Gateway) SelectContext(ctx context.Context, dest Entity, query NamedQuery) error {
	return namedSelectMany(ctx, g.db, g.provider, dest, query)
}

// SelectOne executes a given query and maps a single result to the provided entity.
func (g *Gateway) SelectOne(dest Entity, query NamedQuery) error {
	return namedSelectOne(context.Background(), g.db, g.provider, dest, query)
}

// SelectOneContext executes a given query and maps a single result to the provided entity.
func (g *Gateway) SelectOneContext(ctx context.Context, dest Entity, query NamedQuery) error {
	return namedSelectOne(ctx, g.db, g.provider, dest, query)
}

// Query executes a given query and returns an instance of rows cursor.
func (g *Gateway) Query(query NamedQuery) (*Rows, error) {
	return namedQueryRows(context.Background(), g.db, g.provider, query)
}

// QueryContext executes a given query and returns an instance of rows cursor.
func (g *Gateway) QueryContext(ctx context.Context, query NamedQuery) (*Rows, error) {
	return namedQueryRows(ctx, g.db, g.provider, query)
}

// QueryRow executes a given query and returns an instance of row.
func (g *Gateway) QueryRow(query NamedQuery) (*Row, error) {
	return namedQueryRow(context.Background(), g.db, g.provider, query)
}

// QueryRowContext executes a given query and returns an instance of row.
func (g *Gateway) QueryRowContext(ctx context.Context, query NamedQuery) (*Row, error) {
	return namedQueryRow(ctx, g.db, g.provider, query)
}

// Exec executes a given query. It returns a result that provides information
// about the affected rows.
func (g *Gateway) Exec(query NamedQuery) (Result, error) {
	return namedExec(context.Background(), g.db, g.provider, query)
}

// ExecContext executes a given query. It returns a result that provides information
// about the affected rows.
func (g *Gateway) ExecContext(ctx context.Context, query NamedQuery) (Result, error) {
	return namedExec(ctx, g.db, g.provider, query)
}
