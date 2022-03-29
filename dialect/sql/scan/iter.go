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
		meta:  meta,
		value: value,
		index: -1,
	}
}

// Next progress
func (iter *Iterator) Next() bool {
	count := len(iter.meta.Tree.Children)

	// if we are at the end return
	if count <= 0 {
		return false
	}

	if next := iter.index + 1; next < count {
		iter.index = next
		// done
		return true
	}

	return false
}

// Column returns the column
func (iter *Iterator) Column() *Column {
	// current field
	parent := iter.meta.Tree.Children[iter.index]

	column := &Column{
		Name:    parent.Name,
		Options: parent.Options,
	}

	if name, ok := parent.Options["foreign_key"]; ok {
		column.Name = name
	}

	return column
}

// Value returns the underlying value
func (iter *Iterator) Value() reflect.Value {
	parent := iter.meta.Tree.Children[iter.index]
	// fetch the value
	value := iter.value.FieldByIndex(parent.Index)

	if name, ok := parent.Options["reference_key"]; ok {
		value = reflect.Indirect(value)
		// prepare the meta mapper
		meta := mapper.TypeMap(value.Type())
		// find the actual value
		value = value.FieldByIndex(meta.GetByPath(name).Index)
	}

	// fetch the field value
	return value
}
