package orm_test

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/orm"
	"github.com/phogolabs/parcello"
)

var _ = Describe("Migrate", func() {
	It("passes the migrations to underlying migration executor", func() {
		dir, err := ioutil.TempDir("", "orm_generator")
		Expect(err).To(BeNil())

		url := filepath.Join(dir, "orm.db")
		db, err := orm.Open("sqlite3", url)
		Expect(err).To(BeNil())
		Expect(db.Migrate(parcello.Dir(dir))).To(Succeed())
	})
})

var _ = Describe("UnmarshalRQLParam", func() {
	type Order struct {
		Price uint `rql:"filter"`
	}

	It("unmarshals the RQL param successfully", func() {
		query := map[string]interface{}{
			"filter": map[string]interface{}{
				"price": map[string]int{
					"$gt": 20,
					"$lt": 30,
				},
			},
		}

		data, err := json.Marshal(query)
		Expect(err).To(BeNil())

		param, err := orm.UnmarshalRQLParam(Order{}, data)
		Expect(err).To(BeNil())
		Expect(param.FilterExp).To(Equal("(price > ? AND price < ?)"))
		Expect(param.FilterArgs).To(HaveLen(2))
		Expect(param.FilterArgs[0]).To(BeNumerically("==", 20))
		Expect(param.FilterArgs[1]).To(BeNumerically("==", 30))
	})

	Context("when the model is not provided", func() {
		It("returns an error", func() {
			param, err := orm.UnmarshalRQLParam(nil, []byte{})
			Expect(param).To(BeNil())
			Expect(err).To(MatchError("rql: 'Model' is a required field"))
		})
	})

	Context("when the RQL param is invalid", func() {
		It("returns an error", func() {
			param, err := orm.UnmarshalRQLParam(Order{}, []byte("wrong"))
			Expect(param).To(BeNil())
			Expect(err).To(MatchError("decoding buffer to Query: parse error: syntax error near offset 0 of 'wrong'"))
		})
	})
})
