package oak_test

import (
	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/oak"
)

var _ = Describe("Stmt", func() {
	It("prepares the command correctly", func() {
		stmt := oak.SQL("SELECT * FROM users WHERE id = ?", 1)
		query, params := stmt.Query()
		Expect(query).To(Equal("SELECT * FROM users WHERE id = ?"))
		Expect(params).To(ContainElement(1))
	})
})

var _ = Describe("NamedStmt", func() {
	Context("when the parameter is struct", func() {
		It("return the query", func() {
			type ObjP struct {
				Id int `db:"id"`
			}
			stmt := oak.NamedSQL("SELECT * FROM users WHERE id = :id", &ObjP{Id: 1})
			query, params := stmt.NamedQuery()
			Expect(query).To(Equal("SELECT * FROM users WHERE id = :id"))
			Expect(params).To(HaveLen(1))
			Expect(params).To(HaveKeyWithValue("id", 1))
		})
	})

	Context("when the parameter is named", func() {
		It("return the query", func() {
			stmt := oak.NamedSQL("SELECT * FROM users WHERE id = :id", sql.Named("id", 1))
			query, params := stmt.NamedQuery()
			Expect(query).To(Equal("SELECT * FROM users WHERE id = :id"))
			Expect(params).To(HaveLen(1))
			Expect(params).To(HaveKeyWithValue("id", 1))
		})
	})
})
