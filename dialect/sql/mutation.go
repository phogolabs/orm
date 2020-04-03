package sql

import (
	"strings"

	"github.com/phogolabs/orm/dialect/sql/scan"
)

// NewDelete creates a Mutation that deletes the entity with given primary key.
func NewDelete(table string, src interface{}) *DeleteBuilder {
	keys, err := scan.Columns(src)
	if err != nil {
		panic(err)
	}

	var (
		deleter = Delete(table)
		columns = []string{}
	)

	for _, key := range keys {
		columns = append(columns, key.Name)
	}

	values, err := scan.Values(src, columns...)
	if err != nil {
		panic(err)
	}

	for index, key := range keys {
		if key.HasOption("primary_key") {
			deleter.Where(EQ(key.Name, values[index]))
		}
	}

	if deleter.where == nil {
		return nil
	}

	return deleter
}

// NewInsert creates a Mutation that will save the entity src into the db
func NewInsert(table string, src interface{}) *InsertBuilder {
	keys, err := scan.Columns(src)
	if err != nil {
		panic(err)
	}

	var (
		inserter = Insert(table)
		columns  = []string{}
	)

	for _, key := range keys {
		columns = append(columns, key.Name)
	}

	values, err := scan.Values(src, columns...)
	if err != nil {
		panic(err)
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
		updater   = Update(table)
		meta, err = scan.Columns(src)
	)

	if err != nil {
		panic(err)
	}

	if len(columns) == 0 {
		for _, column := range meta {
			columns = append(columns, column.Name)
		}
	}

	hasOption := func(name string, option string) bool {
		for _, column := range meta {
			if strings.EqualFold(column.Name, name) {
				return column.HasOption(option)
			}
		}

		return false
	}

	values, err := scan.Values(src, columns...)
	if err != nil {
		panic(err)
	}

	for index, name := range columns {
		value := values[index]

		if !hasOption(name, "read_only") {
			if scan.IsNil(value) {
				updater.SetNull(name)
			} else {
				updater.Set(name, value)
			}
		}

		if hasOption(name, "primary_key") {
			updater.Where(EQ(name, value))
		}
	}

	if updater.Empty() {
		return nil
	}

	updater = updater.Returning("*")
	return updater
}
