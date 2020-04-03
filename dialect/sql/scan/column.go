package scan

import (
	"fmt"
	"reflect"
)

// Column information
type Column struct {
	Name    string
	Options map[string]string
}

// HasOption return true if the column has option
func (c *Column) HasOption(name string) bool {
	_, ok := c.Options[name]
	return ok
}

// Columns returns the columns for given opts
func Columns(src interface{}, names ...string) ([]*Column, error) {
	var (
		value = reflect.Indirect(reflect.ValueOf(src))
		kind  = value.Kind()
	)

	if kind != reflect.Struct {
		return nil, fmt.Errorf("sql/scan: invalid type %s. expected struct as an argument", kind)
	}

	var (
		columns = make([]*Column, 0)
		meta    = mapper.TypeMap(value.Type())
	)

	if len(names) == 0 {
		for _, field := range meta.Index {
			names = append(names, field.Name)
		}
	}

	for _, name := range names {
		field, ok := meta.Names[name]

		if !ok {
			continue
		}

		column := &Column{
			Name:    name,
			Options: field.Options,
		}

		columns = append(columns, column)
	}

	return columns, nil
}
