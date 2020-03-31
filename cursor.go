package orm

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

// Position represents a column position within the cursor
type Position struct {
	Column string      `json:"column"`
	Order  string      `json:"order"`
	Value  interface{} `json:"value"`
}

// Equal returns true if the positions are equal
func (p *Position) Equal(pp *Position) bool {
	if !strings.EqualFold(p.Column, pp.Column) {
		return false
	}

	if !strings.EqualFold(p.Order, pp.Order) {
		return false
	}

	return true
}

// Cursor represents a cursor
type Cursor []*Position

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
