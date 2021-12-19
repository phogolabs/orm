package orm

import (
	"github.com/phogolabs/orm/dialect/sql"
)

// View represents a viewer
type View interface {
	Columns() []string
}

var _ View = ViewFunc(nil)

// ViewFunc represents a viewer function
type ViewFunc func() []string

// Columns returns the columns
func (fn ViewFunc) Columns() []string {
	return fn()
}

// OrderBy represents an order by.
type OrderBy interface {
	OrderByExpr() *sql.OrderByExpr
}

var _ OrderBy = OrderByFunc(nil)

// OrderByFunc represents an orderer function
type OrderByFunc func() *sql.OrderByExpr

// OrderByPath returns the orderer
func (fn OrderByFunc) OrderByExpr() *sql.OrderByExpr {
	return fn()
}

// Where returns the predicate
type Where interface {
	Predicate() *sql.Predicate
}

var _ Where = WhereFunc(nil)

// WhereFunc represents a predicator function
type WhereFunc func() *sql.Predicate

// Predicate returns the predicate.
func (fn WhereFunc) Predicate() *sql.Predicate {
	return fn()
}

var _ Where = WhereAnd{}

// WhereAnd represents a collection
type WhereAnd []Where

// And returns an instance of WhereAnd
func And(p ...Where) WhereAnd {
	return WhereAnd(p)
}

// Predicate returns the predicate.
func (x WhereAnd) Predicate() *sql.Predicate {
	items := []*sql.Predicate{}

	for _, factory := range x {
		if predicate := factory.Predicate(); predicate != nil {
			items = append(items, predicate)
		}
	}

	return sql.And(items...)
}

var _ Where = WhereOr{}

// WhereOr represents a collection
type WhereOr []Where

// Or returns an instance of WhereOr
func Or(p ...Where) WhereOr {
	return WhereOr(p)
}

// Predicate returns the predicate.
func (x WhereOr) Predicate() *sql.Predicate {
	items := []*sql.Predicate{}

	for _, factory := range x {
		if predicate := factory.Predicate(); predicate != nil {
			items = append(items, predicate)
		}
	}

	return sql.Or(items...)
}
