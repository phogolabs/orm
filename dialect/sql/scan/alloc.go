package scan

import (
	"fmt"
	"reflect"
	"strings"
)

// Allocator allocates values
type Allocator struct {
	types  []reflect.Type
	setter func(values []interface{}) reflect.Value
}

// Set sets the given values
func (r *Allocator) Set(values ...interface{}) reflect.Value {
	if r.setter != nil {
		return r.setter(values)
	}

	return reflect.Value{}
}

// Allocate allocates values
func (r *Allocator) Allocate() []interface{} {
	values := make([]interface{}, len(r.types))

	for i := range r.types {
		values[i] = reflect.New(r.types[i]).Interface()
	}

	return values
}

// NewAllocator returns allocator  for the given reflect.Type.
func NewAllocator(typ reflect.Type, columns []string) (*Allocator, error) {
	switch k := typ.Kind(); {
	case k == reflect.Interface && typ.NumMethod() == 0:
		fallthrough // interface{}
	case k == reflect.String || k >= reflect.Bool && k <= reflect.Float64:
		return NewAllocatorPrimitive(typ), nil
	case k == reflect.Ptr:
		return NewAllocatorPtr(typ, columns)
	case k == reflect.Struct:
		return NewAllocatorStruct(typ, columns)
	default:
		return nil, fmt.Errorf("sql/scan: unsupported type ([]%s)", k)
	}
}

// NewAllocatorPrimitive allocates primitive type
func NewAllocatorPrimitive(typ reflect.Type) *Allocator {
	return &Allocator{
		types: []reflect.Type{typ},
		setter: func(v []interface{}) reflect.Value {
			return reflect.Indirect(reflect.ValueOf(v[0]))
		},
	}
}

// NewAllocatorStruct returns the a configuration for scanning an sql.Row into a struct.
func NewAllocatorStruct(typ reflect.Type, columns []string) (*Allocator, error) {
	var (
		allocator = &Allocator{}
		meta      = mapper.TypeMap(typ)
		indices   = make([][]int, 0, typ.NumField())
	)

	for _, name := range columns {
		name = strings.ToLower(strings.Split(name, "(")[0])
		field, ok := meta.Names[name]

		if !ok {
			return nil, fmt.Errorf("sql/scan: missing struct field for column: %s", name)
		}

		indices = append(indices, field.Index)
		allocator.types = append(allocator.types, field.Field.Type)
	}

	allocator.setter = func(values []interface{}) reflect.Value {
		value := reflect.New(typ).Elem()

		for index, field := range values {
			value.FieldByIndex(indices[index]).Set(reflect.Indirect(reflect.ValueOf(field)))
		}

		return value
	}

	return allocator, nil
}

// NewAllocatorPtr wraps the underlying type with rowScan.
func NewAllocatorPtr(typ reflect.Type, columns []string) (*Allocator, error) {
	typ = typ.Elem()

	allocator, err := NewAllocator(typ, columns)
	if err != nil {
		return nil, err
	}

	wrap := allocator.setter

	allocator.setter = func(vs []interface{}) reflect.Value {
		value := wrap(vs)
		ptrTyp := reflect.PtrTo(value.Type())
		ptr := reflect.New(ptrTyp.Elem())
		ptr.Elem().Set(value)
		return ptr
	}

	return allocator, nil
}
