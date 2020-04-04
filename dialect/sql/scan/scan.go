package scan

import (
	"database/sql"
	"fmt"
	"reflect"
)

// Scanner is the interface that wraps the
// three sql.Rows methods used for scanning.
type Scanner interface {
	Next() bool
	Scan(...interface{}) error
	Columns() ([]string, error)
}

// Row scans one row to the given value. It fails if the rows holds more than 1 row.
func Row(scanner Scanner, src interface{}) error {
	value, err := valueOf(src)
	if err != nil {
		return err
	}

	columns, err := scanner.Columns()
	if err != nil {
		return fmt.Errorf("sql/scan: failed getting column names: %v", err)
	}

	if !scanner.Next() {
		return sql.ErrNoRows
	}

	allocator, err := NewAllocator(value.Type(), columns)
	if err != nil {
		return err
	}

	if expected, actual := len(columns), len(allocator.types); expected > actual {
		return fmt.Errorf("sql/scan: columns do not match (%d > %d)", expected, actual)
	}

	values := allocator.Allocate()
	if err := scanner.Scan(values...); err != nil {
		return fmt.Errorf("sql/scan: failed scanning rows: %v", err)
	}

	next := allocator.Create(values)
	allocator.Set(value, next, columns)

	if scanner.Next() {
		return fmt.Errorf("sql/scan: expect exactly one row in result set")
	}

	return nil
}

// Rows scans the given ColumnScanner (basically, sql.Row or sql.Rows) into the given slice.
func Rows(scanner Scanner, src interface{}) error {
	value, err := valueOf(src)
	if err != nil {
		return err
	}

	columns, err := scanner.Columns()
	if err != nil {
		return fmt.Errorf("sql/scan: failed getting column names: %v", err)
	}

	if kind := value.Kind(); kind != reflect.Slice {
		return fmt.Errorf("sql/scan: invalid type %s. expected slice as an argument", kind)
	}

	allocator, err := NewAllocator(value.Type().Elem(), columns)
	if err != nil {
		return err
	}

	if expected, actual := len(columns), len(allocator.types); expected > actual {
		return fmt.Errorf("sql/scan: columns do not match (%d > %d)", expected, actual)
	}

	var (
		count = value.Len()
		index = 0
	)

	for scanner.Next() {
		values := allocator.Allocate()

		if err := scanner.Scan(values...); err != nil {
			return fmt.Errorf("sql/scan: failed scanning rows: %v", err)
		}

		switch {
		case index < count:
			allocator.Set(value.Index(index), allocator.Create(values), columns)
			index++
		default:
			next := reflect.Append(value, allocator.Create(values))
			value.Set(next)
		}
	}

	return nil
}
