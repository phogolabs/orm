package dialect

import (
	"context"
	"time"

	"github.com/phogolabs/log"
)

// Logger represents a logger
type Logger = log.Logger

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
	var (
		start  = time.Now()
		logger = d.logger
	)

	err := d.Driver.Exec(ctx, query, args, v)

	logger = logger.WithField("sql.query", query)
	logger = logger.WithField("sql.param", args)
	logger = logger.WithField("sql.duration", time.Since(start).String())

	if err != nil {
		logger.WithError(err).Errorf("query.exec fail")
		return err
	}

	logger.Infof("query.exec success")
	return nil
}

// Query logs its params and calls the underlying driver Query method.
func (d *LoggerDriver) Query(ctx context.Context, query string, args, v interface{}) error {
	var (
		start  = time.Now()
		logger = d.logger
	)

	err := d.Driver.Query(ctx, query, args, v)

	logger = logger.WithField("sql.query", query)
	logger = logger.WithField("sql.param", args)
	logger = logger.WithField("sql.duration", time.Since(start).String())

	if err != nil {
		logger.WithError(err).Errorf("query.exec fail")
		return err
	}

	logger.Infof("query.exec success")
	return nil
}

// Tx adds an log-id for the transaction and calls the underlying driver Tx command.
func (d *LoggerDriver) Tx(ctx context.Context) (Tx, error) {
	logger := d.logger.WithField("sql.tx", time.Now().Unix())

	tx, err := d.Driver.Tx(ctx)
	if err != nil {
		logger.WithError(err).Errorf("tx.start fail")
		return nil, err
	}

	logger.Infof("tx.start success")
	return &LoggerTx{tx, logger, ctx}, nil
}

// LoggerTx is a transaction implementation that logs all transaction operations.
type LoggerTx struct {
	Tx                     // underlying transaction.
	logger Logger          // log function. defaults to fmt.Println.
	ctx    context.Context // underlying transaction context.
}

// Exec logs its params and calls the underlying transaction Exec method.
func (d *LoggerTx) Exec(ctx context.Context, query string, args, v interface{}) error {
	var (
		start  = time.Now()
		logger = d.logger
	)

	err := d.Tx.Exec(ctx, query, args, v)

	logger = logger.WithField("sql.query", query)
	logger = logger.WithField("sql.param", args)
	logger = logger.WithField("sql.duration", time.Since(start).String())

	if err != nil {
		logger.WithError(err).Errorf("query.exec fail")
		return err
	}

	logger.Infof("query.exec success")
	return nil
}

// Query logs its params and calls the underlying transaction Query method.
func (d *LoggerTx) Query(ctx context.Context, query string, args, v interface{}) error {
	var (
		start  = time.Now()
		logger = d.logger
	)

	err := d.Tx.Query(ctx, query, args, v)

	logger = logger.WithField("sql.query", query)
	logger = logger.WithField("sql.param", args)
	logger = logger.WithField("sql.duration", time.Since(start).String())

	if err != nil {
		logger.WithError(err).Errorf("query.exec fail")
		return err
	}

	logger.Infof("query.exec success")
	return nil
}

// Commit logs this step and calls the underlying transaction Commit method.
func (d *LoggerTx) Commit() error {
	var (
		logger = d.logger
		err    = d.Tx.Commit()
	)

	if err != nil {
		logger = d.logger.WithError(err)
		logger.Errorf("tx.commit fail")
		return err
	}

	logger.Infof("tx.commit success")
	return nil
}

// Rollback logs this step and calls the underlying transaction Rollback method.
func (d *LoggerTx) Rollback() error {
	var (
		logger = d.logger
		err    = d.Tx.Rollback()
	)

	if err != nil {
		logger = d.logger.WithError(err)
		logger.Errorf("tx.rollback fail")
		return err
	}

	logger.Infof("tx.rollback success")
	return nil
}
