// Package oak provides a wrapper to work with loukoum built queries as well
// maitaining database version by creating, executing and reverting SQL
// migrations.
//
// The package allows executing embedded SQL statements from script for a given
// name.
package oak

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/phogolabs/parcello"
	"github.com/phogolabs/prana"
	"github.com/phogolabs/prana/sqlexec"
	"github.com/phogolabs/prana/sqlmigr"
)

// Dir implements FileSystem using the native file system restricted to a
// specific directory tree.
type Dir = parcello.Dir

// FileSystem provides with primitives to work with the underlying file system
type FileSystem = parcello.FileSystem

// P is a shortcut to a map. It facilitates passing named params to a named
// commands and queries
type P = sqlexec.P

// NamedStmt is a prepared statement that executes named queries.  Prepare it
// how you would execute a NamedQuery, but pass in a struct or map when executing.
type NamedStmt = sqlx.NamedStmt

// Entity is a destination object for given select operation.
type Entity = interface{}

// Rows is a wrapper around sql.Rows which caches costly reflect operations
// during a looped StructScan
type Rows = sqlx.Rows

// Row is a reimplementation of sql.Row in order to gain access to the underlying
// sql.Rows.Columns() data, necessary for StructScan.
type Row = sqlx.Row

// A Result summarizes an executed SQL command.
type Result = sql.Result

// TxFunc is a transaction function
type TxFunc func(tx *Tx) error

// Query represents an SQL Query that can be executed by Gateway.
type Query = sqlexec.Query

// ParseURL parses a URL and returns the database driver and connection string to the database
var ParseURL = prana.ParseURL

// Preparer prepares query for execution
type Preparer interface {
	// PrepareNamed returns a prepared named statement
	PrepareNamed(query string) (*NamedStmt, error)
}

// Migrate runs all pending migration
func Migrate(gateway *Gateway, fileSystem FileSystem) error {
	return sqlmigr.RunAll(gateway.db, fileSystem)
}

// SQL create a new command from raw query
func SQL(query string, params ...sqlexec.Param) Query {
	return sqlexec.SQL(query, params...)
}

// NamedSQL create a new command from raw query
func NamedSQL(query string, param sqlexec.Param) Query {
	return sqlexec.NamedSQL(query, param)
}
