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

// SortBy parses the given parts as asc and desc clauses
func (s *Selector) SortBy(parts ...string) *Selector {
	for _, order := range parts {
		position := CursorPositionFrom(order)

		if position == nil {
			continue
		}

		s = s.OrderBy(position.String())
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
func (pq *Paginator) Cursor(src interface{}) *Cursor {
	var (
		cursor = Cursor{}
		value  = reflect.Indirect(reflect.ValueOf(src))
	)

	if value.Kind() == reflect.Slice {
		count := value.Len()

		if count == 0 {
			return &cursor
		}

		value = value.Index(count - 1)
		src = value.Interface()
	}

	var (
		positions = pq.positions()
		columns   = []string{}
	)

	for _, position := range positions {
		columns = append(columns, position.Column)
	}

	values, err := scan.Values(src, columns...)
	if err != nil {
		panic(err)
	}

	for index, position := range pq.positions() {
		if index >= len(values) {
			panic("the order clauses should have valid cursor positions")
		}

		nextPosition := &CursorPosition{
			Column: position.Column,
			Order:  position.Order,
			Value:  values[index],
		}

		cursor = append(cursor, nextPosition)
	}

	return &cursor
}

func (pq *Paginator) order(cursor []*CursorPosition) error {
	var (
		positions = pq.positions()
		pcount    = len(positions)
		pindex    = 0
		ccount    = len(cursor)
		pagingKey = CursorPositionFrom(pq.key)
	)

	if pagingKey == nil {
		return fmt.Errorf("sql: pagination column not provided")
	}

	for cindex, candidate := range cursor {
		switch {
		case pcount == 0, cindex > pindex:
			pq.selector = pq.selector.OrderBy(candidate.String())
		case !candidate.Equal(positions[pindex]):
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
		if candidate := positions[pindex]; !candidate.Equal(pagingKey) {
			pq.selector = pq.selector.OrderBy(pagingKey.String())
		}
	}

	return nil
}

func (pq *Paginator) where(cursor []*CursorPosition) *Predicate {
	if count := len(cursor); count == 0 {
		return nil
	}

	var (
		position     = cursor[0]
		predicateCmp *Predicate
		predicateEQ  = EQ(position.Column, position.Value)
	)

	switch position.Order {
	case "asc":
		predicateCmp = GT(position.Column, position.Value)
	case "desc":
		predicateCmp = LT(position.Column, position.Value)
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

func (pq *Paginator) positions() []*CursorPosition {
	const separator = ","

	positions := []*CursorPosition{}

	for _, descriptor := range pq.selector.order {
		for _, order := range strings.Split(descriptor, separator) {
			position := CursorPositionFrom(order)

			if position == nil {
				continue
			}

			positions = append(positions, position)
		}
	}

	return positions
}

// Cursor represents a cursor
type Cursor []*CursorPosition

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

// CursorPosition represents an order position
type CursorPosition struct {
	Column string      `json:"column"`
	Order  string      `json:"order"`
	Value  interface{} `json:"value"`
}

// CursorPositionFrom returns a positions
func CursorPositionFrom(order string) *CursorPosition {
	var (
		position *CursorPosition
		parts    = strings.Fields(order)
	)

	switch len(parts) {
	case 0:
		return nil
	case 1:
		name := strings.ToLower(Unident(parts[0]))

		switch name[0] {
		case '-':
			position = &CursorPosition{
				Column: name[1:],
				Order:  "desc",
			}
		case '+':
			position = &CursorPosition{
				Column: name[1:],
				Order:  "asc",
			}
		default:
			position = &CursorPosition{
				Column: name,
				Order:  "asc",
			}
		}
	case 2:
		name := strings.ToLower(Unident(parts[0]))

		switch strings.ToLower(parts[1]) {
		case "asc":
			position = &CursorPosition{
				Column: name,
				Order:  "asc",
			}
		case "desc":
			position = &CursorPosition{
				Column: name,
				Order:  "desc",
			}
		default:
			return nil
		}
	default:
		return nil
	}

	return position
}

// Equal return true if the positions are equal
func (p *CursorPosition) Equal(pp *CursorPosition) bool {
	return strings.EqualFold(p.Column, pp.Column) && strings.EqualFold(p.Order, pp.Order)
}

// String returns the position as string
func (p *CursorPosition) String() string {
	if p.Order == "asc" {
		return Asc(p.Column)
	}

	return Desc(p.Column)
}

// Unident return the string unidented
func Unident(v string) string {
	return strings.Replace(v, "`", "", -1)
}
