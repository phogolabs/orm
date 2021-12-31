package sql_test

import (
	"github.com/phogolabs/orm/dialect/sql"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pager", func() {
	var query *sql.Selector

	BeforeEach(func() {
		query = sql.Select().
			From(sql.Table("users")).
			Where(sql.Like("name", "john")).
			OrderBy(sql.Asc("name")).
			Limit(100)
	})

	Describe("StartAt", func() {
		It("returns a pager", func() {
			pager := query.Pager()

			query, args := pager.Query()
			Expect(query).To(Equal("SELECT * FROM `users` WHERE `name` LIKE ? ORDER BY `name` ASC LIMIT ?"))
			Expect(pager.Token()).To(BeEmpty())
			Expect(args).To(HaveLen(2))
		})

		Context("when the token is provided", func() {
			It("returns a pager", func() {
				pager := query.StartAt("W3siYyI6Im5hbWUiLCJvIjoiYXNjIiwidiI6IkJyb3duIn1d")
				query, args := pager.Query()
				Expect(query).To(Equal("SELECT * FROM `users` WHERE `name` LIKE ? AND ((`name` > ?) OR (`name` = ?)) ORDER BY `name` ASC LIMIT ?"))
				Expect(pager.Token()).To(Equal("W3siYyI6Im5hbWUiLCJvIjoiYXNjIiwidiI6IkJyb3duIn1d"))
				Expect(args).To(HaveLen(4))
			})

		})

		Context("when the query is not sorted", func() {
			BeforeEach(func() {
				query = sql.Select().
					From(sql.Table("users")).
					Limit(100)
			})

			It("returns an error", func() {
				pager := query.Pager()
				Expect(pager.Error()).To(MatchError("sql: query should have at least one order by clause"))
			})
		})
	})

	Describe("Scan", func() {
		type User struct {
			ID   int    `db:"id"`
			Name string `db:"name"`
		}

		It("returns the next page token", func() {
			users := []*User{
				{ID: 4, Name: "Mike"},
				{ID: 3, Name: "Peter"},
				{ID: 2, Name: "Brown"},
			}

			pager := query.Limit(2).Pager()
			Expect(pager.Scan(&users)).To(Succeed())
			Expect(pager.Token()).To(Equal("W3siYyI6Im5hbWUiLCJvIjoiYXNjIiwidiI6IkJyb3duIn1d"))
			Expect(users).To(HaveLen(2))
			Expect(users[0].Name).To(Equal("Mike"))
			Expect(users[1].Name).To(Equal("Peter"))
		})

		Context("when the target is not a slice", func() {
			It("returns an error", func() {
				user := &User{}
				pager := query.Limit(2).Pager()
				Expect(pager.Scan(&user)).To(MatchError("dialect/sql: invalid type **sql_test.User. expect []interface{}"))
			})
		})
	})

	Describe("SetDialect", func() {
		It("sets the dialect", func() {
			pager := query.Pager()
			pager.SetDialect("postgres")
			Expect(pager.Dialect()).To(Equal("postgres"))
		})
	})

	Describe("Query", func() {
		It("returns the actual query", func() {
			pager := query.OrderBy("id").Pager()
			query, args := pager.Query()
			Expect(query).To(Equal("SELECT * FROM `users` WHERE `name` LIKE ? ORDER BY `name` ASC, `id` ASC LIMIT ?"))
			Expect(args).To(HaveLen(2))
		})
	})
})
