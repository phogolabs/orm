package scan_test

import (
	"time"

	"github.com/phogolabs/orm/dialect/sql/scan"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Iter", func() {
	type Student struct {
		ID        string    `db:"id,primary_key"`
		Name      string    `db:"name"`
		CreatedAt time.Time `db:"created_at,read_only"`
	}

	It("iterates over the fields", func() {
		student := &Student{
			ID:        "001",
			Name:      "Jack",
			CreatedAt: time.Now(),
		}

		iter := scan.IteratorOf(student)

		Expect(iter.Next()).To(BeTrue())

		column := iter.Column()
		Expect(column).NotTo(BeNil())
		Expect(column.Name).To(Equal("id"))
		Expect(column.HasOption("primary_key")).To(BeTrue())

		Expect(iter.Next()).To(BeTrue())

		column = iter.Column()
		Expect(column).NotTo(BeNil())
		Expect(column.Name).To(Equal("name"))
		Expect(column.Options).To(BeEmpty())

		Expect(iter.Next()).To(BeTrue())
		column = iter.Column()
		Expect(column).NotTo(BeNil())
		Expect(column.Name).To(Equal("created_at"))
		Expect(column.HasOption("read_only")).To(BeTrue())

		Expect(iter.Next()).To(BeFalse())
	})

})
