package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

var (
	_ NamedQuery = &Stmt{}
)

// Stmt represents a single command from SQL sqlexec.
type Stmt struct {
	routine string
	query   string
	params  []Param
}

// SQL create a new command from raw query
func SQL(query string, params ...Param) NamedQuery {
	return &Stmt{
		query:  query,
		params: params,
	}
}

// Routine create a new routine for given name
func Routine(routine string, params ...Param) NamedQuery {
	return &Stmt{
		routine: routine,
		params:  params,
	}
}

// NamedQuery prepares prepares the command for execution.
func (cmd *Stmt) NamedQuery() (string, map[string]Param) {
	return cmd.prepareQuery(), cmd.prepareParams()
}

func (cmd *Stmt) prepareQuery() string {
	return sqlx.Rebind(sqlx.NAMED, cmd.query)
}

func (cmd *Stmt) prepareParams() map[string]Param {
	params := make(map[string]interface{})
	index := 1

	for _, param := range cmd.params {
		if mapper, ok := param.(Mapper); ok {
			param = mapper.Map()
		}

		switch arg := param.(type) {
		case sql.NamedArg:
			params[arg.Name] = arg.Value
		case map[string]interface{}:
			for key, value := range arg {
				params[key] = value
			}
		default:
			for key, value := range cmd.bind(arg) {
				if key == "?" {
					key = fmt.Sprintf("arg%d", index)
					index++
				}
				params[key] = value
			}
		}
	}

	return params
}

func (cmd *Stmt) bind(param Param) map[string]interface{} {
	value := reflect.ValueOf(param)

	for value = reflect.ValueOf(param); value.Kind() == reflect.Ptr; {
		value = value.Elem()
	}

	switch {
	case value.Kind() == reflect.Struct:
		return cmd.reflect(value)
	default:
		return map[string]interface{}{
			"?": reflect.ValueOf(param).Interface(),
		}
	}
}

func (cmd *Stmt) reflect(v reflect.Value) map[string]interface{} {
	params := make(map[string]interface{})
	mapper := reflectx.NewMapper("db")

	for key, value := range mapper.FieldMap(v) {
		key = strings.ToLower(key)
		params[key] = value.Interface()
	}

	return params
}
