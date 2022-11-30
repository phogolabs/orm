package sql

import (
	"context"
	"database/sql"
	"io/fs"

	"github.com/jmoiron/sqlx"
	"github.com/phogolabs/prana/sqlexec"
	"github.com/phogolabs/prana/sqlmigr"
)

// ErrNoRows is returned by Scan when QueryRow doesn't return a
// row. In such a case, QueryRow returns a placeholder *Row value that
// defers this error until a Scan.
var ErrNoRows = sql.ErrNoRows

// FileSystem represents the SQL filesystem
type FileSystem = fs.FS

// Provider represents a routine provider
type Provider = sqlexec.Provider

// Migrate runs the migrations
func (d Driver) Migrate(storage FileSystem) error {
	db := sqlx.NewDb(d.DB(), d.name)
	// execute the migration
	return sqlmigr.RunAll(db, storage)
}

// Ping pings the server
func (d Driver) Ping(ctx context.Context) error {
	return d.DB().PingContext(ctx)
}
