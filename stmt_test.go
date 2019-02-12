package orm_test

import (
	"database/sql"

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
		var param *orm.RQLParam

		BeforeEach(func() {
			param = &orm.RQLParam{
				Offset:     10,
				Limit:      10,
				Sort:       "age asc, name",
				FilterExp:  "name = ? AND age >= ?",
				FilterArgs: []interface{}{"John", 22},
			}
		})

		It("returns the query", func() {
			stmt := orm.RQL("users", param)
			query, params := stmt.NamedQuery()
			Expect(query).To(Equal("SELECT * FROM users WHERE name = :arg1 AND age >= :arg2 ORDER BY age asc, name LIMIT 10 OFFSET 10"))
			Expect(params).To(HaveKeyWithValue("arg1", "John"))
			Expect(params).To(HaveKeyWithValue("arg2", 22))
		})

		Context("when the offset is not provided", func() {
			BeforeEach(func() {
				param.Offset = 0
			})

			It("returns the query", func() {
				stmt := orm.RQL("users", param)
				query, params := stmt.NamedQuery()
				Expect(query).To(Equal("SELECT * FROM users WHERE name = :arg1 AND age >= :arg2 ORDER BY age asc, name LIMIT 10"))
				Expect(params).To(HaveKeyWithValue("arg1", "John"))
				Expect(params).To(HaveKeyWithValue("arg2", 22))
			})
		})

		Context("when the limit is not provided", func() {
			BeforeEach(func() {
				param.Limit = 0
			})

			It("returns the query", func() {
				stmt := orm.RQL("users", param)
				query, params := stmt.NamedQuery()
				Expect(query).To(Equal("SELECT * FROM users WHERE name = :arg1 AND age >= :arg2 ORDER BY age asc, name OFFSET 10"))
				Expect(params).To(HaveKeyWithValue("arg1", "John"))
				Expect(params).To(HaveKeyWithValue("arg2", 22))
			})
		})

		Context("when the order is not provided", func() {
			BeforeEach(func() {
				param.Sort = ""
			})

			It("returns the query", func() {
				stmt := orm.RQL("users", param)
				query, params := stmt.NamedQuery()
				Expect(query).To(Equal("SELECT * FROM users WHERE name = :arg1 AND age >= :arg2 LIMIT 10 OFFSET 10"))
				Expect(params).To(HaveKeyWithValue("arg1", "John"))
				Expect(params).To(HaveKeyWithValue("arg2", 22))
			})
		})

		Context("when the filter is not provided", func() {
			BeforeEach(func() {
				param.FilterExp = ""
			})

			It("returns the query", func() {
				stmt := orm.RQL("users", param)
				query, params := stmt.NamedQuery()
				Expect(query).To(Equal("SELECT * FROM users ORDER BY age asc, name LIMIT 10 OFFSET 10"))
				Expect(params).To(HaveKeyWithValue("arg1", "John"))
				Expect(params).To(HaveKeyWithValue("arg2", 22))
			})
		})
	})
})
