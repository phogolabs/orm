package scan_test

import (
	"github.com/phogolabs/orm/dialect/sql/scan"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NamedQuery", func() {
	Context("when the query uses questionmark parameters", func() {
		It("renames the query", func() {
			query, params := scan.NamedQuery("SELECT * FROM id = ? AND name LIKE %root% OR name = ?")
			Expect(query).To(Equal("SELECT * FROM id = :arg0 AND name LIKE %root% OR name = :arg1"))
			Expect(params).To(HaveLen(2))
			Expect(params[0]).To(Equal("arg0"))
			Expect(params[1]).To(Equal("arg1"))
		})
	})

	Context("when the query uses named parameters", func() {
		It("renames the query", func() {
			query, params := scan.NamedQuery("SELECT * FROM id = :arg0 AND name LIKE %root% OR name = :arg1")
			Expect(query).To(Equal("SELECT * FROM id = :arg0 AND name LIKE %root% OR name = :arg1"))
			Expect(params).To(HaveLen(2))
			Expect(params[0]).To(Equal("arg0"))
			Expect(params[1]).To(Equal("arg1"))
		})

		Context("when the args are not separated by spaces", func() {
			It("renames the query", func() {
				query, params := scan.NamedQuery("SELECT * FROM id IN (:arg0,:arg1)")
				Expect(query).To(Equal("SELECT * FROM id IN (:arg0,:arg1)"))
				Expect(params).To(HaveLen(2))
				Expect(params[0]).To(Equal("arg0"))
				Expect(params[1]).To(Equal("arg1"))
			})
		})

		Context("when the args have underline symbol in their name", func() {
			It("renames the query", func() {
				query, params := scan.NamedQuery("SELECT * FROM id IN (:category_id,:category_name)")
				Expect(query).To(Equal("SELECT * FROM id IN (:category_id,:category_name)"))
				Expect(params).To(HaveLen(2))
				Expect(params[0]).To(Equal("category_id"))
				Expect(params[1]).To(Equal("category_name"))
			})
		})
	})

	Context("when the query is mixed", func() {
		It("renames the query", func() {
			query, params := scan.NamedQuery("SELECT * FROM id IN (?,:name)")
			Expect(query).To(Equal("SELECT * FROM id IN (:arg0,:name)"))
			Expect(params).To(HaveLen(2))
			Expect(params[0]).To(Equal("arg0"))
			Expect(params[1]).To(Equal("name"))
		})
	})

	Context("when the arguments are duplicated", func() {
		It("renames the query", func() {
			query, params := scan.NamedQuery("SELECT * FROM id = :arg0 AND name LIKE :arg1 OR name = :arg1")
			Expect(query).To(Equal("SELECT * FROM id = :arg0 AND name LIKE :arg1 OR name = :arg1"))
			Expect(params).To(HaveLen(3))
			Expect(params[0]).To(Equal("arg0"))
			Expect(params[1]).To(Equal("arg1"))
			Expect(params[2]).To(Equal("arg1"))
		})
	})
})
