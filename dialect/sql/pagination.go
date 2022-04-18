package sql

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/phogolabs/orm/dialect/sql/scan"
)

// PaginateTable paginates a given selector
type PaginateTable struct {
	selector *Selector
	cursor   *cursor
	err      error
}

// PaginateBy paginates the given selecttor
func (x *Selector) PaginateBy(args ...string) *PaginateTable {
	if len(args) == 0 {
		args = append(args, "")
	}

	paginator := &PaginateTable{
		selector: x.Clone(),
		cursor:   &cursor{},
	}

	return paginator.seek(args[0])
}

// Token returns the token
func (pg *PaginateTable) Token() string {
	return pg.cursor.text()
}

// Err returns the underlying error
func (pg *PaginateTable) Err() error {
	err := pg.selector.Err()
	// check the paginator error
	if err == nil {
		err = pg.err
	}

	return err
}

// Dialect returns the dialect
func (pg *PaginateTable) Dialect() string {
	return pg.selector.dialect
}

// SetDialect sets the dialect
func (pg *PaginateTable) SetDialect(dialect string) {
	pg.selector.SetDialect(dialect)
}

// Query returns the query
func (pg *PaginateTable) Query() (string, []interface{}) {
	return pg.selector.Query()
}

// Scan scans the target
func (pg *PaginateTable) Scan(target interface{}) error {
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
		if int(count) < *limit {
			// reset the token
			pg.cursor = &cursor{}
			return nil
		}
	}

	// get the last item
	item := value.Index(count - 1)
	// remove the last item
	value.Set(value.Slice(0, count-1))
	// selector order
	order := vectorInOrder(pg.selector)

	columns := []string{}
	// extract the column names
	for _, elem := range order {
		columns = append(columns, elem.column)
	}

	// extract the column values
	values, err := scan.Values(item.Interface(), columns...)
	if err != nil {
		return err
	}

	pg.cursor = &cursor{}

	// calculate the cursor
	for index, elem := range order {
		if index >= len(values) {
			return fmt.Errorf("sql: the order clause %q should have valid cursor vector", elem.column)
		}

		pg.cursor.add(&vector{
			Column: elem.column,
			Order:  elem.order,
			Value:  values[index],
		})
	}

	return nil
}

func (pg *PaginateTable) seek(token string) *PaginateTable {
	if limit := pg.selector.limit; limit != nil {
		next := *limit + 1
		pg.selector.limit = &next
	}

	if err := pg.cursor.decode(token); err != nil {
		pg.err = err
		return pg
	}

	if err := pg.cursor.where(pg.selector); err != nil {
		pg.err = err
		return pg
	}

	if err := pg.cursor.order(pg.selector); err != nil {
		pg.err = err
		return pg
	}

	if len(pg.selector.order) == 0 {
		pg.err = fmt.Errorf("sql: query should have at least one order by clause")
	}

	return pg
}

type vector struct {
	Column string      `json:"c"`
	Order  string      `json:"o"`
	Value  interface{} `json:"v"`
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

func (c *cursor) where(selector *Selector) error {
	if predicate := vectorInWhere(*c); predicate != nil {
		selector.Where(predicate)
	}

	return nil
}

func (c *cursor) order(selector *Selector) error {
	order := vectorInOrder(selector)

	var (
		ccount = len(*c)
		pcount = len(order)
		pindex = 0
	)

	for cindex, vector := range *c {
		source := &OrderColumn{
			column: vector.Column,
			order:  vector.Order,
		}

		switch {
		case pcount == 0, cindex > pindex:
			selector.order = append(selector.order, source)
		default:
			if target := order[pindex]; !source.Equal(target) {
				return fmt.Errorf("sql: pagination cursor order by mismatch")
			}
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

func vectorInWhere(vec []*vector) *Predicate {
	if count := len(vec); count == 0 {
		return nil
	}

	var (
		vector       = vec[0]
		predicateCmp *Predicate
		predicateEQ  = EQ(vector.Column, vector.Value)
	)

	switch vector.Order {
	case "asc":
		predicateCmp = GT(vector.Column, vector.Value)
	case "desc":
		predicateCmp = LT(vector.Column, vector.Value)
	}

	predicate := vectorInWhere(vec[1:])

	switch {
	case predicate != nil:
		predicate = Or(predicateCmp, And(predicateEQ, predicate))
	case predicate == nil:
		predicate = Or(predicateCmp, predicateEQ)
	}

	return predicate
}

func vectorInOrder(selector *Selector) []*OrderColumn {
	order := []*OrderColumn{}

	for _, elem := range selector.order {
		switch x := elem.(type) {
		case string:
			order = append(order, OrderColumnBy(x))
		case *Order:
			order = append(order, x.columns...)
		case *OrderColumn:
			order = append(order, x)
		default:
			panic("sql: pagination cursor order by ambiguous")
		}
	}

	return order
}
