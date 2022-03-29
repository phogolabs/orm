package scan_test

import (
	"time"

	"github.com/phogolabs/orm/dialect/sql/scan"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Iter", func() {
	type University struct {
		ID   string `db:"id,primary_key"`
		Name string `db:"name"`
	}

	type Student struct {
		ID         string      `db:"id,primary_key"`
		Name       string      `db:"name"`
		University *University `db:"university,foreign_key=university_id,reference_key=id"`
		CreatedAt  time.Time   `db:"created_at,read_only"`
	}

	It("iterates over the fields", func() {
		student := &Student{
			ID:   "001",
			Name: "Jack",
			University: &University{
				ID:   "007",
				Name: "Sofia University",
			},
			CreatedAt: time.Now(),
		}

		iter := scan.IteratorOf(student)
		Expect(iter.Next()).To(BeTrue())

		column := iter.Column()
		Expect(column).NotTo(BeNil())
		Expect(column.Name).To(Equal("id"))
		Expect(column.HasOption("primary_key")).To(BeTrue())
		Expect(iter.Value().Interface()).To(Equal("001"))
		Expect(iter.Next()).To(BeTrue())

		column = iter.Column()
		Expect(column).NotTo(BeNil())
		Expect(column.Name).To(Equal("name"))
		Expect(column.Options).To(BeEmpty())
		Expect(iter.Value().Interface()).To(Equal("Jack"))
		Expect(iter.Next()).To(BeTrue())

		column = iter.Column()
		Expect(column).NotTo(BeNil())
		Expect(column.Name).To(Equal("university_id"))
		Expect(column.HasOption("foreign_key")).To(BeTrue())
		Expect(column.HasOption("reference_key")).To(BeTrue())
		Expect(iter.Value().Interface()).To(Equal("007"))
		Expect(iter.Next()).To(BeTrue())

		column = iter.Column()
		Expect(column).NotTo(BeNil())
		Expect(column.Name).To(Equal("created_at"))
		Expect(column.HasOption("read_only")).To(BeTrue())
		Expect(iter.Next()).To(BeFalse())
	})
})
