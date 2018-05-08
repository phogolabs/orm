package oak

import (
	"context"
	"io"

	"github.com/jmoiron/sqlx"
	"github.com/phogolabs/prana/sqlexec"
)

// Gateway is connected to a database and can executes SQL queries against it.
type Gateway struct {
	db       *sqlx.DB
	provider *sqlexec.Provider
}

// OpenURL creates a new gateway connecto to the provided URL.
func OpenURL(url string) (*Gateway, error) {
	driver, source, err := ParseURL(url)
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

// LoadRoutinesFromReader loads all script commands from a given directory. Note that all
// scripts should have .sql extension.
func (g *Gateway) LoadRoutinesFromReader(reader io.Reader) error {
	_, err := g.provider.ReadFrom(reader)
	return err
}

// LoadRoutinesFrom loads all script commands from a given directory. Note that all
// scripts should have .sql extension.
func (g *Gateway) LoadRoutinesFrom(fileSystem FileSystem) error {
	return g.provider.ReadDir(fileSystem)
}

// Routine returns a SQL statement for given name and parameters.
func (g *Gateway) Routine(name string, params ...sqlexec.Param) (Query, error) {
	return g.provider.Query(name, params...)
}

// NamedRoutine returns a SQL statement for given name and map the parameters as
// named.
func (g *Gateway) NamedRoutine(name string, param sqlexec.Param) (Query, error) {
	return g.provider.NamedQuery(name, param)
}

// Transaction starts a new transaction. It commits the transaction if
// succeeds, otherwise rollbacks.
func (g *Gateway) Transaction(fn TxFunc) error {
	tx, err := g.Begin()
	if err != nil {
		return err
	}

	if fErr := fn(tx); fErr != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}

// Begin begins a transaction and returns an *Tx
func (g *Gateway) Begin() (*Tx, error) {
	tx, err := g.db.Beginx()
	if err != nil {
		return nil, err
	}

	return &Tx{tx: tx}, nil
}

// Select executes a given query and maps the result to the provided slice of entities.
func (g *Gateway) Select(dest Entity, query Query) error {
	return selectMany(context.Background(), g.db, dest, query)
}

// SelectContext executes a given query and maps the result to the provided slice of entities.
func (g *Gateway) SelectContext(ctx context.Context, dest Entity, query Query) error {
	return selectMany(ctx, g.db, dest, query)
}

// SelectOne executes a given query and maps a single result to the provided entity.
func (g *Gateway) SelectOne(dest Entity, query Query) error {
	return selectOne(context.Background(), g.db, dest, query)
}

// SelectOneContext executes a given query and maps a single result to the provided entity.
func (g *Gateway) SelectOneContext(ctx context.Context, dest Entity, query Query) error {
	return selectOne(ctx, g.db, dest, query)
}

// Query executes a given query and returns an instance of rows cursor.
func (g *Gateway) Query(query Query) (*Rows, error) {
	return queryRows(context.Background(), g.db, query)
}

// QueryContext executes a given query and returns an instance of rows cursor.
func (g *Gateway) QueryContext(ctx context.Context, query Query) (*Rows, error) {
	return queryRows(ctx, g.db, query)
}

// QueryRow executes a given query and returns an instance of row.
func (g *Gateway) QueryRow(query Query) (*Row, error) {
	return queryRow(context.Background(), g.db, query)
}

// QueryRowContext executes a given query and returns an instance of row.
func (g *Gateway) QueryRowContext(ctx context.Context, query Query) (*Row, error) {
	return queryRow(ctx, g.db, query)
}

// Exec executes a given query. It returns a result that provides information
// about the affected rows.
func (g *Gateway) Exec(query Query) (Result, error) {
	return exec(context.Background(), g.db, query)
}

// ExecContext executes a given query. It returns a result that provides information
// about the affected rows.
func (g *Gateway) ExecContext(ctx context.Context, query Query) (Result, error) {
	return exec(ctx, g.db, query)
}
