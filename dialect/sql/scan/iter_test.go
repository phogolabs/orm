package scan_test

import (
	"time"

	"github.com/phogolabs/orm/dialect/sql/scan"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Iter", func() {
	type Address struct {
		Street  string `db:"street"`
		City    string `db:"city"`
		State   string `db:"state"`
		Country string `db:"country"`
	}

	type Student struct {
		ID        string    `db:"id,primary_key"`
		Name      string    `db:"name"`
		Address   *Address  `db:"address,inline,prefix"`
		CreatedAt time.Time `db:"created_at,read_only"`
	}

	It("iterates over the fields", func() {
		student := &Student{
			ID:        "001",
			Name:      "Jack",
			CreatedAt: time.Now(),
			Address: &Address{
				Street:  "5th Avenue",
				City:    "New York City",
				State:   "New York",
				Country: "USA",
			},
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
		Expect(column.Name).To(Equal("address_street"))
		Expect(iter.Next()).To(BeTrue())

		column = iter.Column()
		Expect(column).NotTo(BeNil())
		Expect(column.Name).To(Equal("address_city"))
		Expect(iter.Next()).To(BeTrue())

		column = iter.Column()
		Expect(column).NotTo(BeNil())
		Expect(column.Name).To(Equal("address_state"))
		Expect(iter.Next()).To(BeTrue())

		column = iter.Column()
		Expect(column).NotTo(BeNil())
		Expect(column.Name).To(Equal("address_country"))
		Expect(iter.Next()).To(BeTrue())

		column = iter.Column()
		Expect(column).NotTo(BeNil())
		Expect(column.Name).To(Equal("created_at"))
		Expect(column.HasOption("read_only")).To(BeTrue())
		Expect(iter.Next()).To(BeFalse())
	})
})
