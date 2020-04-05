package orm

import (
	"time"

	"github.com/phogolabs/orm/dialect"
)

// Option represents a Gateway option
type Option interface {
	Apply(*Gateway)
}

// OptionFunc represents a function that can be used to set option
type OptionFunc func(*Gateway)

// Apply applies the option
func (fn OptionFunc) Apply(gateway *Gateway) {
	fn(gateway)
}

// WithLogger sets the logger
func WithLogger(logger dialect.Logger) Option {
	fn := func(g *Gateway) {
		g.driver = dialect.Log(g.driver, logger)
	}

	return OptionFunc(fn)
}

// WithMaxIdleConns sets the maximum number of connections in the idle
// connection pool.
//
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns,
// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit.
//
// If n <= 0, no idle connections are retained.
//
// The default max idle connections is currently 2. This may change in
// a future release.
func WithMaxIdleConns(value int) Option {
	fn := func(g *Gateway) {
		g.db.SetMaxIdleConns(value)
	}

	return OptionFunc(fn)
}

// WithMaxOpenConns sets the maximum number of open connections to the database.
//
// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
// MaxIdleConns, then MaxIdleConns will be reduced to match the new
// MaxOpenConns limit.
//
// If n <= 0, then there is no limit on the number of open connections.
// The default is 0 (unlimited).
func WithMaxOpenConns(value int) Option {
	fn := func(g *Gateway) {
		g.db.SetMaxOpenConns(value)
	}

	return OptionFunc(fn)
}

// WithConnMaxLifetime sets the maximum amount of time a connection may be reused.
//
// Expired connections may be closed lazily before reuse.
//
// If d <= 0, connections are reused forever.
func WithConnMaxLifetime(duration time.Duration) Option {
	fn := func(g *Gateway) {
		g.db.SetConnMaxLifetime(duration)
	}

	return OptionFunc(fn)
}
