package oak

import (
	"database/sql"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx/reflectx"
)

var (
	_ Query      = &Stmt{}
	_ NamedQuery = &NamedStmt{}
)

// Stmt represents a single command from SQL sqlexec.
type Stmt struct {
	query  string
	params []Param
}

// SQL create a new command from raw query
func SQL(query string, params ...Param) Query {
	return &Stmt{
		query:  query,
		params: params,
	}
}

// Query prepares prepares the command for execution.
func (cmd *Stmt) Query() (string, []Param) {
	return cmd.query, cmd.params
}

// NamedStmt is command that can use named parameters
type NamedStmt struct {
	query  string
	params []Param
}

// NamedSQL create a new named command from raw query
func NamedSQL(query string, params ...Param) NamedQuery {
	return &NamedStmt{
		query:  query,
		params: params,
	}
}

// NamedQuery prepares the command for execution.
func (cmd *NamedStmt) NamedQuery() (string, map[string]interface{}) {
	params := make(map[string]interface{})

	for _, param := range cmd.params {
		switch arg := param.(type) {
		case sql.NamedArg:
			params[arg.Name] = arg.Value
		case map[string]interface{}:
			for key, value := range arg {
				params[key] = value
			}
		default:
			for key, value := range cmd.bindArgs(arg) {
				params[key] = value
			}
		}
	}

	return cmd.query, params
}

func (cmd *NamedStmt) bindArgs(param Param) map[string]interface{} {
	params := make(map[string]interface{})
	mapper := reflectx.NewMapper("db")

	v := reflect.ValueOf(param)

	for v = reflect.ValueOf(param); v.Kind() == reflect.Ptr; {
		v = v.Elem()
	}

	for key, value := range mapper.FieldMap(v) {
		key = strings.ToLower(key)
		params[key] = value.Interface()
	}

	return params
}
