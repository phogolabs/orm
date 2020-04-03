package scan_test

import (
	"reflect"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/phogolabs/orm/dialect/sql/scan"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Allocator", func() {
	Describe("Allocate", func() {
		It("allocates new values", func() {
			type User struct {
				ID        string    `db:"id"`
				Name      string    `db:"name"`
				CreatedAt time.Time `db:"created_at"`
			}

			columns := []string{"id", "name", "created_at"}
			allocator, err := scan.NewAllocator(reflect.TypeOf(&User{}), columns)
			Expect(err).NotTo(HaveOccurred())

			values := allocator.Allocate()
			Expect(values).To(HaveLen(3))

			value := allocator.Set("007", "John", time.Now())
			user, ok := value.Interface().(*User)
			Expect(ok).To(BeTrue())
			Expect(user.ID).To(Equal("007"))
			Expect(user.Name).To(Equal("John"))
		})
	})
})
