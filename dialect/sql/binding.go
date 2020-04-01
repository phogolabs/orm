package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/jmoiron/sqlx/reflectx"
)

var mapper = reflectx.NewMapper("db")

func bindParam(param interface{}) map[string]interface{} {
	value := reflect.ValueOf(param)

	switch {
	case value.Kind() == reflect.Ptr:
		return bindParam(value.Elem().Interface())
	case value.Kind() == reflect.Slice:
		return bindParamSlice(value)
	case value.Kind() == reflect.Map:
		return bindParamMap(value)
	case value.Kind() == reflect.Struct:
		return bindParamStruct(value)
	default:
		return map[string]interface{}{
			"?": param,
		}
	}
}

func bindParamStruct(value reflect.Value) map[string]interface{} {
	var (
		params = make(map[string]interface{})
		actual = value.Interface()
	)

	if arg, ok := actual.(sql.NamedArg); ok {
		params[arg.Name] = arg.Value
		return params
	}

	for k, v := range mapper.FieldMap(value) {
		k = strings.ToLower(k)
		params[k] = v.Interface()
	}

	return params
}

func bindParamSlice(value reflect.Value) map[string]interface{} {
	var (
		params = make(map[string]interface{})
		pindex = 0
	)

	for index := 0; index < value.Len(); index++ {
		kv := bindParam(value.Index(index).Interface())

		for k, v := range kv {
			if k == "?" {
				k = fmt.Sprintf("arg%d", pindex)
				pindex++
			}

			params[k] = v
		}
	}

	return params
}

func bindParamMap(value reflect.Value) map[string]interface{} {
	if kv, ok := value.Interface().(map[string]interface{}); ok {
		return kv
	}

	return make(map[string]interface{})
}

func compile(dialect, namedQ string) (query string, names []string, err error) {
	allowedBindRunes := []*unicode.RangeTable{
		unicode.Letter,
		unicode.Digit,
	}

	qs := []byte(namedQ)
	names = make([]string, 0, 10)
	rebound := make([]byte, 0, len(qs))

	inName := false
	last := len(qs) - 1
	currentVar := 1
	name := make([]byte, 0, 10)

	for i, b := range qs {
		// a ':' while we're in a name is an error
		if b == ':' {
			// if this is the second ':' in a '::' escape sequence, append a ':'
			if inName && i > 0 && qs[i-1] == ':' {
				rebound = append(rebound, ':')
				inName = false
				continue
			} else if inName {
				err = errors.New("unexpected `:` while reading named param at " + strconv.Itoa(i))
				return
			}
			inName = true
			name = []byte{}
		} else if inName && i > 0 && b == '=' && len(name) == 0 {
			rebound = append(rebound, ':', '=')
			inName = false
			continue
			// if we're in a name, and this is an allowed character, continue
		} else if inName && (unicode.IsOneOf(allowedBindRunes, rune(b)) || b == '_' || b == '.') && i != last {
			// append the byte to the name if we are in a name and not on the last byte
			name = append(name, b)
			// if we're in a name and it's not an allowed character, the name is done
		} else if inName {
			inName = false
			// if this is the final byte of the string and it is part of the name, then
			// make sure to add it to the name
			if i == last && unicode.IsOneOf(allowedBindRunes, rune(b)) {
				name = append(name, b)
			}
			// add the string representation to the names list
			names = append(names, string(name))
			// add a proper bindvar for the bindType
			switch dialect {
			// oracle only supports named type bind vars even for positional
			case "":
				rebound = append(rebound, ':')
				rebound = append(rebound, name...)
			case "mysql":
				rebound = append(rebound, '?')
			case "postgres", "postgresql":
				rebound = append(rebound, '$')
				for _, b := range strconv.Itoa(currentVar) {
					rebound = append(rebound, byte(b))
				}
				currentVar++
			}
			// add this byte to string unless it was not part of the name
			if i != last {
				rebound = append(rebound, b)
			} else if !unicode.IsOneOf(allowedBindRunes, rune(b)) {
				rebound = append(rebound, b)
			}
		} else {
			// this is a normal byte and should just go onto the rebound query
			rebound = append(rebound, b)
		}
	}

	return string(rebound), names, nil
}
