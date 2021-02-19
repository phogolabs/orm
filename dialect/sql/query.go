package sql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/phogolabs/orm/dialect/sql/scan"
)

var _ Procedure = &RoutineQuery{}

// RoutineQuery represents a named routine
type RoutineQuery struct {
	name    string
	dialect string
	args    []interface{}
	stmt    *NamedQuery
}

// Routine create a new routine for given name
func Routine(name string, args ...interface{}) *RoutineQuery {
	return &RoutineQuery{
		name: name,
		args: args,
	}
}

// Name returns the name of the procedure
func (r *RoutineQuery) Name() string {
	return r.name
}

// Total returns the total count of parameters
func (r *RoutineQuery) Total() int {
	if r.stmt == nil {
		return len(r.args)
	}
	return r.stmt.Total()
}

// Dialect returns the dialect
func (r *RoutineQuery) Dialect() string {
	return r.dialect
}

// SetDialect sets the dialect
func (r *RoutineQuery) SetDialect(dialect string) {
	r.dialect = dialect
}

// SetQuery sets the query
func (r *RoutineQuery) SetQuery(value string) {
	r.stmt = Query(value, r.args...)
	r.stmt.SetDialect(r.dialect)
}

// Query returns the query representation of the element
// and its arguments (if any).
func (r *RoutineQuery) Query() (string, []interface{}) {
	if r.stmt == nil {
		return "", r.args
	}
	// return the underlying information
	return r.stmt.Query()
}

// Error returns the underlying error
func (r *RoutineQuery) Error() error {
	if r.stmt == nil {
		return nil
	}
	// return the underlying error
	return r.stmt.err
}

// A NamedArg is a named argument. NamedArg values may be used as
// arguments to Query or Exec and bind to the corresponding named
// parameter in the SQL statement.
//
// For a more concise way to create NamedArg values, see
// the Named function.
type NamedArg = sql.NamedArg

var _ Statement = &NamedQuery{}

// NameQuery is a named query that uses named arguments
type NamedQuery struct {
	err     error
	dialect string
	query   string
	args    []sql.NamedArg
}

// NamedQuery create a new named query
func Query(query string, params ...interface{}) *NamedQuery {
	query, columns := scan.NamedQuery(query)
	// scane the arguments
	args, err := scan.Args(params, columns...)
	// create the query
	querier := &NamedQuery{
		err:   err,
		query: query,
	}

	for index, name := range columns {
		param := NamedArg{
			Name:  name,
			Value: args[index],
		}
		querier.args = append(querier.args, param)
	}

	return querier
}

// Dialect returns the dialect
func (r *NamedQuery) Dialect() string {
	return r.dialect
}

// SetDialect sets the dialect
func (r *NamedQuery) SetDialect(dialect string) {
	r.dialect = dialect
}

// Total returns the total count of parameters
func (r *NamedQuery) Total() int {
	return len(r.args)
}

// Error returns the error
func (r *NamedQuery) Error() error {
	return r.err
}

// Query returns the routine
// Query returns the query representation of the element
// and its arguments (if any).
func (r *NamedQuery) Query() (string, []interface{}) {
	var (
		query = r.query
		args  = make([]interface{}, 0)
	)

	for index, param := range r.args {
		target := fmt.Sprintf(":%v", param.Name)

		switch r.dialect {
		case "postgres":
			name := fmt.Sprintf("$%d", index+1)
			query = strings.Replace(query, target, name, 1)
			args = append(args, param.Value)
		case "mysql", "sqlite":
			name := "?"
			query = strings.Replace(query, target, name, 1)
			args = append(args, param.Value)
		case "oci8", "ora", "goracle", "godror":
			args = append(args, param)
		case "sqlserver":
			name := fmt.Sprintf("@%v", param.Name)
			query = strings.Replace(query, target, name, 1)
			args = append(args, param)
		default:
			args = append(args, param)
		}
	}

	return query, args
}
