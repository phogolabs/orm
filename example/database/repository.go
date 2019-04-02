// Package database contains an object model of database schema 'default'
// Auto-generated at Tue, 02 Apr 2019 10:26:05 CEST
package database

import (
	"github.com/phogolabs/orm"
	"github.com/phogolabs/orm/example/database/model"
)

// UserRepository represents a repository for 'users'
type UserRepository struct {
	// Gateway connects the repository to the underlying database
	Gateway *orm.Gateway
}

// SelectAll returns all User from the database
func (r *UserRepository) SelectAll() ([]*model.User, error) {
	records := []*model.User{}
	routine := orm.Routine("select-all-users")

	if err := r.Gateway.Select(&records, routine); err != nil {
		return nil, err
	}

	return records, nil
}

// SelectByPK returns a record of User for given primary key
func (r *UserRepository) SelectByPK(id int) (*model.User, error) {
	param := orm.Map{
		"id": id,
	}

	routine := orm.Routine("select-user-by-pk", param)
	record := &model.User{}

	if err := r.Gateway.SelectOne(record, routine); err != nil {
		return nil, err
	}

	return record, nil
}

// SearchAll returns all  from the database for given RQL query
func (r *UserRepository) SearchAll(query *orm.RQLQuery) ([]*model.User, error) {
	records := []*model.User{}
	routine := orm.RQL("users", query)

	if err := r.Gateway.Select(&records, routine); err != nil {
		return nil, err
	}

	return records, nil
}

// Insert inserts a record of type User into the database
func (r *UserRepository) Insert(row *model.User) error {
	routine := orm.Routine("insert-user", row)
	_, err := r.Gateway.Exec(routine)
	return err
}

// UpdateByPK updates a record of type User for given primary key
func (r *UserRepository) UpdateByPK(row *model.User) error {
	routine := orm.Routine("update-user-by-pk", row)
	_, err := r.Gateway.Exec(routine)
	return err
}

// DeleteByPK deletes a record of User for given primary key
func (r *UserRepository) DeleteByPK(id int) error {
	param := orm.Map{
		"id": id,
	}

	routine := orm.Routine("delete-user-by-pk", param)
	_, err := r.Gateway.Exec(routine)
	return err
}
