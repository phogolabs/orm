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
})
