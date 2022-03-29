package sql

import (
	"github.com/phogolabs/orm/dialect/sql/scan"
)

// InsertMutation represents an insert mutation
type InsertMutation struct {
	builder *InsertBuilder
}

// NewInsert creates a Mutation that will save the entity src into the db
func NewInsert(table string) *InsertMutation {
	return &InsertMutation{
		builder: Insert(table),
	}
}

// Entity returns the builder
func (d *InsertMutation) Entity(src interface{}) *InsertBuilder {
	var (
		iterator = scan.IteratorOf(src)
		columns  = make([]string, 0)
		values   = make([]interface{}, 0)
	)

	for iterator.Next() {
		var (
			column = iterator.Column()
			value  = iterator.Value().Interface()
		)

		if column.HasOption("auto") {
			if scan.IsEmpty(value) {
				continue
			}
		}

		columns = append(columns, column.Name)
		values = append(values, value)
	}

	return d.builder.
		Columns(columns...).
		Values(values...)
}

// UpdateMutation represents an update mutation
type UpdateMutation struct {
	builder *UpdateBuilder
}

// NewUpdate creates a Mutation that updates the entity into the db
func NewUpdate(table ...string) *UpdateMutation {
	table = append(table, "")

	return &UpdateMutation{
		builder: Update(table[0]),
	}
}

// Entity returns the builder
func (d *UpdateMutation) Entity(src interface{}, columns ...string) *UpdateBuilder {
	var (
		builder    = d.builder
		empty      = len(columns) == 0
		iterator   = scan.IteratorOf(src)
		updateable = make(map[string]interface{})
	)

	for iterator.Next() {
		var (
			column = iterator.Column()
			value  = iterator.Value().Interface()
		)

		if empty {
			columns = append(columns, column.Name)
		}

		immutable := column.HasOption("read_only") || column.HasOption("immutable") || column.HasOption("primary_key")
		// we can update only immutable columns
		if !immutable {
			updateable[column.Name] = value
		}

		// if the update statement does not have table name
		// means that we are in DO UPDATE case
		if builder.table == "" {
			continue
		}

		if column.HasOption("primary_key") {
			// TODO: we may use immutable & unique column with not null values as part of the where
			builder.Where(EQ(column.Name, value))
		}
	}

	for _, name := range columns {
		if value, ok := updateable[name]; ok {
			if scan.IsNil(value) {
				builder.SetNull(name)
			} else {
				builder.Set(name, value)
			}
		}
	}

	return builder
}

// DeleteMutation represents a delete mutation
type DeleteMutation struct {
	builder *DeleteBuilder
}

// NewDelete creates a Mutation that deletes the entity with given primary key.
func NewDelete(table string) *DeleteMutation {
	return &DeleteMutation{
		builder: Delete(table),
	}
}

// Entity returns the builder
func (d *DeleteMutation) Entity(src interface{}) *DeleteBuilder {
	var (
		builder  = d.builder
		iterator = scan.IteratorOf(src)
	)

	for iterator.Next() {
		if column := iterator.Column(); column.HasOption("primary_key") {
			// where condition
			builder = builder.Where(
				EQ(column.Name, iterator.Value().Interface()),
			)
		}
	}

	return builder
}
