// Package oak provides a wrapper to work with loukoum built queries as well
// maitaining database version by creating, executing and reverting SQL
// migrations.
//
// The package allows executing embedded SQL statements from script for a given
// name.
package oak

import (
	"database/sql"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/phogolabs/parcello"
	"github.com/phogolabs/prana/sqlexec"
	"github.com/phogolabs/prana/sqlmigr"
)

// P is a shortcut to a map. It facilitates passing named params to a named
// commands and queries
type P = map[string]interface{}

// FileSystem provides with primitives to work with the underlying file system
type FileSystem = parcello.FileSystem

// Dir implements FileSystem using the native file system restricted to a
// specific directory tree.
type Dir = parcello.Dir

// Query represents an SQL Query that can be executed by Gateway.
type Query interface {
	// Prepare prepares the query for execution. It returns the actual query and
	// a maps of its arguments.
	Prepare() (string, map[string]interface{})
}

// Preparer prepares query for execution
type Preparer interface {
	// PrepareNamed returns a prepared named statement
	PrepareNamed(query string) (*NamedStmt, error)
}

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

var provider *sqlexec.Provider

func init() {
	provider = &sqlexec.Provider{}
}

// Setup setups the oak environment for us
func Setup(gateway *Gateway, manager *parcello.Manager) error {
	script, err := manager.Root("script")
	if err != nil {
		return err
	}

	if err = LoadSQLCommandsFrom(script); err != nil {
		return err
	}

	migration, err := manager.Root("migration")
	if err != nil {
		return err
	}

	if err := Migrate(gateway, migration); err != nil {
		return err
	}

	return nil
}

// Migrate runs all pending migration
func Migrate(gateway *Gateway, fileSystem FileSystem) error {
	return sqlmigr.RunAll(gateway.db, fileSystem)
}

// LoadSQLCommandsFromReader loads all commands from a given reader.
func LoadSQLCommandsFromReader(r io.Reader) error {
	_, err := provider.ReadFrom(r)
	return err
}

// LoadSQLCommandsFrom loads all script commands from a given directory. Note that all
// scripts should have .sql extension.
func LoadSQLCommandsFrom(fileSystem FileSystem) error {
	return provider.ReadDir(fileSystem)
}

// Command returns a command for given name and parameters. The operation can
// panic if the command cannot be found.
func Command(name string, params ...sqlexec.Param) Query {
	cmd, err := provider.Command(name, params...)

	if err != nil {
		panic(err)
	}

	return cmd
}

// NamedCommand returns a command for given name and map the parameters as
// named. The operation can panic if the command cannot be found.
func NamedCommand(name string, params ...sqlexec.Param) Query {
	cmd, err := provider.NamedCommand(name, params...)

	if err != nil {
		panic(err)
	}

	return cmd
}

// SQL create a new command from raw query
func SQL(query string, params ...sqlexec.Param) Query {
	return sqlexec.SQL(query, params...)
}

// NamedSQL create a new command from raw query
func NamedSQL(query string, params ...sqlexec.Param) Query {
	return sqlexec.NamedSQL(query, params...)
}

// ParseURL parses a URL and returns the database driver and connection string to the database
func ParseURL(conn string) (string, string, error) {
	uri, err := url.Parse(conn)
	if err != nil {
		return "", "", err
	}

	driver := strings.ToLower(uri.Scheme)

	switch driver {
	case "mysql", "sqlite3":
		source := strings.Replace(conn, fmt.Sprintf("%s://", driver), "", -1)
		return driver, source, nil
	default:
		return driver, conn, nil
	}
}
