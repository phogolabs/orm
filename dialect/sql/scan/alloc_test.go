package scan_test

import (
	"reflect"

	"github.com/phogolabs/orm/dialect/sql/scan"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/phogolabs/orm/dialect/sql/scan/mock"
)

var _ = Describe("Allocator", func() {
	Describe("Allocate", func() {
		It("allocates new values", func() {
			var (
				id    = "007"
				group = "555"
				name  = "John"
			)

			columns := []string{"id", "group_id", "name", "created_at"}
			allocator, err := scan.NewAllocator(reflect.TypeOf(&User{}), columns)
			Expect(err).NotTo(HaveOccurred())

			values := allocator.Allocate()
			Expect(values).To(HaveLen(4))

			values[0] = &id
			values[1] = &group
			values[2] = &name

			value := allocator.Create(values)
			user, ok := value.Interface().(*User)
			Expect(ok).To(BeTrue())
			Expect(user.ID).To(Equal(id))
			Expect(user.Name).To(Equal(name))
			Expect(user.Group).NotTo(BeNil())
			Expect(user.Group.ID).To(Equal(group))
		})
	})
})
