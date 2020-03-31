package orm

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Paginator represents a paginator
type Paginator struct {
	table     string
	where     string
	limit     int
	params    map[string]interface{}
	positions []*Position
	err       error
}

// Paginate create a new paginator query
func Paginate(table string) *Paginator {
	return &Paginator{
		table:  table,
		params: make(map[string]interface{}),
		limit:  100,
	}
}

// Order orders the query
func (pq *Paginator) Order(order string) *Paginator {
	const (
		separator = ","
		asc       = "+"
		desc      = "-"
	)

	positions := []*Position{}

	for _, field := range strings.Split(order, separator) {
		field = strings.TrimSpace(field)

		if field == "" {
			continue
		}

		position := &Position{
			Column: field,
			Order:  asc,
		}

		switch {
		case strings.HasPrefix(field, asc):
			position = &Position{
				Column: field[1:],
				Order:  asc,
			}
		case strings.HasPrefix(field, desc):
			position = &Position{
				Column: field[1:],
				Order:  desc,
			}
		}

		positions = append(positions, position)
	}

	pq.append(positions)
	return pq
}

// Where sets the where statement
func (pq *Paginator) Where(cond string, params ...Param) *Paginator {
	pq.where = fmt.Sprintf("(%s)", sqlx.Rebind(sqlx.NAMED, cond))

	for k, v := range prepareParams(params) {
		pq.params[k] = v
	}

	return pq
}

// Limit sets the page limit
func (pq *Paginator) Limit(value int) *Paginator {
	pq.limit = value
	return pq
}

// Start the pagination from a given start
func (pq *Paginator) Start(cursor *Cursor) *Paginator {
	pq.append(*cursor)
	return pq
}

// Cursor returns the cursor for given input
func (pq *Paginator) Cursor(input interface{}) *Cursor {
	var (
		next  = Cursor{}
		value = reflect.Indirect(reflect.ValueOf(input))
	)

	if value.Kind() == reflect.Slice {
		count := value.Len()

		if count == 0 {
			return &next
		}

		value = value.Index(count - 1)
	}

	for _, position := range pq.positions {
		index := &Position{
			Column: position.Column,
			Order:  position.Order,
			Value:  mapper.FieldByName(value, position.Column).Interface(),
		}

		next = append(next, index)
	}

	return &next
}

// NamedQuery prepares prepares the command for execution.
func (pq *Paginator) NamedQuery() (string, map[string]Param) {
	query := &bytes.Buffer{}
	fmt.Fprintf(query, "SELECT * FROM %v", pq.table)

	if where := and(pq.build(0), pq.where); len(where) > 0 {
		fmt.Fprintf(query, " WHERE %s", where)
	}

	if order := pq.order(); len(order) > 0 {
		fmt.Fprintf(query, " ORDER BY %s", order)
	}

	fmt.Fprintf(query, " LIMIT %d", pq.limit)

	return query.String(), pq.params
}

// Error returns any internal errors
func (pq *Paginator) Error() error {
	return pq.err
}

func (pq *Paginator) order() string {
	orderBy := &bytes.Buffer{}

	for index, position := range pq.positions {
		if index > 0 {
			fmt.Fprintf(orderBy, ", ")
		}

		switch position.Order {
		case "+":
			fmt.Fprintf(orderBy, "%v ASC", position.Column)
		case "-":
			fmt.Fprintf(orderBy, "%v DESC", position.Column)
		}
	}

	return orderBy.String()
}

func (pq *Paginator) build(index int) string {
	count := len(pq.positions)

	if count == 0 {
		return ""
	}

	var (
		predicate        = ""
		predicateCompare = ""
		predicateEqual   = ""
	)

	var (
		order  = pq.positions[index].Order
		column = pq.positions[index].Column
		value  = pq.positions[index].Value
	)

	if value != nil {
		predicateEqual = fmt.Sprintf("%s = :%s", column, column)

		switch order {
		case "+":
			predicateCompare = fmt.Sprintf("%s > :%s", column, column)
		case "-":
			predicateCompare = fmt.Sprintf("%s < :%s", column, column)
		default:
			predicateCompare = fmt.Sprintf("%s > :%s", column, column)
		}
	}

	predicate = predicateCompare

	if index < count-1 {
		predicate = or(predicateCompare,
			and(predicateEqual, pq.build(index+1)))
	}

	return predicate
}

func (pq *Paginator) append(positions []*Position) {
	err := fmt.Errorf("orm: invalid cursor position")
	count := len(pq.positions)

	for index, position := range positions {
		switch {
		case count == 0:
			pq.positions = append(pq.positions, position)
			continue
		case index >= count:
			pq.err = err
			return
		case !position.Equal(pq.positions[index]):
			pq.err = err
			return
		}
	}
}
