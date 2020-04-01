package orm

import (
	"github.com/phogolabs/orm/dialect/sql"
	"github.com/phogolabs/parcello"
)

type (
	// FileSystem provides with primitives to work with the underlying file system
	FileSystem = parcello.FileSystem

	// Map represents a key-value map
	Map map[string]interface{}
)

var (
	// SQL represents an SQL command
	SQL = sql.Command
	// Routine represents an SQL routine
	Routine = sql.Routine
)

var (
	// NewDelete creates a Mutation that deletes the entity with given primary key.
	NewDelete = sql.NewDelete

	// NewInsert creates a Mutation that will save the entity src into the db
	NewInsert = sql.NewInsert

	// NewUpdate creates a Mutation that updates the entity into the db
	NewUpdate = sql.NewUpdate
)
