package sql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/phogolabs/orm/dialect/sql/scan"
)

var _ Querier = &RoutineQuerier{}

// RoutineQuerier represents a named routine
type RoutineQuerier struct {
	name string
	args []interface{}
}

// Routine create a new routine for given name
func Routine(name string, args ...interface{}) *RoutineQuerier {
	return &RoutineQuerier{
		name: name,
		args: args,
	}
}

// Name returns the name of the routine
func (r *RoutineQuerier) Name() string {
	return r.name
}

// Query returns the routine
// Query returns the query representation of the element
// and its arguments (if any).
func (r *RoutineQuerier) Query() (string, []interface{}) {
	return r.name, r.args
}

// A NamedArg is a named argument. NamedArg values may be used as
// arguments to Query or Exec and bind to the corresponding named
// parameter in the SQL statement.
//
// For a more concise way to create NamedArg values, see
// the Named function.
type NamedArg = sql.NamedArg

// Named provides a more concise way to create NamedArg values.
var Named = sql.Named

var _ Querier = &NamedQuerier{}

// A NamedQuerier is a named query that uses named arguments
type NamedQuerier struct {
	query   string
	args    []sql.NamedArg
	dialect string
}

// NamedQuery create a new named query
func NamedQuery(query string, params ...interface{}) *NamedQuerier {
	query, columns := scan.NamedQuery(query)

	args, err := scan.Args(params, columns...)
	if err != nil {
		args = params
	}

	querier := &NamedQuerier{
		query: query,
	}

	for index, name := range columns {
		param := sql.Named(name, args[index])
		querier.args = append(querier.args, param)
	}

	return querier
}

// Dialect returns the dialect
func (r *NamedQuerier) Dialect() string {
	return r.dialect
}

// SetDialect sets the dialect
func (r *NamedQuerier) SetDialect(dialect string) {
	r.dialect = dialect
}

// Total returns the total count of parameters
func (r *NamedQuerier) Total() int {
	return len(r.args)
}

// Query returns the routine
// Query returns the query representation of the element
// and its arguments (if any).
func (r *NamedQuerier) Query() (string, []interface{}) {
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
