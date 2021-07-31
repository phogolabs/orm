package scan

import (
	"reflect"

	"github.com/jmoiron/sqlx/reflectx"
)

// Iterator iterates over a struct fields
type Iterator struct {
	prefix string
	value  reflect.Value
	meta   *reflectx.StructMap
	node   *Iterator
	index  int
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
	if iter.node != nil {
		// proceed
		if ok := iter.node.Next(); ok {
			return true
		}
		// go one level up
		iter.node = nil
	}

	count := len(iter.meta.Tree.Children)

	// if we are at the end return
	if count <= 0 {
		return false
	}

	if next := iter.index + 1; next < count {
		iter.index = next

		if column := iter.Column(); column.HasOption("inline") {
			iter.node = iter.createAt(column)
		}

		return true
	}

	return false
}

// Column returns the column
func (iter *Iterator) Column() *Column {
	if iter.node != nil {
		return iter.node.Column()
	}

	// current field
	field := iter.meta.Tree.Children[iter.index]

	column := &Column{
		Name:    field.Name,
		Options: field.Options,
	}

	if len(iter.prefix) > 0 {
		column.Name = iter.prefix + "_" + column.Name
	}

	return column
}

// Value returns the underlying value
func (iter *Iterator) Value() reflect.Value {
	if iter.node != nil {
		return iter.node.Value()
	}

	field := iter.meta.Tree.Children[iter.index]
	// fetch the field value
	return iter.value.FieldByIndex(field.Index)
}

func (iter *Iterator) delta() int {
	return len(iter.meta.Tree.Children) - iter.index - 1
}

func (iter *Iterator) createAt(column *Column) *Iterator {
	node := &Iterator{
		meta:  mapper.TypeMap(iter.Value().Type()),
		value: reflect.Indirect(iter.Value()),
		index: 0,
	}

	// set the prefix
	if column.HasOption("prefix") {
		node.prefix = column.Name
	}

	return node
}
