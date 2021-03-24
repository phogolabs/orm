package database

import (
	"github.com/phogolabs/log"
	"github.com/phogolabs/orm"
	"github.com/phogolabs/orm/example/database/migration"
	"github.com/phogolabs/orm/example/database/routine"
)

// Schema represents the database schema
var Schema = migration.Schema

// Open opens the connection
func Open(url string) (*orm.Gateway, error) {
	logger := log.WithField("component", "database")

	return orm.Connect(url,
		orm.WithLogger(logger),
		orm.WithRoutine(routine.Statement),
	)
}
