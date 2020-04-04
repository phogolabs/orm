package sql

import (
	"github.com/phogolabs/orm/dialect/sql/scan"
)

// NewDelete creates a Mutation that deletes the entity with given primary key.
func NewDelete(table string, src interface{}) *DeleteBuilder {
	var (
		deleter  = Delete(table)
		iterator = scan.IteratorOf(src)
	)

	for iterator.Next() {
		column := iterator.Column()

		if column.HasOption("primary_key") {
			deleter.Where(EQ(column.Name, iterator.Value().Interface()))
		}
	}

	return deleter
}

// NewInsert creates a Mutation that will save the entity src into the db
func NewInsert(table string, src interface{}) *InsertBuilder {
	var (
		inserter = Insert(table)
		iterator = scan.IteratorOf(src)
		columns  = make([]string, 0)
		values   = make([]interface{}, 0)
	)

	for iterator.Next() {
		column := iterator.Column()
		columns = append(columns, column.Name)
		values = append(values, iterator.Value().Interface())
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
		empty      = len(columns) == 0
		updater    = Update(table)
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

	updater = updater.Returning("*")
	return updater
}
