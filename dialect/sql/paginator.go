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
func (s *Selector) PaginateBy(key string) *Paginator {
	return &Paginator{
		selector: s,
		key:      key,
	}
}

// OrderFrom parses the given parts as asc and desc clauses
func (s *Selector) OrderFrom(orderBy ...*Order) *Selector {
	for _, order := range orderBy {
		if order == nil {
			continue
		}

		s = s.OrderBy(order.String())
	}

	return s
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

	if err := pq.order(*cursor); err != nil {
		return nil, err
	}

	if predicate := pq.where(*cursor); predicate != nil {
		pq.selector = pq.selector.Where(predicate)
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

	pagingKey, err := OrderFrom(pq.key)
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
		predicate = predicateCmp
	}

	return predicate
}

func (pq *Paginator) orderBy() ([]*Order, error) {
	const separator = ","

	orders := []*Order{}

	for _, descriptor := range pq.selector.order {
		for _, order := range strings.Split(descriptor, separator) {
			order, err := OrderFrom(order)

			if err != nil {
				return nil, err
			}

			if order != nil {
				orders = append(orders, order)
			}
		}
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

// Order represents a order
type Order struct {
	Column    string `json:"column"`
	Direction string `json:"direction"`
}

// OrderFrom returns an order
func OrderFrom(value string) (*Order, error) {
	var (
		order *Order
		parts = strings.Fields(value)
	)

	switch len(parts) {
	case 0:
		return nil, nil
	case 1:
		name := strings.ToLower(Unident(parts[0]))

		switch name[0] {
		case '-':
			order = &Order{
				Column:    name[1:],
				Direction: "desc",
			}
		case '+':
			order = &Order{
				Column:    name[1:],
				Direction: "asc",
			}
		default:
			order = &Order{
				Column:    name,
				Direction: "asc",
			}
		}
	case 2:
		name := strings.ToLower(Unident(parts[0]))

		switch strings.ToLower(parts[1]) {
		case "asc":
			order = &Order{
				Column:    name,
				Direction: "asc",
			}
		case "desc":
			order = &Order{
				Column:    name,
				Direction: "desc",
			}
		default:
			return nil, fmt.Errorf("sql: unexpected order: %v", order)
		}
	default:
		return nil, fmt.Errorf("sql: unexepcted syntax: %v", order)
	}

	return order, nil
}

// Equal return true if the positions are equal
func (p *Order) Equal(pp *Order) bool {
	return strings.EqualFold(p.Column, pp.Column) && strings.EqualFold(p.Direction, pp.Direction)
}

// String returns the position as string
func (p *Order) String() string {
	if p.Direction == "asc" {
		return Asc(p.Column)
	}

	return Desc(p.Column)
}

// Unident return the string unidented
func Unident(v string) string {
	return strings.Replace(v, "`", "", -1)
}
