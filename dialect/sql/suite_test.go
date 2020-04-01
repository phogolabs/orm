package sql_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/orm/dialect/sql"
)

func TestSQL(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SQL Suite")
}

func SetDialect(dialect string, querier sql.Querier) {
	type state interface {
		SetDialect(string)
	}

	if t, ok := querier.(state); ok {
		t.SetDialect(dialect)
	}
}
