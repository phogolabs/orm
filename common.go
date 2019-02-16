// Package orm provides a wrapper to work with loukoum built queries as well
// maitaining database version by creating, executing and reverting SQL
// migrations.
//
// The package allows executing embedded SQL statements from script for a given
// name.
package orm

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/phogolabs/parcello"
)

type (
	// FileSystem provides with primitives to work with the underlying file system
	FileSystem = parcello.FileSystem

	// Rows is a wrapper around sql.Rows which caches costly reflect operations
	// during a looped StructScan
	Rows = sqlx.Rows

	// Row is a reimplementation of sql.Row in order to gain access to the underlying
	// sql.Rows.Columns() data, necessary for StructScan.
	Row = sqlx.Row

	// A Result summarizes an executed SQL command.
	Result = sql.Result

	// TxFunc is a transaction function
	TxFunc func(tx *Tx) error

	// TxContextFunc is a transaction function
	TxContextFunc func(ctx context.Context, tx *Tx) error

	// Entity is a destination object for given select operation.
	Entity = interface{}

	// Param is a command parameter for given query.
	Param = interface{}
)

// Mapper provides a map of parameters
type Mapper interface {
	// Map returens the parameter map
	Map() map[string]interface{}
}

// Preparer prepares query for execution
type Preparer interface {
	// Rebind rebinds the query
	Rebind(query string) string
	// PrepareNamed returns a prepared named statement
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
}

// NamedQuery returns the underlying query
type NamedQuery interface {
	// NamedQuery prepares the query
	NamedQuery() (string, map[string]Param)
}

// Map is a shortcut to a map. It facilitates passing named params to a named
// commands and queries
type Map map[string]interface{}

// Map returens the parameter map
func (m Map) Map() map[string]interface{} {
	return m
}
