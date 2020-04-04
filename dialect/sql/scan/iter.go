package scan

import (
	"reflect"

	"github.com/jmoiron/sqlx/reflectx"
)

// Iterator iterates over a struct fields
type Iterator struct {
	value reflect.Value
	meta  *reflectx.StructMap
	index int
}

// IteratorOf creates a new iterator
func IteratorOf(src interface{}) *Iterator {
	var (
		value = reflect.Indirect(reflect.ValueOf(src))
		meta  = mapper.TypeMap(value.Type())
	)

	return &Iterator{
		value: value,
		meta:  meta,
		index: -1,
	}
}

// Next progress
func (i *Iterator) Next() bool {
	count := len(i.meta.Index)

	if count == 0 {
		return false
	}

	if next := i.index + 1; next < count {
		i.index = next
		return true
	}

	return false
}

// Column returns the column
func (i *Iterator) Column() *Column {
	field := i.meta.Index[i.index]

	return &Column{
		Name:    field.Name,
		Options: field.Options,
	}
}

// Value returns the underlying value
func (i *Iterator) Value() reflect.Value {
	field := i.meta.Index[i.index]
	return i.value.FieldByIndex(field.Index)
}
