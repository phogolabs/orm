package sql

import (
	"bytes"
	"fmt"
	"strings"
)

var _ Querier = &routine{}

// routine represents a named routine
type routine struct {
	name   string
	params []interface{}
}

// Routine create a new routine for given name
func Routine(name string, params ...interface{}) Querier {
	return &routine{
		name:   name,
		params: params,
	}
}

// Name returns the name of the routine
func (r *routine) Name() string {
	return r.name
}

// Query returns the routine
// Query returns the query representation of the element
// and its arguments (if any).
func (r *routine) Query() (string, []interface{}) {
	return r.name, r.params
}

var _ Querier = &command{}

type command struct {
	query   string
	dialect string
	params  []interface{}
	err     error
}

// Command create a new command from raw query
func Command(query string, params ...interface{}) Querier {
	stmt := &command{
		query:  query,
		params: params,
	}

	stmt.rename()
	return stmt
}

// SetDialect sets the dialect
func (r *command) SetDialect(dialect string) {
	r.dialect = dialect
}

// Query returns the routine
// Query returns the query representation of the element
// and its arguments (if any).
func (r *command) Query() (string, []interface{}) {
	var (
		query, names = r.compile()
		params       = r.bind(names)
	)

	return query, params
}

func (r *command) rename() {
	var (
		query  = r.query
		buffer = &bytes.Buffer{}
		next   int
	)

	for index := strings.Index(query, "?"); index != -1; index = strings.Index(query, "?") {
		fmt.Fprintf(buffer, query[:index])
		fmt.Fprintf(buffer, ":arg%d", next)

		query = query[index+1:]
		next++
	}

	fmt.Fprintf(buffer, query)
	r.query = buffer.String()
}

func (r *command) compile() (query string, columns []string) {
	query, columns, err := compile(r.dialect, r.query)
	if err != nil {
		r.err = err
	}

	return query, columns
}

func (r *command) bind(columns []string) []interface{} {
	if len(columns) == 0 {
		return r.params
	}

	var (
		kv     = bindParam(r.params)
		params = []interface{}{}
	)

	for _, name := range columns {
		if v, ok := kv[name]; ok {
			params = append(params, v)
		}
	}

	return params
}
