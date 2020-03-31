// Package orm provides a wrapper to work with loukoum built queries as well
// maitaining database version by creating, executing and reverting SQL
// migrations.
//
// The package allows executing embedded SQL statements from script for a given
// name.
package orm

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"

	"github.com/a8m/rql"
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

	// RQLQuery is the decoded result of the user input.
	RQLQuery = rql.Query
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

// MapOf creates a map from struct that matches the provided keys. If the key
// list is empty it will use all keys.
func MapOf(v interface{}, k ...string) Map {
	return nil
}

// Map returens the parameter map
func (m Map) Map() map[string]interface{} {
	return m
}

var _ error = ErrorCollector{}

// ErrorCollector is a slice of errors
type ErrorCollector []error

// Error returns the error
func (errs ErrorCollector) Error() string {
	buffer := &bytes.Buffer{}

	for index, err := range errs {
		if index > 0 {
			fmt.Fprint(buffer, "; ")
		}

		fmt.Fprintf(buffer, err.Error())
	}

	return buffer.String()
}

// Unwrap unwrapps the collector
func (errs ErrorCollector) Unwrap() error {
	count := len(errs)

	switch {
	case count == 0:
		return nil
	case count == 1:
		return errs[0]
	default:
		return errs
	}
}
