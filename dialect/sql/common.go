package sql

import "strings"

// Unident return the string unidented
func Unident(v string) string {
	return strings.Replace(v, "`", "", -1)
}

// If cond return predicat
func If(cond bool, predicate *Predicate) *Predicate {
	if cond {
		return predicate
	}

	return nil
}
