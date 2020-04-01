package sql

import (
	"reflect"

	"github.com/jmoiron/sqlx/reflectx"
)

// NewDelete creates a Mutation that deletes the entity with given primary key.
func NewDelete(table string, src interface{}) *DeleteBuilder {
	var (
		value   = reflect.Indirect(reflect.ValueOf(src))
		meta    = mapper.TypeMap(value.Type())
		deleter = Delete(table)
		count   = 0
	)

	for _, field := range meta.Index {
		if _, ok := field.Options["primary_key"]; ok {
			value := reflectx.FieldByIndexes(value, field.Index)
			deleter.Where(EQ(field.Name, value.Interface()))
			count++
		}
	}

	if count == 0 {
		deleter = nil
	}

	return deleter
}

// NewInsert creates a Mutation that will save the entity src into the db
func NewInsert(table string, src interface{}) *InsertBuilder {
	var (
		value    = reflect.Indirect(reflect.ValueOf(src))
		meta     = mapper.TypeMap(value.Type())
		inserter = Insert(table)
		columns  []string
		values   []interface{}
	)

	for _, field := range meta.Index {
		value := reflectx.FieldByIndexes(value, field.Index)
		values = append(values, value.Interface())
		columns = append(columns, field.Name)
	}

	inserter = inserter.
		Columns(columns...).
		Values(values...).
		Returning("*")

	return inserter
}

// NewUpdate creates a Mutation that updates the entity into the db
func NewUpdate(table string, src interface{}, columns ...string) *UpdateBuilder {
	var (
		value   = reflect.Indirect(reflect.ValueOf(src))
		meta    = mapper.TypeMap(value.Type())
		updater = Update(table)
		count   = 0
	)

	if len(columns) == 0 {
		for _, field := range meta.Index {
			columns = append(columns, field.Name)
		}
	}

	for _, name := range columns {
		field := meta.GetByPath(name)

		if _, ok := field.Options["read_only"]; ok {
			continue
		}

		value := mapper.FieldByName(value, name)

		if value.Kind() == reflect.Ptr {
			if value.IsNil() || value.Elem().IsZero() {
				updater = updater.SetNull(field.Name)
				continue
			}
		}

		updater = updater.Set(field.Name, value.Interface())
	}

	for _, field := range meta.Index {
		if _, ok := field.Options["primary_key"]; ok {
			value := reflectx.FieldByIndexes(value, field.Index)
			updater.Where(EQ(field.Name, value.Interface()))
			count++
		}
	}

	if count == 0 {
		updater = nil
	}

	updater = updater.Returning("*")
	return updater
}
