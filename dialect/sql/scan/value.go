package scan

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/jmoiron/sqlx/reflectx"
)

var (
	mapper       = reflectx.NewMapper("db")
	namedArgType = reflect.TypeOf(sql.NamedArg{})
)

// IsNil returns true if the value is nil
func IsNil(src interface{}) bool {
	if src == nil {
		return true
	}

	value := reflect.ValueOf(src)

	switch value.Kind() {
	case reflect.Ptr:
		return value.IsNil()
	case reflect.Chan:
		return value.IsNil()
	case reflect.Func:
		return value.IsNil()
	case reflect.Map:
		return value.IsNil()
	case reflect.Interface:
		return value.IsNil()
	case reflect.Slice:
		return value.IsNil()
	}

	return false
}

// IsEmpty returns true if the value is empty
func IsEmpty(src interface{}) bool {
	if src == nil {
		return true
	}

	value := reflect.ValueOf(src)

	switch value.Kind() {
	case reflect.Ptr:
		return value.IsNil()
	case reflect.Chan:
		return value.IsNil()
	case reflect.Func:
		return value.IsNil()
	case reflect.Map:
		return value.IsNil()
	case reflect.Interface:
		return value.IsNil()
	case reflect.Slice:
		return value.IsNil()
	}

	return value.IsZero()
}

// Args returns the arguments
func Args(src []interface{}, columns ...string) ([]interface{}, error) {
	if len(columns) == 0 {
		return src, nil
	}

	args := make([]interface{}, 0)

	for _, arg := range src {
		param := reflect.ValueOf(arg)
		param = reflect.Indirect(param)

		switch k := param.Type().Kind(); {
		case k == reflect.String || k >= reflect.Bool && k <= reflect.Float64:
			args = append(args, arg)
		default:
			values, err := valuesOf(param, columns)
			if err != nil {
				return nil, err
			}
			args = append(args, values...)
		}
	}

	return args, nil
}

// Values scans a struct and returns the values associated with the columns
// provided. Only simple value types are supported (i.e. Bool, Ints, Uints,
// Floats, Interface, String, NamedArg)
func Values(src interface{}, columns ...string) ([]interface{}, error) {
	// prepare the input
	for index, column := range columns {
		columns[index] = unquote(column)
	}

	value, err := valueOf(src)
	if err != nil {
		return nil, err
	}

	return valuesOf(value, columns)
}

func valuesOf(value reflect.Value, columns []string) ([]interface{}, error) {
	switch value.Kind() {
	case reflect.Struct:
		switch value.Type() {
		case namedArgType:
			return valuesOfNamedArg(value, columns)
		default:
			return valuesOfStruct(value, columns)
		}
	case reflect.Map:
		return valuesOfMap(value, columns)
	default:
		return nil, fmt.Errorf("sql/scan: invalid type %s. expected struct or map as an argument", value.Kind())
	}
}

func valuesOfStruct(target reflect.Value, columns []string) ([]interface{}, error) {
	kind := target.Kind()

	if kind != reflect.Struct {
		return nil, fmt.Errorf("sql/scan: invalid type %s. expected struct as an argument", kind)
	}

	values := make([]interface{}, 0)

	if len(columns) == 0 {
		meta := mapper.TypeMap(target.Type())

		iter := &Iterator{
			meta:  meta,
			value: target,
			index: -1,
		}

		for iter.Next() {
			column := iter.Column()
			columns = append(columns, column.Name)
		}
	}

	for _, name := range columns {
		if field := fieldByName(target.Type(), name); field != nil {
			// find the value
			value := valueByIndex(target, field.Index).Interface()
			// append it
			values = append(values, value)
		}
	}

	return values, nil
}

func valuesOfNamedArg(target reflect.Value, columns []string) ([]interface{}, error) {
	var (
		values = []interface{}{}
		arg    = target.Interface().(sql.NamedArg)
	)

	if len(columns) == 0 {
		columns = append(columns, arg.Name)
	}

	for _, column := range columns {
		if column == arg.Name {
			values = append(values, arg.Value)
			break
		}
	}

	return values, nil
}

func valuesOfMap(value reflect.Value, columns []string) ([]interface{}, error) {
	var (
		empty = reflect.Value{}
		kind  = value.Kind()
	)

	if kind != reflect.Map {
		return nil, fmt.Errorf("sql/scan: invalid type %s. expected map as an argument", kind)
	}

	if keyKind := value.Type().Key().Kind(); keyKind != reflect.String {
		return nil, fmt.Errorf("sql/scan: invalid type %s. expected string as an key", keyKind)
	}

	if len(columns) == 0 {
		for _, key := range value.MapKeys() {
			columns = append(columns, key.String())
		}
	}

	values := make([]interface{}, 0)

	for _, name := range columns {
		param := value.MapIndex(reflect.ValueOf(name))

		if param == empty {
			continue
		}

		values = append(values, param.Interface())
	}

	return values, nil
}

func valueOf(src interface{}) (reflect.Value, error) {
	var (
		empty = reflect.Value{}
		typ   = reflect.TypeOf(src)
		kind  = typ.Kind()
	)

	if kind != reflect.Ptr {
		return empty, fmt.Errorf("sql/scan: invalid type %s. expected pointer as an argument", kind)
	}

	value := reflect.ValueOf(src)
	value = reflect.Indirect(value)
	return value, nil
}
