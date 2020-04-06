package sql_test

import (
	"github.com/phogolabs/orm/dialect/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Order", func() {
	Describe("OrderFrom", func() {
		It("returns the position from a given string", func() {
			position, err := sql.OrderFrom("ID ASC")
			Expect(err).NotTo(HaveOccurred())
			Expect(position).NotTo(BeNil())
			Expect(position.Column).To(Equal("id"))
			Expect(position.Direction).To(Equal("asc"))
		})

		Context("when the order is sanizied", func() {
			It("returns the position from a given string", func() {
				position, err := sql.OrderFrom(sql.Asc("ID"))
				Expect(err).NotTo(HaveOccurred())
				Expect(position).NotTo(BeNil())
				Expect(position.Column).To(Equal("id"))
				Expect(position.Direction).To(Equal("asc"))
			})
		})

		Context("when the order is not provided", func() {
			It("returns the position from a given string", func() {
				position, err := sql.OrderFrom("ID")
				Expect(err).NotTo(HaveOccurred())
				Expect(position).NotTo(BeNil())
				Expect(position.Column).To(Equal("id"))
				Expect(position.Direction).To(Equal("asc"))
			})
		})

		Context("when the string is empty", func() {
			It("returns the position from a given string", func() {
				position, err := sql.OrderFrom("")
				Expect(err).NotTo(HaveOccurred())
				Expect(position).To(BeNil())
			})
		})
	})

	Describe("Equal", func() {
		It("returns true", func() {
			position := &sql.Order{
				Column:    "id",
				Direction: "asc",
			}

			Expect(position.Equal(position)).To(BeTrue())
		})

		It("returns false", func() {
			position := &sql.Order{
				Column:    "id",
				Direction: "asc",
			}

			Expect(position.Equal(&sql.Order{})).To(BeFalse())
		})
	})

	Describe("String", func() {
		It("returns the position as a string", func() {
			position := &sql.Order{
				Column:    "id",
				Direction: "asc",
			}

			Expect(position.String()).To(Equal("`id` ASC"))
		})
	})
})
