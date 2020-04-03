package orm

import (
	"context"

	"github.com/phogolabs/orm/dialect/sql"
	"github.com/phogolabs/parcello"
)

type (
	// FileSystem provides with primitives to work with the underlying file system
	FileSystem = parcello.FileSystem

	// Map represents a key-value map
	Map map[string]interface{}
)

// GatewayQuerier executes the commands
type GatewayQuerier interface {
	// All executes the query and returns a list of entities.
	All(ctx context.Context, q sql.Querier, v interface{}) error

	// Only returns the only entity in the query, returns an error if not
	// exactly one entity was returned.
	Only(ctx context.Context, q sql.Querier, v interface{}) error

	// First returns the first entity in the query. Returns *NotFoundError
	// when no records were found.
	First(ctx context.Context, q sql.Querier, v interface{}) error

	// Query executes a query that returns rows, typically a SELECT in SQL.
	// It scans the result into the pointer v. In SQL, you it's usually *sql.Rows.
	Query(ctx context.Context, q sql.Querier) (*sql.Rows, error)

	// Exec executes a query that doesn't return rows. For example, in SQL, INSERT
	// or UPDATE.  It scans the result into the pointer v. In SQL, you it's usually
	// sql.Result.
	Exec(ctx context.Context, q sql.Querier) (sql.Result, error)
}

var (
	// SQL represents an SQL command
	SQL = sql.NamedQuery
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
