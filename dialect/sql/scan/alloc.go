package scan

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx/reflectx"
)

// Allocator allocates values
type Allocator struct {
	types  []reflect.Type
	create func(values []interface{}) reflect.Value
}

// Create sets the given values
func (r *Allocator) Create(values []interface{}) reflect.Value {
	return r.create(values)
}

// Allocate allocates values
func (r *Allocator) Allocate() []interface{} {
	values := make([]interface{}, len(r.types))

	for index := range r.types {
		values[index] = reflect.New(r.types[index]).Interface()
	}

	return values
}

// Set sets the values
func (r *Allocator) Set(value, next reflect.Value, columns []string) {
	switch {
	case value.Kind() == reflect.Ptr:
		r.Set(value.Elem(), next.Elem(), columns)
	case value.Kind() == reflect.Struct:
		for _, name := range columns {
			field := fieldByName(value.Type(), name)
			// copy the value from the source to target
			source := next.FieldByIndex(field.Index)
			target := valueByIndex(value, field.Index)
			// set the value
			target.Set(source)
		}
	default:
		value.Set(next)
	}
}

// NewAllocator returns allocator  for the given reflect.Type.
func NewAllocator(target reflect.Type, columns []string) (*Allocator, error) {
	switch k := target.Kind(); {
	case k == reflect.Interface && target.NumMethod() == 0:
		fallthrough // interface{}
	case k == reflect.String || k >= reflect.Bool && k <= reflect.Float64:
		return NewAllocatorPrimitive(target), nil
	case k == reflect.Ptr:
		return NewAllocatorPtr(target, columns)
	case k == reflect.Struct:
		return NewAllocatorStruct(target, columns)
	default:
		return nil, fmt.Errorf("sql/scan: unsupported type ([]%s)", k)
	}
}

// NewAllocatorPrimitive allocates primitive type
func NewAllocatorPrimitive(typ reflect.Type) *Allocator {
	return &Allocator{
		types: []reflect.Type{typ},
		create: func(v []interface{}) reflect.Value {
			return reflect.Indirect(reflect.ValueOf(v[0]))
		},
	}
}

// NewAllocatorStruct returns the a configuration for scanning an sql.Row into a struct.
func NewAllocatorStruct(target reflect.Type, columns []string) (*Allocator, error) {
	var (
		types   = []reflect.Type{}
		indices = make([][]int, 0, target.NumField())
	)

	for _, name := range columns {
		name = strings.ToLower(strings.Split(name, "(")[0])

		field := fieldByName(target, name)
		// check if the field is nil
		if field == nil {
			return nil, fmt.Errorf("sql/scan: missing struct field for column: %s", name)
		}

		indices = append(indices, field.Index)
		types = append(types, field.Field.Type)
	}

	allocator := &Allocator{
		types: types,
		create: func(values []interface{}) reflect.Value {
			row := reflect.New(target).Elem()

			for index, value := range values {
				vector := indices[index]
				column := valueByIndex(row, vector)
				column.Set(reflect.Indirect(reflect.ValueOf(value)))
			}

			return row
		},
	}

	return allocator, nil
}

// NewAllocatorPtr wraps the underlying type with rowScan.
func NewAllocatorPtr(target reflect.Type, columns []string) (*Allocator, error) {
	target = target.Elem()

	allocator, err := NewAllocator(target, columns)
	if err != nil {
		return nil, err
	}

	create := allocator.create

	allocator.create = func(vs []interface{}) reflect.Value {
		value := create(vs)
		ptrTyp := reflect.PtrTo(value.Type())
		ptr := reflect.New(ptrTyp.Elem())
		ptr.Elem().Set(value)
		return ptr
	}

	return allocator, nil
}

func valueByIndex(target reflect.Value, vector []int) reflect.Value {
	if len(vector) == 1 {
		return target.Field(vector[0])
	}

	for depth, index := range vector {
		if depth > 0 && target.Kind() == reflect.Ptr {
			valType := target.Type().Elem()

			if valType.Kind() == reflect.Struct && target.IsNil() {
				// set the value
				target.Set(reflect.New(valType))
			}

			target = target.Elem()
		}

		// field
		target = target.Field(index)
	}

	return target
}

func fieldByName(target reflect.Type, name string) *reflectx.FieldInfo {
	meta := mapper.TypeMap(target)

	if field, ok := meta.Names[name]; ok {
		return field
	}

	find := func(parent *reflectx.FieldInfo, key string) *reflectx.FieldInfo {
		if field := fieldByName(parent.Field.Type, key); field != nil {
			// translate the field
			index := append(meta.Tree.Index, parent.Index...)
			index = append(index, field.Index...)
			// traverse
			return meta.GetByTraversal(index)
		}

		return nil
	}

	trim := func(parent *reflectx.FieldInfo, name string) string {
		name = strings.TrimPrefix(name, parent.Name)
		name = strings.TrimPrefix(name, "_")
		// done
		return name
	}

	for _, parent := range meta.Tree.Children {
		if key, ok := parent.Options["foreign_key"]; ok && key == name {
			if field, ok := parent.Options["reference_key"]; ok {
				return find(parent, field)
			}
		}

		if _, ok := parent.Options["prefix"]; ok {
			return find(parent, trim(parent, name))
		}
	}

	return nil
}
