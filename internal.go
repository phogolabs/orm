package orm

import (
	"reflect"
	"strings"
)

type (
	commandable interface {
		Name() string
	}

	translatable interface {
		SetDialect(string)
	}

	errorable interface {
		Err() error
	}
)

func nameOf(value reflect.Type) string {
	switch value.Kind() {
	case reflect.Ptr:
		return nameOf(value.Elem())
	case reflect.Struct:
		return strings.ToLower(value.Name())
	case reflect.Slice:
		return nameOf(value.Elem())
	case reflect.Array:
		return nameOf(value.Elem())
	case reflect.Map:
		return "map"
	default:
		return value.Name()
	}
}
