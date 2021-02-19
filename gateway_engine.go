package orm

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/phogolabs/orm/dialect"
	"github.com/phogolabs/orm/dialect/sql"
	"github.com/phogolabs/orm/dialect/sql/scan"
	"github.com/phogolabs/prana/sqlexec"
)

var _ Querier = &engine{}

type engine struct {
	provider *sqlexec.Provider
	querier  dialect.ExecQuerier
	dialect  string
}

// All executes the query and returns a list of entities.
func (g *engine) All(ctx context.Context, q sql.Statement, v interface{}) error {
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
func (g *engine) Only(ctx context.Context, q sql.Statement, v interface{}) error {
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
func (g *engine) First(ctx context.Context, q sql.Statement, v interface{}) error {
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
func (g *engine) Query(ctx context.Context, q sql.Statement) (*sql.Rows, error) {
	query, params, err := g.compile(q)
	if err != nil {
		return nil, g.wrap(err)
	}

	rows := &sql.Rows{}

	if err := g.querier.Query(ctx, query, params, rows); err != nil {
		return nil, g.wrap(err)
	}

	return rows, nil
}

// Exec executes a query that doesn't return rows. For example, in SQL, INSERT
// or UPDATE.  It scans the result into the pointer v. In SQL, you it's usually
// sql.Result.
func (g *engine) Exec(ctx context.Context, q sql.Statement) (sql.Result, error) {
	query, params, err := g.compile(q)
	if err != nil {
		return nil, g.wrap(err)
	}

	var result sql.Result

	if err := g.querier.Exec(ctx, query, params, &result); err != nil {
		return nil, g.wrap(err)
	}

	return result, nil
}

func (g *engine) compile(stmt sql.Statement) (string, []interface{}, error) {
	// find the command if any
	if routine, ok := stmt.(sql.Procedure); ok {
		// get the actual SQL query
		query, err := g.provider.Query(routine.Name())
		// if getting the query fails
		if err != nil {
			return "", nil, g.wrap(err)
		}
		// sets the routine's query
		routine.SetQuery(query)
	}

	// set the dialect
	if query, ok := stmt.(sql.Translatable); ok {
		query.SetDialect(g.dialect)
	}

	// compile the query
	query, params := stmt.Query()

	// check for errors
	if reporter, ok := stmt.(sql.Errorable); ok {
		if err := reporter.Error(); err != nil {
			return "", nil, g.wrap(err)
		}
	}

	return query, params, nil
}

func (g *engine) wrap(err error) error {
	if err == nil {
		return nil
	}

	var (
		// error as string
		errm = err.Error()
		// known errors
		errors = [...]string{
			// MySQL 1062 error (ER_DUP_ENTRY).
			"Error 1062",
			// SQLite.
			"UNIQUE constraint failed: %s",
			// PostgreSQL.
			"pq: duplicate key value violates unique constraint %q",
			// check constraint
			"pq: violates check constraint %q",
			// new row check constraint
			"pq: new row for relation %q violates check constraint %q",
		}
	)

	for _, message := range errors {
		// name of the constrain
		var name string
		// scane the name
		fmt.Sscanf(errm, message, &name, &name)
		// check
		if len(name) > 0 || strings.Contains(errm, message) {
			// return the constraint
			return &ConstraintError{
				name: name,
				wrap: err,
			}
		}
	}

	return err
}

func nameOf(value reflect.Type) string {
	switch value.Kind() {
	case reflect.Ptr:
		return nameOf(value.Elem())
	case reflect.Struct:
		return strings.ToLower(value.Name())
	case reflect.Slice:
		return nameOf(value.Elem())
	case reflect.Array:
		return nameOf(value.Elem())
	case reflect.Map:
		return "map"
	default:
		return value.Name()
	}
}
