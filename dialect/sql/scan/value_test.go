package scan_test

import (
	"database/sql"
	"time"

	"github.com/phogolabs/orm/dialect/sql/scan"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/phogolabs/orm/dialect/sql/scan/mock"
)

var _ = Describe("Values", func() {
	var entity *User

	BeforeEach(func() {
		email := "john.doe@example.com"
		// mock it
		entity = &User{
			ID:        "007",
			Name:      "John Doe",
			Email:     &email,
			CreatedAt: time.Now(),
			Group: &Group{
				ID:          "555",
				Name:        "guest",
				Description: "My Group",
				CreatedAt:   time.Now(),
			},
		}
	})

	Context("when the source is struct", func() {
		It("scans the values successfully", func() {
			values, err := scan.Values(entity, "name", "email", "group_id")
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(HaveLen(3))
			Expect(values[0]).To(Equal(entity.Name))
			Expect(values[1]).To(Equal(entity.Email))
			Expect(values[2]).To(Equal(entity.Group.ID))
		})

		Context("when the column is path", func() {
			It("scans the values successfully", func() {
				values, err := scan.Values(entity, "group_id", "group_name", "group_description")
				Expect(err).NotTo(HaveOccurred())
				Expect(values).To(HaveLen(3))
				Expect(values[0]).To(Equal(entity.Group.ID))
				Expect(values[1]).To(Equal(entity.Group.Name))
				Expect(values[2]).To(Equal(entity.Group.Description))
			})
		})

		Context("when the column is not found", func() {
			It("scans the values successfully", func() {
				values, err := scan.Values(entity, "name", "email", "description")
				Expect(err).NotTo(HaveOccurred())
				Expect(values).To(HaveLen(2))
				Expect(values[0]).To(Equal(entity.Name))
				Expect(values[1]).To(Equal(entity.Email))
			})
		})

		Context("when the columns are not provided", func() {
			It("scans the values successfully", func() {
				values, err := scan.Values(entity)
				Expect(err).NotTo(HaveOccurred())
				Expect(values).To(HaveLen(6))
				Expect(values[0]).To(Equal(entity.ID))
				Expect(values[1]).To(Equal(entity.Name))
				Expect(values[2]).To(Equal(entity.Email))
				Expect(values[3]).To(Equal(entity.Group.ID))
				Expect(values[4]).To(Equal(entity.CreatedAt))
			})
		})
	})

	Context("when the source is NamedArg", func() {
		It("scans the values successfully", func() {
			arg := sql.Named("name", "John Doe")
			values, err := scan.Values(&arg, "name")
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(HaveLen(1))
			Expect(values[0]).To(Equal(arg.Value))
		})

		Context("when the columns are not provided", func() {
			It("scans the values successfully", func() {
				arg := sql.Named("name", "John Doe")
				values, err := scan.Values(&arg)
				Expect(err).NotTo(HaveOccurred())
				Expect(values).To(HaveLen(1))
				Expect(values[0]).To(Equal(arg.Value))
			})
		})

		Context("when the column is not found", func() {
			It("scans the values successfully", func() {
				arg := sql.Named("name", "John Doe")
				values, err := scan.Values(&arg, "age")
				Expect(err).NotTo(HaveOccurred())
				Expect(values).To(HaveLen(0))
			})
		})
	})

	Context("when the source is a map", func() {
		var dictionary map[string]interface{}

		BeforeEach(func() {
			dictionary = map[string]interface{}{
				"id":         "007",
				"name":       "John Doe",
				"email":      entity.Email,
				"created_at": time.Now(),
				"deleted_at": nil,
			}
		})

		It("scans the values successfully", func() {
			values, err := scan.Values(&dictionary, "id", "name")
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(HaveLen(2))
			Expect(values[0]).To(Equal(dictionary["id"]))
			Expect(values[1]).To(Equal(dictionary["name"]))
		})

		Context("when the key is not found", func() {
			It("scans the values successfully", func() {
				values, err := scan.Values(&dictionary, "id", "address")
				Expect(err).NotTo(HaveOccurred())
				Expect(values).To(HaveLen(1))
				Expect(values[0]).To(Equal(dictionary["id"]))
			})
		})

		Context("when the columns are not provided", func() {
			It("scans the values successfully", func() {
				values, err := scan.Values(&dictionary)
				Expect(err).NotTo(HaveOccurred())
				Expect(values).To(HaveLen(5))
			})
		})

		Context("when the map key is not supported", func() {
			It("returns an error", func() {
				dictionary := map[int]interface{}{}
				values, err := scan.Values(&dictionary)
				Expect(err).To(MatchError("sql/scan: invalid type int. expected string as an key"))
				Expect(values).To(BeEmpty())
			})
		})
	})

	Context("when the source is not supported", func() {
		It("returns an error", func() {
			src := "pointer"
			values, err := scan.Values(&src)
			Expect(err).To(MatchError("sql/scan: invalid type string. expected struct or map as an argument"))
			Expect(values).To(BeEmpty())
		})
	})

	Context("when the source is not pointer", func() {
		It("returns an error", func() {
			values, err := scan.Values("not pointer")
			Expect(err).To(MatchError("sql/scan: invalid type string. expected pointer as an argument"))
			Expect(values).To(BeEmpty())
		})
	})
})

var _ = Describe("Args", func() {
	Context("when no columns are provided", func() {
		It("returns the actual arguments", func() {
			args := []interface{}{1, "root", time.Now()}
			values, err := scan.Args(args)
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(Equal(args))
		})
	})

	Context("when the columns are provided", func() {
		It("returns the actual arguments", func() {
			args := []interface{}{
				sql.NamedArg{Name: "id", Value: 1},
				sql.NamedArg{Name: "entity", Value: "root"},
				sql.NamedArg{Name: "password", Value: "swordfish"},
			}

			values, err := scan.Args(args, "id", "entity", "password")
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(HaveLen(3))
			Expect(values[0]).To(Equal(1))
			Expect(values[1]).To(Equal("root"))
			Expect(values[2]).To(Equal("swordfish"))
		})
	})
})
