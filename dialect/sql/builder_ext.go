package sql

import (
	"bytes"
	"fmt"
	"strings"
)

// Name returns the name
func (s *Selector) Name() string {
	switch view := s.from.(type) {
	case *WithBuilder:
		return view.name
	case *SelectTable:
		return view.name
	case *Selector:
		return view.as
	default:
		panic(fmt.Sprintf("unhandled TableView type %T", s.from))
	}
}

// TableViews returns the table views
func (s *Selector) TableViews() []TableView {
	views := []TableView{}
	views = append(views, s.from)

	for _, item := range s.joins {
		views = append(views, item.table)
	}

	return views
}

// Name returns the name
func (s *SelectTable) Name() string {
	return s.name
}

// SelectorFunc represents a selector function
type SelectorFunc func(*Selector)

// Selection represents a selection
type Selection struct {
	fn SelectorFunc
}

// SelectionBy returns a new selection
func SelectionBy(fn SelectorFunc) *Selection {
	return &Selection{fn: fn}
}

// SelectionWith creates a selection
func SelectionWith(selection ...*Selection) *Selection {
	fn := func(selector *Selector) {
		for _, element := range selection {
			// execute the selection
			element.Select(selector)
		}
	}

	return &Selection{fn: fn}
}

// Select selects the selector
func (s *Selection) Select(selector *Selector) {
	if s != nil {
		s.fn(selector)
	}
}

// Projection represents an order
type Projection struct {
	columns []*ProjectColumn
	dialect string
	total   int
}

// ProjectBy returns the view
func (s *SelectTable) ProjectBy(clause *Projection) *Projection {
	output := &Projection{
		dialect: clause.dialect,
		total:   clause.total,
	}

	for _, element := range clause.columns {
		column := element.Clone()
		column.name = s.C(column.name)
		// append the column
		output.columns = append(output.columns, column)
	}

	return output
}

// ProjectBy returns the view
func (w *WithBuilder) ProjectBy(clause *Projection) *Projection {
	output := &Projection{
		dialect: clause.dialect,
		total:   clause.total,
	}

	for _, element := range clause.columns {
		column := element.Clone()
		column.name = w.C(column.name)
		// append the column
		output.columns = append(output.columns, column)
	}

	return output
}

// ProjectBy returns a seelect by
func ProjectBy(clauses ...string) *Projection {
	selection := &Projection{}

	for _, clause := range clauses {
		exprs := strings.Split(clause, ",")
		// parse each element
		for _, expr := range exprs {
			if column := ProjectColumnBy(expr); column != nil {
				// add the expression the collection
				selection.columns = append(selection.columns, column)
			}
		}
	}

	return selection
}

// SetDialect sets the dialect
func (x *Projection) SetDialect(dialect string) {
	x.dialect = dialect
}

// Dialect returns the dialect
func (x *Projection) Dialect() string {
	return x.dialect
}

// Total returns the total value
func (x *Projection) Total() int {
	return x.total
}

// SetTotal sets the totla arguments
func (x *Projection) SetTotal(total int) {
	x.total = total
}

// Query returns the order by clause
func (x *Projection) Query() (string, []interface{}) {
	queriers := []Querier{}

	for _, column := range x.columns {
		column.SetDialect(x.dialect)
		column.SetTotal(x.total)
		// add the column to the collection
		queriers = append(queriers, column)
	}

	b := &Builder{dialect: x.dialect, total: x.total}
	b.JoinComma(queriers...)
	return b.String(), nil
}

// As creates an alias
func (x *Projection) As(prefix string, extra ...string) *Projection {
	for _, column := range x.columns {
		name := column.name
		// prepare the name
		name = strings.Replace(name, `"`, "", -1)
		name = strings.Replace(name, "`", "", -1)
		// split the name
		parts := strings.Split(name, ".")
		index := len(parts) - 1
		// find the relative name
		name = parts[index]

		head := []string{}
		head = append(head, prefix)
		head = append(head, extra...)
		head = append(head, name)
		// prepare
		name = strings.Join(head, ".")
		// alias
		column.As(name)
	}

	return x
}

// ProjectColumn is a column selector.
type ProjectColumn struct {
	dialect string
	alias   string
	name    string
	total   int
}

// ProjectColumnBy returns a new column selector.
func ProjectColumnBy(name string) *ProjectColumn {
	return &ProjectColumn{
		name: name,
	}
}

// Clone clones the column
func (x *ProjectColumn) Clone() *ProjectColumn {
	return &ProjectColumn{
		dialect: x.dialect,
		alias:   x.alias,
		name:    x.name,
		total:   x.total,
	}
}

// SetDialect sets the dialect
func (x *ProjectColumn) SetDialect(dialect string) {
	x.dialect = dialect
}

// Dialect returns the dialect
func (x *ProjectColumn) Dialect() string {
	return x.dialect
}

// Total returns the total value
func (x *ProjectColumn) Total() int {
	return x.total
}

// SetTotal sets the totla arguments
func (x *ProjectColumn) SetTotal(total int) {
	x.total = total
}

// As adds the AS clause to the column selector.
func (x *ProjectColumn) As(alias string) *ProjectColumn {
	x.alias = alias
	return x
}

// Query returns the select column clause
func (x *ProjectColumn) Query() (string, []interface{}) {
	b := &Builder{dialect: x.dialect}
	b.Ident(x.name)
	if x.alias != "" {
		b.Pad().WriteString("AS")
		b.Pad().Ident(x.alias)
	}
	return b.String(), nil
}

// Order represents an order
type Order struct {
	columns []*OrderColumn
	dialect string
	total   int
}

// Columns returns the column nams
func (x *Order) Columns() []string {
	columns := []string{}
	// extract the column names
	for _, elem := range x.columns {
		columns = append(columns, elem.column)
	}

	return columns
}

// String returns the order as string
func (x *Order) String() string {
	data, err := x.MarshalText()
	if err != nil {
		panic(err)
	}
	return string(data)
}

// MarshalText encodes the receiver into UTF-8-encoded text and returns the result.
func (x *Order) MarshalText() ([]byte, error) {
	buffer := &bytes.Buffer{}

	for _, column := range x.columns {
		data, err := column.MarshalText()
		if err != nil {
			return nil, err
		}

		if buffer.Len() > 0 {
			buffer.WriteString(", ")
		}

		if _, err := buffer.Write(data); err != nil {
			return nil, err
		}
	}

	return buffer.Bytes(), nil
}

// UnmarshalText must be able to decode the form generated by MarshalText.
// UnmarshalText must copy the text if it wishes to retain the text after
// returning.
func (x *Order) UnmarshalText(data []byte) error {
	separator := []byte(",")
	// partition the text
	clauses := bytes.Split(data, separator)
	// parse each element
	for _, expr := range clauses {
		column := &OrderColumn{}
		// unmarsha the column
		if err := column.UnmarshalText(expr); err != nil {
			return err
		}
		// add the expression the collection
		x.columns = append(x.columns, column)
	}

	return nil
}

// Order returns the order
func (s *Selector) Order() *Order {
	columns := []*OrderColumn{}

	for _, elem := range s.order {
		switch x := elem.(type) {
		case string:
			columns = append(columns, OrderColumnBy(x))
		case *Order:
			columns = append(columns, x.columns...)
		case *OrderColumn:
			columns = append(columns, x)
		default:
			panic("sql: order column ambiguous")
		}
	}

	return &Order{
		columns: columns,
		total:   s.total,
		dialect: s.dialect,
	}
}

// OrderBy returns the order
func (s *SelectTable) OrderBy(clause *Order) *Order {
	order := &Order{
		dialect: clause.dialect,
		total:   clause.total,
	}

	for _, element := range clause.columns {
		element = element.Clone()
		element.column = s.C(element.column)
		// append the column
		order.columns = append(order.columns, element)
	}

	return order
}

// OrderBy returns the order
func (w *WithBuilder) OrderBy(clause *Order) *Order {
	order := &Order{
		dialect: clause.dialect,
		total:   clause.total,
	}

	for _, element := range clause.columns {
		element = element.Clone()
		element.column = w.C(element.column)
		// append the column
		order.columns = append(order.columns, element)
	}

	return order
}

// OrderBy returns an order by
func OrderBy(clauses ...string) *Order {
	order := &Order{}

	for _, clause := range clauses {
		data := []byte(clause)
		order.UnmarshalText(data)
	}

	return order
}

// OrderWith returns an order from
func OrderWith(orders ...*Order) *Order {
	order := &Order{}

	for _, element := range orders {
		order.dialect = element.dialect
		order.total += element.total

		for _, clause := range element.columns {
			order.columns = append(order.columns, clause.Clone())
		}
	}

	return order
}

// Prepend prepends the clauses
func (x *Order) Prepend(clauses ...string) *Order {
	head := OrderBy(clauses...)
	// return a new order
	return OrderWith(head, x)
}

// Map maps the order based on the provided mapping
func (x *Order) Map(mapping map[string]string) *Order {
	order := &Order{
		dialect: x.dialect,
		total:   x.total,
	}

	for _, element := range x.columns {
		name := element.column
		// prepare the name
		name = strings.Replace(name, `"`, "", -1)
		name = strings.Replace(name, "`", "", -1)
		// split the name
		parts := strings.Split(name, ".")
		index := len(parts) - 1
		// find the relative name
		name = parts[index]

		if nick, ok := mapping[element.column]; ok {
			element = element.Clone()
			element.column = strings.Replace(element.column, name, nick, -1)
			// append the order
			order.columns = append(order.columns, element)
		}
	}

	return order
}

// SetDialect sets the dialect
func (x *Order) SetDialect(dialect string) {
	x.dialect = dialect
}

// Dialect returns the dialect
func (x *Order) Dialect() string {
	return x.dialect
}

// Total returns the total value
func (x *Order) Total() int {
	return x.total
}

// SetTotal sets the totla arguments
func (x *Order) SetTotal(total int) {
	x.total = total
}

// Query returns the order by clause
func (x *Order) Query() (string, []interface{}) {
	queriers := []Querier{}

	for _, column := range x.columns {
		column.SetDialect(x.dialect)
		column.SetTotal(x.total)
		// add the column to the collection
		queriers = append(queriers, column)
	}

	b := &Builder{dialect: x.dialect, total: x.total}
	b.JoinComma(queriers...)
	return b.String(), nil
}

// OrderColumn represents an order by column
type OrderColumn struct {
	dialect string
	column  string
	order   string
	total   int
	err     error
}

// MarshalText encodes the receiver into UTF-8-encoded text and returns the result.
func (x *OrderColumn) MarshalText() ([]byte, error) {
	if x.err == nil {
		text := fmt.Sprintf("%s %s", x.column, x.order)
		return []byte(text), nil
	}

	return nil, x.err
}

// UnmarshalText must be able to decode the form generated by MarshalText.
// UnmarshalText must copy the text if it wishes to retain the text after
// returning.
func (x *OrderColumn) UnmarshalText(data []byte) error {
	data = bytes.ToLower(data)
	data = bytes.TrimSpace(data)

	var (
		name  string
		order string
	)

	elem := bytes.Fields(data)
	// prepare the name
	if len(elem) > 0 {
		name = string(elem[0])
		name = strings.Replace(name, "`", "", -1)
	}

	switch len(elem) {
	case 0:
		return nil
	case 1:
		order = string(name[0])
		// convert the expression
		switch order {
		case "+":
			order = "asc"
			name = name[1:]
		case "-":
			order = "desc"
			name = name[1:]
		default:
			order = "asc"
		}
	case 2:
		order = string(elem[1])
	}

	switch order {
	case "asc":
		*x = OrderColumn{
			column: name,
			order:  order,
		}
	case "desc":
		*x = OrderColumn{
			column: name,
			order:  order,
		}
	default:
		*x = OrderColumn{
			err: fmt.Errorf("expression %q is not valid", string(data)),
		}

		return x.err
	}

	return nil
}

// OrderColumnBy returns a order column
func OrderColumnBy(expr string) *OrderColumn {
	data := []byte(expr)
	// unmarshal
	x := &OrderColumn{}
	x.UnmarshalText(data)
	// done!
	return x
}

// Clone clones the column
func (x *OrderColumn) Clone() *OrderColumn {
	return &OrderColumn{
		dialect: x.dialect,
		column:  x.column,
		order:   x.order,
		total:   x.total,
		err:     x.err,
	}
}

// SetDialect sets the dialect
func (x *OrderColumn) SetDialect(dialect string) {
	x.dialect = dialect
}

// Dialect returns the dialect
func (x *OrderColumn) Dialect() string {
	return x.dialect
}

// Total returns the total value
func (x *OrderColumn) Total() int {
	return x.total
}

// SetTotal sets the totla arguments
func (x *OrderColumn) SetTotal(total int) {
	x.total = total
}

// Err returns the error
func (x *OrderColumn) Err() error {
	return x.err
}

// Equal returns true the expressions are equal; otherwise false.
func (x *OrderColumn) Equal(y *OrderColumn) bool {
	return strings.EqualFold(x.column, y.column) && strings.EqualFold(x.order, y.order)
}

// Query returns the order by clause
func (x *OrderColumn) Query() (string, []interface{}) {
	b := &Builder{dialect: x.dialect, total: x.total}
	b.Ident(x.column)

	switch x.order {
	case "asc":
		b.WriteString(" ASC")
	case "desc":
		b.WriteString(" DESC")
	default:
		b.WriteString(" ASC")
	}

	return b.String(), nil
}

// SetDialect sets the dialect
func (n Queries) SetDialect(dialect string) {
	type QuerierDialect interface {
		SetDialect(dialect string)
	}

	for _, querier := range n {
		if r, ok := querier.(QuerierDialect); ok {
			r.SetDialect(dialect)
		}
	}
}

// ReturningExpr adds the `RETURNING` clause to the insert statement. PostgreSQL only.
func (i *InsertBuilder) ReturningExpr(queriers ...Querier) *InsertBuilder {
	type Projection interface {
		Columns() []string
	}

	i.returning = []string{}

	for _, querier := range queriers {
		if projection, ok := querier.(Projection); ok {
			i.returning = append(i.returning, projection.Columns()...)
		} else {
			column, _ := querier.Query()
			// append the actual column
			i.returning = append(i.returning, column)
		}
	}

	return i
}

// Or sets the next coming predicate with OR operator (disjunction).
func (u *UpdateBuilder) Or() *UpdateBuilder {
	u.or = true
	return u
}

// Not sets the next coming predicate with not.
func (u *UpdateBuilder) Not() *UpdateBuilder {
	u.not = true
	return u
}

// Returning adds the `RETURNING` clause to the insert statement. PostgreSQL only.
func (u *UpdateBuilder) Returning(columns ...string) *UpdateBuilder {
	u.returning = columns
	return u
}

// ReturningExpr adds the `RETURNING` clause to the insert statement. PostgreSQL only.
func (u *UpdateBuilder) ReturningExpr(queriers ...Querier) *UpdateBuilder {
	type Projection interface {
		Columns() []string
	}

	u.returning = []string{}

	for _, querier := range queriers {
		if selection, ok := querier.(Projection); ok {
			u.returning = append(u.returning, selection.Columns()...)
		} else {
			column, _ := querier.Query()
			// append the actual column
			u.returning = append(u.returning, column)
		}
	}

	return u
}

// Or sets the next coming predicate with OR operator (disjunction).
func (d *DeleteBuilder) Or() *DeleteBuilder {
	d.or = true
	return d
}

// Not sets the next coming predicate with not.
func (d *DeleteBuilder) Not() *DeleteBuilder {
	d.not = true
	return d
}

// Returning adds the `RETURNING` clause to the insert statement. PostgreSQL only.
func (d *DeleteBuilder) Returning(columns ...string) *DeleteBuilder {
	d.returning = columns
	return d
}

// ReturningExpr adds the `RETURNING` clause to the insert statement. PostgreSQL only.
func (d *DeleteBuilder) ReturningExpr(queriers ...Querier) *DeleteBuilder {
	type Projection interface {
		Columns() []string
	}

	d.returning = []string{}

	for _, querier := range queriers {
		if selection, ok := querier.(Projection); ok {
			d.returning = append(d.returning, selection.Columns()...)
		} else {
			column, _ := querier.Query()
			// append the actual column
			d.returning = append(d.returning, column)
		}
	}

	return d
}
