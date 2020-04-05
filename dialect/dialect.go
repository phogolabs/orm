// Copyright 2019-present Facebook Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package dialect

import (
	"context"
	"database/sql/driver"
	"time"
)

// Dialect names for external usage.
const (
	MySQL    = "mysql"
	SQLite   = "sqlite3"
	Postgres = "postgres"
	Gremlin  = "gremlin"
)

// ExecQuerier wraps the 2 database operations.
type ExecQuerier interface {
	// Exec executes a query that doesn't return rows. For example, in SQL, INSERT or UPDATE.
	// It scans the result into the pointer v. In SQL, you it's usually sql.Result.
	Exec(ctx context.Context, query string, args, v interface{}) error
	// Query executes a query that returns rows, typically a SELECT in SQL.
	// It scans the result into the pointer v. In SQL, you it's usually *sql.Rows.
	Query(ctx context.Context, query string, args, v interface{}) error
}

// Driver is the interface that wraps all necessary operations for ent clients.
type Driver interface {
	ExecQuerier
	// Tx starts and returns a new transaction.
	// The provided context is used until the transaction is committed or rolled back.
	Tx(context.Context) (Tx, error)
	// Close closes the underlying connection.
	Close() error
	// Dialect returns the dialect name of the driver.
	Dialect() string
}

// Tx wraps the Exec and Query operations in transaction.
type Tx interface {
	ExecQuerier
	driver.Tx
}

type nopTx struct {
	Driver
}

func (nopTx) Commit() error   { return nil }
func (nopTx) Rollback() error { return nil }

// NopTx returns a Tx with a no-op Commit / Rollback methods wrapping
// the provided Driver d.
func NopTx(d Driver) Tx {
	return nopTx{d}
}

// Logger represents a logger
type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
}

// LoggerDriver is a driver that logs all driver operations.
type LoggerDriver struct {
	Driver
	logger Logger
}

// Log gets a driver and an optional logging function, and returns
// a new debugged-driver that prints all outgoing operations.
func Log(d Driver, logger Logger) Driver {
	return &LoggerDriver{d, logger}
}

// Exec logs its params and calls the underlying driver Exec method.
func (d *LoggerDriver) Exec(ctx context.Context, query string, args, v interface{}) error {
	d.logger.Infof("driver.Query: query=%v", query)
	d.logger.Debugf("driver.Query: args=%v", args)
	return d.Driver.Exec(ctx, query, args, v)
}

// Query logs its params and calls the underlying driver Query method.
func (d *LoggerDriver) Query(ctx context.Context, query string, args, v interface{}) error {
	d.logger.Infof("driver.Query: query=%v", query)
	d.logger.Debugf("driver.Query: args=%v", args)
	return d.Driver.Query(ctx, query, args, v)
}

// Tx adds an log-id for the transaction and calls the underlying driver Tx command.
func (d *LoggerDriver) Tx(ctx context.Context) (Tx, error) {
	tx, err := d.Driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	id := time.Now().Unix()
	d.logger.Infof("driver.Tx(%d): started", id)
	return &LoggerTx{tx, id, d.Logger, ctx}, nil
}

// LoggerTx is a transaction implementation that logs all transaction operations.
type LoggerTx struct {
	Tx                     // underlying transaction.
	id     int64           // transaction logging id.
	logger Logger          // log function. defaults to fmt.Println.
	ctx    context.Context // underlying transaction context.
}

// Exec logs its params and calls the underlying transaction Exec method.
func (d *LoggerTx) Exec(ctx context.Context, query string, args, v interface{}) error {
	d.logger.Infof("Tx(%d).Exec: query=%v args=%v", d.id, query)
	d.logger.Debugf("Tx(%d).Exec: args=%v", d.id, args)
	return d.Tx.Exec(ctx, query, args, v)
}

// Query logs its params and calls the underlying transaction Query method.
func (d *LoggerTx) Query(ctx context.Context, query string, args, v interface{}) error {
	d.logger.Infof("Tx(%d).Query: query=%v", d.id, query)
	d.logger.Debugf("Tx(%d).Query: args=%v", args)
	return d.Tx.Query(ctx, query, args, v)
}

// Commit logs this step and calls the underlying transaction Commit method.
func (d *LoggerTx) Commit() error {
	d.logger.Infof("Tx(%d): committed", d.id)
	return d.Tx.Commit()
}

// Rollback logs this step and calls the underlying transaction Rollback method.
func (d *LoggerTx) Rollback() error {
	d.logger.Infof("Tx(%d): rollbacked", d.id)
	return d.Tx.Rollback()
}
