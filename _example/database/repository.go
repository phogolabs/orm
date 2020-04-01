// Package database contains an repository of database schema ''
// Auto-generated at Fri, 12 Apr 2019 13:09:42 CEST
package database

import (
	"context"

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
	return r.SelectAllContext(context.TODO())
}

// SelectAllContext returns all User from the database
func (r *UserRepository) SelectAllContext(ctx context.Context) ([]*model.User, error) {
	records := []*model.User{}
	routine := orm.Routine("select-all-users")

	if err := r.Gateway.SelectContext(ctx, &records, routine); err != nil {
		return nil, err
	}

	return records, nil
}

// SelectByPK returns a record of User for given primary key
func (r *UserRepository) SelectByPK(id int) (*model.User, error) {
	return r.SelectByPKContext(context.TODO(), id)
}

// SelectByPKContext returns a record of User for given primary key
func (r *UserRepository) SelectByPKContext(ctx context.Context, id int) (*model.User, error) {
	param := orm.Map{
		"id": id,
	}

	routine := orm.Routine("select-user-by-pk", param)
	record := &model.User{}

	if err := r.Gateway.SelectOneContext(ctx, record, routine); err != nil {
		return nil, err
	}

	return record, nil
}

// SearchAll returns all  from the database for given RQL query
func (r *UserRepository) SearchAll(query *orm.RQLQuery) ([]*model.User, error) {
	return r.SearchAllContext(context.TODO(), query)
}

// SearchAllContext returns all  from the database for given RQL query
func (r *UserRepository) SearchAllContext(ctx context.Context, query *orm.RQLQuery) ([]*model.User, error) {
	records := []*model.User{}
	routine := orm.RQL("users", query)

	if err := r.Gateway.SelectContext(ctx, &records, routine); err != nil {
		return nil, err
	}

	return records, nil
}

// Insert inserts a record of type User into the database
func (r *UserRepository) Insert(row *model.User) error {
	return r.InsertContext(context.TODO(), row)
}

// InsertContext inserts a record of type User into the database
func (r *UserRepository) InsertContext(ctx context.Context, row *model.User) error {
	routine := orm.Routine("insert-user", row)
	_, err := r.Gateway.ExecContext(ctx, routine)
	return err
}

// UpdateByPK updates a record of type User for given primary key
func (r *UserRepository) UpdateByPK(row *model.User) error {
	return r.UpdateByPKContext(context.TODO(), row)
}

// UpdateByPKContext updates a record of type User for given primary key
func (r *UserRepository) UpdateByPKContext(ctx context.Context, row *model.User) error {
	routine := orm.Routine("update-user-by-pk", row)
	_, err := r.Gateway.ExecContext(ctx, routine)
	return err
}

// DeleteByPK deletes a record of User for given primary key
func (r *UserRepository) DeleteByPK(id int) error {
	return r.DeleteByPKContext(context.TODO(), id)
}

// DeleteByPKContext deletes a record of User for given primary key
func (r *UserRepository) DeleteByPKContext(ctx context.Context, id int) error {
	param := orm.Map{
		"id": id,
	}

	routine := orm.Routine("delete-user-by-pk", param)
	_, err := r.Gateway.ExecContext(ctx, routine)
	return err
}
