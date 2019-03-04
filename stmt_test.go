package orm_test

import (
	"database/sql"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/orm"
)

var _ = Describe("Stmt", func() {
	Context("when the parameter is struct", func() {
		It("return the query", func() {
			stmt := orm.SQL("SELECT * FROM users WHERE id = :id", &ObjP{Id: 1})
			query, params := stmt.NamedQuery()
			Expect(query).To(Equal("SELECT * FROM users WHERE id = :id"))
			Expect(params).To(HaveLen(1))
			Expect(params).To(HaveKeyWithValue("id", 1))
		})
	})

	Context("when the parameter is param mapper", func() {
		It("return the query", func() {
			stmt := orm.SQL("SELECT * FROM users WHERE id = :id", &ObjM{Id: 1})
			query, params := stmt.NamedQuery()
			Expect(query).To(Equal("SELECT * FROM users WHERE id = :id"))
			Expect(params).To(HaveLen(2))
			Expect(params).To(HaveKeyWithValue("id", 1))
			Expect(params).To(HaveKeyWithValue("name", "jack"))
		})
	})

	Context("when the parameter is map", func() {
		It("return the query", func() {
			stmt := orm.SQL("SELECT * FROM users WHERE id = :id", map[string]interface{}{"id": 1})
			query, params := stmt.NamedQuery()
			Expect(query).To(Equal("SELECT * FROM users WHERE id = :id"))
			Expect(params).To(HaveLen(1))
			Expect(params).To(HaveKeyWithValue("id", 1))
		})
	})

	Context("when the parameter is named", func() {
		It("return the query", func() {
			stmt := orm.SQL("SELECT * FROM users WHERE id = :id", sql.Named("id", 1))
			query, params := stmt.NamedQuery()
			Expect(query).To(Equal("SELECT * FROM users WHERE id = :id"))
			Expect(params).To(HaveLen(1))
			Expect(params).To(HaveKeyWithValue("id", 1))
		})
	})

	Context("when the RQL is used", func() {
		var param *orm.RQLQuery

		type User struct {
			Name string `rql:"filter,sort"`
			Age  int    `rql:"filter,sort"`
		}

		prepare := func(stmt interface{}) {
			v := reflect.ValueOf(stmt)
			method := v.MethodByName("Prepare")

			p := reflect.ValueOf(&User{})
			values := method.Call([]reflect.Value{p})
			Expect(values).To(HaveLen(1))
			Expect(values[0].Interface()).NotTo(HaveOccurred())
		}

		BeforeEach(func() {
			param = &orm.RQLQuery{
				Offset: 10,
				Limit:  10,
				Filter: orm.Map{
					"age": orm.Map{
						"$gte": 22,
					},
				},
				Sort: []string{"+age", "-name"},
			}
		})

		It("returns the query", func() {
			stmt := orm.RQL("users", param)
			prepare(stmt)

			query, params := stmt.NamedQuery()
			Expect(query).To(Equal("SELECT * FROM users WHERE age >= :arg1 ORDER BY age asc, name desc LIMIT 10 OFFSET 10"))
			Expect(params).To(HaveKeyWithValue("arg1", 22))
		})

		Context("when the offset is not provided", func() {
			BeforeEach(func() {
				param.Offset = 0
			})

			It("returns the query", func() {
				stmt := orm.RQL("users", param)
				prepare(stmt)

				query, params := stmt.NamedQuery()
				Expect(query).To(Equal("SELECT * FROM users WHERE age >= :arg1 ORDER BY age asc, name desc LIMIT 10"))
				Expect(params).To(HaveKeyWithValue("arg1", 22))
			})
		})

		Context("when the limit is not provided", func() {
			BeforeEach(func() {
				param.Limit = 0
			})

			It("returns the query", func() {
				stmt := orm.RQL("users", param)
				prepare(stmt)

				query, params := stmt.NamedQuery()
				Expect(query).To(Equal("SELECT * FROM users WHERE age >= :arg1 ORDER BY age asc, name desc LIMIT 25 OFFSET 10"))
				Expect(params).To(HaveKeyWithValue("arg1", 22))
			})
		})

		Context("when the order is not provided", func() {
			BeforeEach(func() {
				param.Sort = []string{}
			})

			It("returns the query", func() {
				stmt := orm.RQL("users", param)
				prepare(stmt)

				query, params := stmt.NamedQuery()
				Expect(query).To(Equal("SELECT * FROM users WHERE age >= :arg1 LIMIT 10 OFFSET 10"))
				Expect(params).To(HaveKeyWithValue("arg1", 22))
			})
		})

		Context("when the filter is not provided", func() {
			BeforeEach(func() {
				param.Filter = nil
			})

			It("returns the query", func() {
				stmt := orm.RQL("users", param)
				prepare(stmt)

				query, params := stmt.NamedQuery()
				Expect(query).To(Equal("SELECT * FROM users ORDER BY age asc, name desc LIMIT 10 OFFSET 10"))
				Expect(params).To(HaveLen(0))
			})
		})
	})
})
