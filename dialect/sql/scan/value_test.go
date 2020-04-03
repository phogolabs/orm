package scan_test

import (
	"database/sql"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/google/uuid"
	"github.com/phogolabs/orm/dialect/sql/scan"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Values", func() {
	Context("when the source is struct", func() {
		user := &User{
			ID:        uuid.New(),
			Name:      "John Doe",
			Email:     pointer.ToString("john.doe@example.com"),
			DeletedAt: nil,
			CreatedAt: time.Now(),
		}

		It("scans the values successfully", func() {
			values, err := scan.Values(user, "name", "email")
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(HaveLen(2))
			Expect(values[0]).To(Equal(user.Name))
			Expect(values[1]).To(Equal(user.Email))
		})

		Context("when the column is not found", func() {
			It("scans the values successfully", func() {
				values, err := scan.Values(user, "name", "email", "address")
				Expect(err).NotTo(HaveOccurred())
				Expect(values).To(HaveLen(2))
				Expect(values[0]).To(Equal(user.Name))
				Expect(values[1]).To(Equal(user.Email))
			})
		})

		Context("when the columns are not provided", func() {
			It("scans the values successfully", func() {
				values, err := scan.Values(user)
				Expect(err).NotTo(HaveOccurred())
				Expect(values).To(HaveLen(5))
				Expect(values[0]).To(Equal(user.ID))
				Expect(values[1]).To(Equal(user.Name))
				Expect(values[2]).To(Equal(user.Email))
				Expect(values[3]).To(Equal(user.DeletedAt))
				Expect(values[4]).To(Equal(user.CreatedAt))
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
		kv := map[string]interface{}{
			"id":         uuid.New(),
			"name":       "John Doe",
			"email":      pointer.ToString("john.doe@example.com"),
			"deleted_at": nil,
			"created_at": time.Now(),
		}

		It("scans the values successfully", func() {
			values, err := scan.Values(&kv, "id", "name")
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(HaveLen(2))
			Expect(values[0]).To(Equal(kv["id"]))
			Expect(values[1]).To(Equal(kv["name"]))
		})

		Context("when the key is not found", func() {
			It("scans the values successfully", func() {
				values, err := scan.Values(&kv, "id", "address")
				Expect(err).NotTo(HaveOccurred())
				Expect(values).To(HaveLen(1))
				Expect(values[0]).To(Equal(kv["id"]))
			})
		})

		Context("when the columns are not provided", func() {
			It("scans the values successfully", func() {
				values, err := scan.Values(&kv)
				Expect(err).NotTo(HaveOccurred())
				Expect(values).To(HaveLen(5))
			})
		})

		Context("when the map key is not supported", func() {
			It("returns an error", func() {
				kv := map[int]interface{}{}
				values, err := scan.Values(&kv)
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
				sql.NamedArg{Name: "user", Value: "root"},
				sql.NamedArg{Name: "password", Value: "swordfish"},
			}

			values, err := scan.Args(args, "id", "user", "password")
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(HaveLen(3))
			Expect(values[0]).To(Equal(1))
			Expect(values[1]).To(Equal("root"))
			Expect(values[2]).To(Equal("swordfish"))
		})
	})
})
