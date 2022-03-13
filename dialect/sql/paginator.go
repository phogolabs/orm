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
	cursor   *cursor
	err      error
}

// Token returns the token
func (pg *Paginator) Token() string {
	return pg.cursor.text()
}

// Error returns the underlying error
func (pg *Paginator) Error() error {
	return pg.err
}

// Dialect returns the dialect
func (pg *Paginator) Dialect() string {
	return pg.selector.dialect
}

// SetDialect sets the dialect
func (pg *Paginator) SetDialect(dialect string) {
	pg.selector.SetDialect(dialect)
}

// Query returns the query
func (pg *Paginator) Query() (string, []interface{}) {
	return pg.selector.Query()
}

// Cursor returns the next cursor
func (pg *Paginator) Scan(target interface{}) error {
	value := reflect.Indirect(reflect.ValueOf(target))

	if value.Kind() != reflect.Slice {
		return fmt.Errorf("dialect/sql: invalid type %T. expect []interface{}", target)
	}

	count := value.Len()
	// if we do not have any items we should not proceed
	if count == 0 {
		// reset the token
		pg.cursor = &cursor{}
		return nil
	}

	if limit := pg.selector.limit; limit != nil {
		if uint64(count) < *limit {
			// reset the token
			pg.cursor = &cursor{}
			return nil
		}
	}

	// get the last item
	item := value.Index(count - 1)
	// remove the last item
	value.Set(value.Slice(0, count-1))

	columns := []string{}
	// if the order is present
	if pg.selector.order != nil {
		// extract the column names
		for _, orderBy := range pg.selector.order.columns {
			columns = append(columns, orderBy.Name)
		}
	}
	// extract the column values
	values, err := scan.Values(item.Interface(), columns...)
	if err != nil {
		return err
	}

	pg.cursor = &cursor{}

	if pg.selector.order != nil {
		// calculate the cursor
		for index, clause := range pg.selector.order.columns {
			if index >= len(values) {
				return fmt.Errorf("sql: the order clauses should have valid cursor vector")
			}

			pg.cursor.add(&vector{
				Column:    clause.Name,
				Direction: clause.Direction,
				Value:     values[index],
			})
		}
	}

	return nil
}

func (pg *Paginator) seek(token string) *Paginator {
	if limit := pg.selector.limit; limit != nil {
		next := *limit + 1
		pg.selector.limit = &next
	}

	if err := pg.cursor.decode(token); err != nil {
		pg.err = err
		return pg
	}

	if err := pg.cursor.predicate(pg.selector); err != nil {
		pg.err = err
		return pg
	}

	if err := pg.cursor.order(pg.selector); err != nil {
		pg.err = err
		return pg
	}

	if pg.selector.order == nil || len(pg.selector.order.columns) == 0 {
		pg.err = fmt.Errorf("sql: query should have at least one order by clause")
	}

	return pg
}

type vector struct {
	Column    string      `json:"c"`
	Direction string      `json:"o"`
	Value     interface{} `json:"v"`
}

type cursor []*vector

func (c *cursor) decode(value string) error {
	if value == "" {
		return nil
	}

	if n := len(value) % 4; n != 0 {
		value += strings.Repeat("=", 4-n)
	}

	data, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, c)
}

func (c *cursor) text() string {
	if len(*c) == 0 {
		return ""
	}

	data, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}

func (c *cursor) add(v *vector) {
	*c = append(*c, v)
}

func (c *cursor) predicate(selector *Selector) error {
	if predicate := c.where(*c); predicate != nil {
		selector.Where(predicate)
	}

	return nil
}

func (c *cursor) where(vec []*vector) *Predicate {
	if count := len(vec); count == 0 {
		return nil
	}

	var (
		vector       = vec[0]
		predicateCmp *Predicate
		predicateEQ  = EQ(vector.Column, vector.Value)
	)

	switch vector.Direction {
	case "asc":
		predicateCmp = GT(vector.Column, vector.Value)
	case "desc":
		predicateCmp = LT(vector.Column, vector.Value)
	}

	predicate := c.where(vec[1:])

	switch {
	case predicate != nil:
		predicate = Or(predicateCmp, And(predicateEQ, predicate))
	case predicate == nil:
		predicate = Or(predicateCmp, predicateEQ)
	}

	return predicate
}

func (c *cursor) order(selector *Selector) error {
	var (
		ccount = len(*c)
		pcount = selector.order.count()
		pindex = 0
	)

	for cindex, vector := range *c {
		orderBy := &OrderByColumn{
			Name:      vector.Column,
			Direction: vector.Direction,
		}

		switch {
		case pcount == 0, cindex > pindex:
			selector.order = selector.order.add(orderBy)
		case !orderBy.Equal(selector.order.columns[pindex]):
			return fmt.Errorf("sql: pagination cursor position mismatch")
		}

		if pindex+1 < pcount {
			pindex++
		}

		if cindex != ccount-1 {
			continue
		}
	}

	return nil
}
