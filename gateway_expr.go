package orm

import (
	"github.com/phogolabs/orm/dialect/sql"
)

// Viewer represents a viewer
type Viewer interface {
	Columns() []string
}

// ViewerFunc represents a viewer function
type ViewerFunc func() []string

// Columns returns the columns
func (fn ViewerFunc) Columns() []string {
	return fn()
}

// Orderer represents an order by.
type Orderer interface {
	OrderByPath() *sql.OrderByPath
}

// OrdererFunc represents an orderer function
type OrdererFunc func() *sql.OrderByPath

// OrderByPath returns the orderer
func (fn OrdererFunc) OrderByPath() *sql.OrderByPath {
	return fn()
}

// Predicator returns the predicate
type Predicator interface {
	Predicate() *sql.Predicate
}

// PredicatorFunc represents a predicator function
type PredicatorFunc func() *sql.Predicate

// Predicate returns the predicate.
func (fn PredicatorFunc) Predicate() *sql.Predicate {
	return fn()
}

// PredicatorAnd represents a collection
type PredicatorAnd []Predicator

// And returns an instance of PredicatorAnd
func And(p ...Predicator) PredicatorAnd {
	return PredicatorAnd(p)
}

// Predicate returns the predicate.
func (x PredicatorAnd) Predicate() *sql.Predicate {
	items := []*sql.Predicate{}

	for _, factory := range x {
		if predicate := factory.Predicate(); predicate != nil {
			items = append(items, predicate)
		}
	}

	return sql.And(items...)
}

// PredicatorOr represents a collection
type PredicatorOr []Predicator

// Or returns an instance of PredicatorOr
func Or(p ...Predicator) PredicatorOr {
	return PredicatorOr(p)
}

// Predicate returns the predicate.
func (x PredicatorOr) Predicate() *sql.Predicate {
	items := []*sql.Predicate{}

	for _, factory := range x {
		if predicate := factory.Predicate(); predicate != nil {
			items = append(items, predicate)
		}
	}

	return sql.Or(items...)
}
