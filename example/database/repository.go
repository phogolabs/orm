// Package database contains an repository of database schema ''
// Auto-generated at Fri, 12 Apr 2019 13:09:42 CEST
package database

import (
	"context"

	"github.com/phogolabs/orm"
	"github.com/phogolabs/orm/example/database/ent"
)

// UserRepository represents a repository for 'users'
type UserRepository struct {
	// Gateway connects the repository to the underlying database
	Gateway *orm.Gateway
}

// AllUsers returns all User from the database
func (r *UserRepository) AllUsers(ctx context.Context) ([]*ent.User, error) {
	var (
		entities = []*ent.User{}
		routine  = orm.Routine("select-all-users")
	)

	if err := r.Gateway.All(ctx, routine, &entities); err != nil {
		return nil, err
	}

	return entities, nil
}

// InsertUser inserts a record of type User into the database
func (r *UserRepository) InsertUser(ctx context.Context, entity *ent.User) error {
	routine := orm.Routine("insert-user", entity)
	_, err := r.Gateway.Exec(ctx, routine)
	return err
}
