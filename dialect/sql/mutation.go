package sql

import (
	"github.com/phogolabs/orm/dialect/sql/scan"
)

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
		deleter  = d.builder
		iterator = scan.IteratorOf(src)
	)

	for iterator.Next() {
		column := iterator.Column()

		if column.HasOption("primary_key") {
			deleter = deleter.Where(EQ(column.Name, iterator.Value().Interface()))
		}
	}

	return deleter
}

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
func NewUpdate(table string) *UpdateMutation {
	return &UpdateMutation{
		builder: Update(table),
	}
}

// Entity returns the builder
func (d *UpdateMutation) Entity(src interface{}, columns ...string) *UpdateBuilder {
	var (
		updater    = d.builder
		empty      = len(columns) == 0
		iterator   = scan.IteratorOf(src)
		updateable = make(map[string]interface{})
	)

	for iterator.Next() {
		column := iterator.Column()

		if empty {
			columns = append(columns, column.Name)
		}

		if !column.HasOption("read_only") {
			updateable[column.Name] = iterator.Value().Interface()
		}

		if column.HasOption("primary_key") {
			updater.Where(EQ(column.Name, iterator.Value().Interface()))
		}
	}

	for _, name := range columns {
		if value, ok := updateable[name]; ok {
			if scan.IsNil(value) {
				updater.SetNull(name)
			} else {
				updater.Set(name, value)
			}
		}
	}

	return updater
}
