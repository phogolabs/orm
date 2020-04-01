package sql_test

import (
	"github.com/phogolabs/orm/dialect/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Routine", func() {
	It("creates new routine successfully", func() {
		type Namer interface {
			Name() string
		}

		routine := sql.Routine("my-test-routine", 5432)

		namer, ok := routine.(Namer)
		Expect(ok).To(BeTrue())
		Expect(namer.Name()).To(Equal("my-test-routine"))

		query, params := routine.Query()
		Expect(query).To(Equal("my-test-routine"))
		Expect(params).To(HaveLen(1))
		Expect(params).To(ContainElement(5432))
	})
})

var _ = Describe("Command", func() {
	Context("when the provided argument is single", func() {
		It("creates new command successfully", func() {
			routine := sql.Command("SELECT * FROM users WHERE id = ?", 5432)

			query, params := routine.Query()
			Expect(query).To(Equal("SELECT * FROM users WHERE id = :arg0"))
			Expect(params).To(HaveLen(1))
			Expect(params).To(ContainElement(5432))
		})
	})

	Context("when the provided argument is a slice", func() {
		It("creates new command successfully", func() {
			routine := sql.Command("SELECT * FROM users WHERE id = ? AND category_id > ? AND category_name = ?", 1, 77, "fruits")
			query, params := routine.Query()

			Expect(query).To(Equal("SELECT * FROM users WHERE id = :arg0 AND category_id > :arg1 AND category_name = :arg2"))
			Expect(params).To(HaveLen(3))
			Expect(params).To(ContainElement(1))
			Expect(params).To(ContainElement(77))
			Expect(params).To(ContainElement("fruits"))
		})

		Context("when the query is named", func() {
		})
	})

	Context("when the provided argument is a map", func() {
		It("creates new command successfully", func() {
			param := map[string]interface{}{
				"category_id":   99,
				"category_name": "nuts",
				"id":            1234,
			}

			routine := sql.Command("SELECT * FROM users WHERE id = :id AND category_id > :category_id AND category_name = :category_name", param)
			query, params := routine.Query()

			Expect(query).To(Equal("SELECT * FROM users WHERE id = :id AND category_id > :category_id AND category_name = :category_name"))
			Expect(params).To(HaveLen(3))
			Expect(params).To(ContainElement(1234))
			Expect(params).To(ContainElement(99))
			Expect(params).To(ContainElement("nuts"))
		})
	})

	Context("when the provided argument is a struct", func() {
		type User struct {
			ID           int    `db:"id"`
			CategoryID   int    `db:"category_id"`
			CategoryName string `db:"category_name"`
		}

		It("creates new command successfully", func() {
			u := &User{
				ID:           1234,
				CategoryID:   99,
				CategoryName: "nuts",
			}

			routine := sql.Command("SELECT * FROM users WHERE id = :id AND category_id > :category_id AND category_name = :category_name", u)

			SetDialect("postgres", routine)

			query, params := routine.Query()

			Expect(query).To(Equal("SELECT * FROM users WHERE id = $1 AND category_id > $2 AND category_name = $3"))
			Expect(params).To(HaveLen(3))
			Expect(params).To(ContainElement(1234))
			Expect(params).To(ContainElement(99))
			Expect(params).To(ContainElement("nuts"))
		})
	})
})
