package orm

import (
	"context"
	"reflect"
	"strings"

	"github.com/phogolabs/orm/dialect"
	"github.com/phogolabs/orm/dialect/sql"
	"github.com/phogolabs/orm/dialect/sql/scan"
	"github.com/phogolabs/prana/sqlexec"
)

var _ GatewayQuerier = &ExecGateway{}

// ExecGateway is connected to a database and can executes SQL queries against it.
type ExecGateway struct {
	driver   dialect.ExecQuerier
	provider *sqlexec.Provider
	dialect  string
}

// All executes the query and returns a list of entities.
func (g *ExecGateway) All(ctx context.Context, q sql.Querier, v interface{}) error {
	rows, err := g.Query(ctx, q)
	if err != nil {
		return err
	}
	defer rows.Close()

	err = scan.Rows(rows, v)
	return g.wrap(err)
}

// Only returns the only entity in the query, returns an error if not
// exactly one entity was returned.
func (g *ExecGateway) Only(ctx context.Context, q sql.Querier, v interface{}) error {
	rows, err := g.Query(ctx, q)
	if err != nil {
		return err
	}
	defer rows.Close()

	err = scan.Row(rows, v)

	switch {
	case err == sql.ErrNoRows:
		return &NotFoundError{nameOf(reflect.TypeOf(v))}
	case err == scan.ErrOneRow:
		return &NotSingularError{nameOf(reflect.TypeOf(v))}
	default:
		return g.wrap(err)
	}
}

// First returns the first entity in the query. Returns *NotFoundError
// when no user was found.
func (g *ExecGateway) First(ctx context.Context, q sql.Querier, v interface{}) error {
	rows, err := g.Query(ctx, q)
	if err != nil {
		return g.wrap(err)
	}
	defer rows.Close()

	err = scan.Row(rows, v)

	switch {
	case err == sql.ErrNoRows:
		return &NotFoundError{nameOf(reflect.TypeOf(v))}
	case err == scan.ErrOneRow:
		return nil
	default:
		return g.wrap(err)
	}
}

// Query executes a query that returns rows, typically a SELECT in SQL.
// It scans the result into the pointer v. In SQL, you it's usually *sql.Rows.
func (g *ExecGateway) Query(ctx context.Context, q sql.Querier) (*sql.Rows, error) {
	query, params, err := g.compile(q)
	if err != nil {
		return nil, g.wrap(err)
	}

	rows := &sql.Rows{}

	if err := g.driver.Query(ctx, query, params, rows); err != nil {
		return nil, g.wrap(err)
	}

	return rows, nil
}

// Exec executes a query that doesn't return rows. For example, in SQL, INSERT
// or UPDATE.  It scans the result into the pointer v. In SQL, you it's usually
// sql.Result.
func (g *ExecGateway) Exec(ctx context.Context, q sql.Querier) (sql.Result, error) {
	query, params, err := g.compile(q)
	if err != nil {
		return nil, g.wrap(err)
	}

	var result sql.Result

	if err := g.driver.Exec(ctx, query, params, &result); err != nil {
		return nil, g.wrap(err)
	}

	return result, nil
}

func (g *ExecGateway) compile(querier sql.Querier) (string, []interface{}, error) {
	// find the command if any
	if _, ok := querier.(commandable); ok {
		name, params := querier.Query()

		query, err := g.provider.Query(name)
		if err != nil {
			return "", nil, g.wrap(err)
		}

		querier = sql.NamedQuery(query, params...)
	}

	// set the dialect
	if syntax, ok := querier.(translatable); ok {
		syntax.SetDialect(g.dialect)
	}

	query, params := querier.Query()

	// check for errors
	if reporter, ok := querier.(errorable); ok {
		if err := reporter.Err(); err != nil {
			return "", nil, g.wrap(err)
		}
	}

	return query, params, nil
}

func (g *ExecGateway) wrap(err error) error {
	if err == nil {
		return nil
	}

	var (
		msg = err.Error()
		// error format per dialect.
		errors = [...]string{
			// MySQL 1062 error (ER_DUP_ENTRY).
			"Error 1062",
			// SQLite.
			"UNIQUE constraint failed",
			// PostgreSQL.
			"duplicate key value violates unique constraint",
			"violates check constraint",
		}
	)

	for index := range errors {
		if strings.Contains(msg, errors[index]) {
			return &ConstraintError{msg, err}
		}
	}

	return err
}
