package sql

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/phogolabs/orm/dialect/sql/scan"
)

// Paginator paginates a given selector
type Paginator struct {
	selector *Selector
	key      string
}

// PaginateBy a unique column
func (selector *Selector) PaginateBy(key string) *Paginator {
	return &Paginator{
		selector: selector,
		key:      key,
	}
}

// SetDialect sets the dialect
func (pq *Paginator) SetDialect(dialect string) {
	pq.selector.SetDialect(dialect)
}

// Seek seeks the paginator
func (pq *Paginator) Seek(cursor *Cursor) (*Paginator, error) {
	if cursor == nil {
		cursor = &Cursor{}
	}

	if limit := pq.selector.limit; limit != nil {
		next := *limit + 1
		pq.selector.limit = &next
	}

	if err := pq.order(*cursor); err != nil {
		return nil, err
	}

	if predicate := pq.where(*cursor); predicate != nil {
		pq.selector.Where(predicate)
	}

	return pq, nil
}

// Query returns the query
func (pq *Paginator) Query() (string, []interface{}) {
	return pq.selector.Query()
}

// Cursor returns the next cursor
func (pq *Paginator) Cursor(src interface{}) (*Cursor, error) {
	var (
		cursor = Cursor{}
		value  = reflect.Indirect(reflect.ValueOf(src))
	)

	if value.Kind() == reflect.Slice {
		count := value.Len()

		if count == 0 {
			return &cursor, nil
		}

		if limit := pq.selector.limit; limit != nil {
			if count < *limit {
				return &cursor, nil
			}
		}

		value = value.Index(count - 1)
		src = value.Interface()
	}

	var (
		columns     = []string{}
		orders, err = pq.orderBy()
	)

	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		columns = append(columns, order.Column)
	}

	values, err := scan.Values(src, columns...)
	if err != nil {
		return nil, err
	}

	for index, position := range orders {
		if index >= len(values) {
			return nil, fmt.Errorf("sql: the order clauses should have valid cursor vector")
		}

		nextPosition := &Vector{
			Column: position.Column,
			Order:  position.Direction,
			Value:  values[index],
		}

		cursor = append(cursor, nextPosition)
	}

	return &cursor, nil
}

// Window returns a window
func (pq *Paginator) Window(src interface{}) interface{} {
	value := reflect.Indirect(reflect.ValueOf(src))

	if value.Kind() == reflect.Slice {
		count := value.Len()

		if count == 0 {
			return src
		}

		if limit := pq.selector.limit; limit != nil {
			if count < *limit {
				return src
			}

			return value.Slice(0, count-1).Interface()
		}
	}

	return nil
}

func (pq *Paginator) order(cursor []*Vector) error {
	orders, err := pq.orderBy()
	if err != nil {
		return err
	}

	var (
		ccount = len(cursor)
		pcount = len(orders)
		pindex = 0
	)

	pagingKey, err := DecodeOrder(pq.key)
	switch {
	case err != nil:
		return err
	case pagingKey == nil:
		return fmt.Errorf("sql: pagination column not provided")
	}

	for cindex, vector := range cursor {
		candidate := vector.OrderBy()

		switch {
		case pcount == 0, cindex > pindex:
			pq.selector = pq.selector.OrderBy(candidate.String())
		case !candidate.Equal(orders[pindex]):
			return fmt.Errorf("sql: pagination cursor position mismatch")
		}

		if pindex+1 < pcount {
			pindex++
		}

		if cindex != ccount-1 {
			continue
		}

		if !candidate.Equal(pagingKey) {
			return fmt.Errorf("sql: pagination column should be placed at the end")
		}
	}

	switch {
	case ccount == 0 && pcount == 0:
		pq.selector = pq.selector.OrderBy(pagingKey.String())
	case ccount == 0 && pcount != 0:
		if candidate := orders[pindex]; !candidate.Equal(pagingKey) {
			pq.selector = pq.selector.OrderBy(pagingKey.String())
		}
	}

	return nil
}

func (pq *Paginator) where(cursor []*Vector) *Predicate {
	if count := len(cursor); count == 0 {
		return nil
	}

	var (
		vector       = cursor[0]
		predicateCmp *Predicate
		predicateEQ  = EQ(vector.Column, vector.Value)
	)

	switch vector.Order {
	case "asc":
		predicateCmp = GT(vector.Column, vector.Value)
	case "desc":
		predicateCmp = LT(vector.Column, vector.Value)
	}

	predicate := pq.where(cursor[1:])

	switch {
	case predicate != nil:
		predicate = Or(predicateCmp, And(predicateEQ, predicate))
	case predicate == nil:
		predicate = Or(predicateCmp, predicateEQ)
	}

	return predicate
}

func (pq *Paginator) orderBy() ([]*Order, error) {
	orders := []*Order{}

	for _, descriptor := range pq.selector.order {
		orderBy, err := DecodeOrderBy(descriptor)
		if err != nil {
			return nil, err
		}

		orders = append(orders, *orderBy...)
	}

	return orders, nil
}

// Vector represents an order vector
type Vector struct {
	Column string      `json:"column"`
	Order  string      `json:"order"`
	Value  interface{} `json:"value"`
}

// OrderBy representation
func (v *Vector) OrderBy() *Order {
	return &Order{
		Column:    v.Column,
		Direction: v.Order,
	}
}

// Cursor represents a cursor
type Cursor []*Vector

// DecodeCursor decodes the cursor
func DecodeCursor(token string) (*Cursor, error) {
	cursor := &Cursor{}

	if token == "" {
		return cursor, nil
	}

	if n := len(token) % 4; n != 0 {
		token += strings.Repeat("=", 4-n)
	}

	data, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cursor); err != nil {
		return nil, err
	}

	return cursor, nil
}

// String returns the cursor as a string
func (c *Cursor) String() string {
	if len(*c) == 0 {
		return ""
	}

	data, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}
