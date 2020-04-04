package scan_test

import (
	"reflect"
	"time"

	"github.com/phogolabs/orm/dialect/sql/scan"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Allocator", func() {
	type User struct {
		ID        string    `db:"id"`
		Name      string    `db:"name"`
		CreatedAt time.Time `db:"created_at"`
	}

	Describe("Allocate", func() {
		It("allocates new values", func() {
			var (
				id   = "007"
				name = "John"
			)

			columns := []string{"id", "name", "created_at"}
			allocator, err := scan.NewAllocator(reflect.TypeOf(&User{}), columns)
			Expect(err).NotTo(HaveOccurred())

			values := allocator.Allocate()
			Expect(values).To(HaveLen(3))

			values[0] = &id
			values[1] = &name

			value := allocator.Create(values)
			user, ok := value.Interface().(*User)
			Expect(ok).To(BeTrue())
			Expect(user.ID).To(Equal(id))
			Expect(user.Name).To(Equal(name))
		})
	})
})
