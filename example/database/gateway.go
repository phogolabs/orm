package database

import (
	"github.com/phogolabs/log"
	"github.com/phogolabs/orm"
	"github.com/phogolabs/orm/example/database/migration"
	"github.com/phogolabs/orm/example/database/routine"
)

func NewGateway(url string) (*orm.Gateway, error) {
	logger := log.WithField("component", "database")

	gateway, err := orm.Connect(url,
		orm.WithLogger(logger),
		orm.WithRoutine(routine.Query),
	)

	if err != nil {
		return nil, err
	}

	if err = gateway.Migrate(migration.Schema); err != nil {
		return nil, err
	}

	return gateway, nil
}
