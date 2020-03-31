package orm

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/a8m/rql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/phogolabs/prana/sqlexec"
)

var mapper = reflectx.NewMapper("db")

var (
	_ NamedQuery = &stmt{}
	_ NamedQuery = &routine{}
	_ NamedQuery = &query{}
)

type stmt struct {
	query  string
	params []Param
}

// SQL create a new command from raw query
func SQL(query string, params ...Param) NamedQuery {
	return &stmt{
		query:  query,
		params: params,
	}
}

// NamedQuery prepares prepares the command for execution.
func (cmd *stmt) NamedQuery() (string, map[string]Param) {
	namedQuery := sqlx.Rebind(sqlx.NAMED, cmd.query)
	namedParams := prepareParams(cmd.params)
	return namedQuery, namedParams
}

type routine struct {
	name   string
	query  string
	params []Param
}

// Routine create a new routine for given name
func Routine(name string, params ...Param) NamedQuery {
	return &routine{
		name:   name,
		params: params,
	}
}

// NamedQuery prepares prepares the command for execution.
func (cmd *routine) NamedQuery() (string, map[string]Param) {
	namedQuery := sqlx.Rebind(sqlx.NAMED, cmd.query)
	namedParams := prepareParams(cmd.params)
	return namedQuery, namedParams
}

func (cmd *routine) Prepare(provider *sqlexec.Provider) error {
	query, err := provider.Query(cmd.name)
	if err != nil {
		return err
	}

	cmd.query = query
	return nil
}

type query struct {
	table string
	query *rql.Query
	param *rql.Params
}

// RQL create a new command from raw query
func RQL(table string, param *RQLQuery) NamedQuery {
	return &query{
		table: table,
		query: param,
	}
}

// NamedQuery prepares prepares the command for execution.
func (cmd *query) NamedQuery() (string, map[string]Param) {
	if cmd.param == nil {
		return "", nil
	}

	namedQuery := sqlx.Rebind(sqlx.NAMED, cmd.build(cmd.param))
	namedParams := prepareParams(cmd.param.FilterArgs)
	return namedQuery, namedParams
}

func (cmd *query) Prepare(model interface{}) error {
	model = reflectSliceElem(model)

	parser, err := rql.NewParser(rql.Config{
		Model:    model,
		FieldSep: ".",
	})

	if err != nil {
		return err
	}

	body, err := json.Marshal(cmd.query)
	if err != nil {
		return err
	}

	cmd.param, err = parser.Parse(body)
	return err
}

func (cmd *query) build(param *rql.Params) string {
	buffer := &bytes.Buffer{}

	fmt.Fprintf(buffer, "SELECT * FROM %s", cmd.table)

	if param != nil {
		if param.FilterExp != "" {
			fmt.Fprintf(buffer, " WHERE %v", param.FilterExp)
		}

		if param.Sort != "" {
			fmt.Fprintf(buffer, " ORDER BY %s", param.Sort)
		}

		if param.Limit > 0 {
			fmt.Fprintf(buffer, " LIMIT %d", param.Limit)
		}

		if param.Offset > 0 {
			fmt.Fprintf(buffer, " OFFSET %d", param.Offset)
		}
	}

	return buffer.String()
}

func prepareParams(values []Param) map[string]Param {
	params := make(map[string]interface{})
	index := 1

	for _, param := range values {
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
			for key, value := range bindParam(arg) {
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

func bindParam(param Param) map[string]interface{} {
	value := reflect.ValueOf(param)

	for value = reflect.ValueOf(param); value.Kind() == reflect.Ptr; {
		value = value.Elem()
	}

	switch {
	case value.Kind() == reflect.Struct:
		return reflectMap(value)
	default:
		return map[string]interface{}{
			"?": reflect.ValueOf(param).Interface(),
		}
	}
}

func reflectSliceElem(v interface{}) interface{} {
	mType := reflect.TypeOf(v)

	if mType.Kind() == reflect.Ptr {
		mType = mType.Elem()
	}

	if mType.Kind() == reflect.Slice {
		mType = mType.Elem()
	}

	if mType.Kind() == reflect.Ptr {
		mType = mType.Elem()
	}

	return reflect.New(mType).Interface()
}

func reflectMap(v reflect.Value) map[string]interface{} {
	params := make(map[string]interface{})

	for key, value := range mapper.FieldMap(v) {
		key = strings.ToLower(key)
		params[key] = value.Interface()
	}

	return params
}
