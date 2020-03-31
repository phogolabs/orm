package orm_test

import (
	"github.com/phogolabs/orm"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Paginator", func() {
	It("returns the query for the first page", func() {
		paginator := orm.Paginate("users").
			Order("+name,+id").
			Limit(200)

		query, params := paginator.NamedQuery()
		Expect(query).To(Equal("SELECT * FROM users ORDER BY name ASC, id ASC LIMIT 200"))
		Expect(params).To(BeEmpty())
	})

	It("returns the query for the second page", func() {
		cursor := &orm.Cursor{
			&orm.Position{Column: "name", Order: "+", Value: "Alpha"},
			&orm.Position{Column: "id", Order: "+", Value: "007"},
		}

		paginator := orm.Paginate("users").
			Start(cursor).
			Order("+name,+id").
			Limit(200)

		query, params := paginator.NamedQuery()
		Expect(query).To(Equal("SELECT * FROM users WHERE (name > :name OR (name = :name AND id > :id)) ORDER BY name ASC, id ASC LIMIT 200"))
		Expect(params).To(BeEmpty())
	})

	Context("when the where clause provided", func() {
		It("returns the query for the first page", func() {
			paginator := orm.Paginate("users").
				Where("name > ?", "John").
				Order("+name,+id").
				Limit(200)

			query, params := paginator.NamedQuery()
			Expect(query).To(Equal("SELECT * FROM users WHERE (name > :arg1) ORDER BY name ASC, id ASC LIMIT 200"))
			Expect(params).To(HaveKeyWithValue("arg1", "John"))
		})

		It("returns the query for the second page", func() {
			cursor := &orm.Cursor{
				&orm.Position{Column: "name", Order: "+", Value: "Alpha"},
				&orm.Position{Column: "id", Order: "+", Value: "007"},
			}

			paginator := orm.Paginate("users").
				Start(cursor).
				Where("name > ?", "John").
				Order("+name,+id").
				Limit(200)

			query, params := paginator.NamedQuery()
			Expect(query).To(Equal("SELECT * FROM users WHERE ((name > :name OR (name = :name AND id > :id)) AND (name > :arg1)) ORDER BY name ASC, id ASC LIMIT 200"))
			Expect(params).To(HaveKeyWithValue("arg1", "John"))
		})
	})

	Context("when the token is mismatched", func() {
		It("returns an error", func() {
			cursor := &orm.Cursor{
				&orm.Position{Column: "id", Order: "+", Value: "007"},
				&orm.Position{Column: "name", Order: "+", Value: "Alpha"},
			}

			paginator := orm.Paginate("users").
				Start(cursor).
				Order("+name,+id").
				Limit(200)

			query, params := paginator.NamedQuery()
			Expect(query).To(Equal("SELECT * FROM users WHERE (id > :id OR (id = :id AND name > :name)) ORDER BY id ASC, name ASC LIMIT 200"))
			Expect(params).To(BeEmpty())
			Expect(paginator.Error()).To(MatchError("orm: invalid cursor position"))
		})
	})
})
