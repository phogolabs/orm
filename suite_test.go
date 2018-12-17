package orm_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOAK(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OAK Suite")
}

type ObjP struct {
	Id int `db:"id"`
}

type ObjM struct {
	Id int
}

func (m *ObjM) Map() map[string]interface{} {
	param := make(map[string]interface{})
	param["id"] = m.Id
	param["name"] = "jack"
	return param
}
