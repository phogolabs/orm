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
			entity := &User{
				ID:   "007",
				Name: "John Doe",
				Group: &Group{
					ID:          "555",
					Name:        "guest",
					Description: "My Group",
				},
			}

			columns := []string{"id", "group_id", "name", "created_at", "group_name", "group_description"}

			allocator, err := scan.NewAllocator(reflect.TypeOf(&User{}), columns)
			Expect(err).NotTo(HaveOccurred())

			values := allocator.Allocate()
			Expect(values).To(HaveLen(6))

			values[0] = &entity.ID
			values[1] = &entity.Group.ID
			values[2] = &entity.Name
			values[3] = &entity.CreatedAt
			values[4] = &entity.Group.Name
			values[5] = &entity.Group.Description

			record, ok := allocator.Create(values).Interface().(*User)
			Expect(ok).To(BeTrue())
			Expect(record).To(Equal(entity))
		})
	})
})
