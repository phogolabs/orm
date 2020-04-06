package sql

import (
	"fmt"
	"strings"
)

// Order represents a order
type Order struct {
	Column    string `json:"column"`
	Direction string `json:"direction"`
}

// OrderOf parses the given parts as asc and desc clauses
func (selector *Selector) OrderOf(orderBy ...*Order) *Selector {
	for _, order := range orderBy {
		if order == nil {
			continue
		}

		selector = selector.OrderBy(order.String())
	}

	return selector
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
