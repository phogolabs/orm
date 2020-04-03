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

	return scan.Rows(rows, v)
}

// Only returns the only entity in the query, returns an error if not
// exactly one entity was returned.
func (g *ExecGateway) Only(ctx context.Context, q sql.Querier, v interface{}) error {
	var (
		value     = reflect.New(reflect.SliceOf(reflect.TypeOf(v)))
		rows, err = g.Query(ctx, q)
	)

	if err != nil {
		return err
	}

	if err := scan.Rows(rows, value.Interface()); err != nil {
		return err
	}

	switch value.Elem().Len() {
	case 1:
		source := reflect.Indirect(value.Elem().Index(0))
		target := reflect.Indirect(reflect.ValueOf(v))
		target.Set(source)
	case 0:
		return &NotFoundError{nameOf(value.Type())}
	default:
		return &NotSingularError{nameOf(value.Type())}
	}

	return nil
}

// First returns the first entity in the query. Returns *NotFoundError
// when no user was found.
func (g *ExecGateway) First(ctx context.Context, q sql.Querier, v interface{}) error {
	var (
		value     = reflect.New(reflect.SliceOf(reflect.TypeOf(v)))
		rows, err = g.Query(ctx, q)
	)

	if err != nil {
		return err
	}

	if err := scan.Rows(rows, value.Interface()); err != nil {
		return err
	}

	if value := value.Elem(); value.Len() == 0 {
		return &NotFoundError{nameOf(value.Type())}
	}

	source := reflect.Indirect(value.Elem().Index(0))
	target := reflect.Indirect(reflect.ValueOf(v))
	target.Set(source)
	return nil
}

// Query executes a query that returns rows, typically a SELECT in SQL.
// It scans the result into the pointer v. In SQL, you it's usually *sql.Rows.
func (g *ExecGateway) Query(ctx context.Context, q sql.Querier) (*sql.Rows, error) {
	query, params, err := g.compile(q)

	if err != nil {
		return nil, err
	}

	var rows = &sql.Rows{}

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
		return nil, err
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
			return "", nil, err
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
			return "", nil, err
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
			"Error 1062",               // MySQL 1062 error (ER_DUP_ENTRY).
			"UNIQUE constraint failed", // SQLite.
			"duplicate key value violates unique constraint", // PostgreSQL.
		}
	)

	for index := range errors {
		if strings.Contains(msg, errors[index]) {
			return &ConstraintError{msg, err}
		}
	}

	return err
}
