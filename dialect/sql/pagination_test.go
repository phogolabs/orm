package sql_test

import (
	"github.com/phogolabs/orm/dialect/sql"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PaginateBy", func() {
	var query *sql.Selector

	BeforeEach(func() {
		query = sql.Select().
			From(sql.Table("users")).
			Where(sql.Like("name", "john")).
			OrderExpr(sql.OrderColumnBy("name")).
			Limit(100)
	})

	It("returns a paginator", func() {
		paginator := query.PaginateBy()

		query, args := paginator.Query()
		Expect(query).To(Equal("SELECT * FROM `users` WHERE `name` LIKE ? ORDER BY `name` ASC LIMIT 101"))
		Expect(paginator.Cursor()).To(BeNil())
		Expect(args).To(HaveLen(1))
	})

	Context("when the cursor is provided", func() {
		var cursor *sql.Cursor

		BeforeEach(func() {
			cursor = &sql.Cursor{
				OrderBy: sql.OrderBy("name"),
				WhereAt: []interface{}{"john"},
			}
		})

		It("returns a paginator", func() {
			paginator := query.PaginateBy(cursor)
			query, args := paginator.Query()
			Expect(query).To(Equal("SELECT * FROM `users` WHERE `name` LIKE ? AND (`name` > ? OR `name` = ?) ORDER BY `name` ASC LIMIT 101"))
			Expect(paginator.Cursor()).To(Equal(cursor))
			Expect(args).To(HaveLen(3))
		})
	})

	Context("when the query is not sorted", func() {
		BeforeEach(func() {
			query = sql.Select().
				From(sql.Table("users")).
				Limit(100)
		})

		It("returns an error", func() {
			paginator := query.PaginateBy()
			Expect(paginator.Err()).To(MatchError("sql: query should have at least one order by clause"))
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

			paginator := query.Limit(2).PaginateBy()
			Expect(paginator.Scan(&users)).To(Succeed())

			cursor := paginator.Cursor()
			Expect(cursor).NotTo(BeNil())
			Expect(cursor.OrderBy).NotTo(BeNil())
			Expect(cursor.OrderBy.String()).To(Equal("name asc"))
			Expect(cursor.WhereAt).To(HaveLen(1))
			Expect(cursor.WhereAt).To(ContainElement("Brown"))

			Expect(users).To(HaveLen(2))
			Expect(users[0].Name).To(Equal("Mike"))
			Expect(users[1].Name).To(Equal("Peter"))
		})

		Context("when the target is not a slice", func() {
			It("returns an error", func() {
				user := &User{}
				paginator := query.Limit(2).PaginateBy()
				Expect(paginator.Scan(&user)).To(MatchError("sql: invalid type **sql_test.User. expect []interface{}"))
			})
		})
	})

	Describe("SetDialect", func() {
		It("sets the dialect", func() {
			paginator := query.PaginateBy()
			paginator.SetDialect("postgres")
			Expect(paginator.Dialect()).To(Equal("postgres"))
		})
	})

	Describe("Query", func() {
		It("returns the actual query", func() {
			paginator := query.OrderExpr(sql.OrderColumnBy("id")).PaginateBy()
			query, args := paginator.Query()
			Expect(query).To(Equal("SELECT * FROM `users` WHERE `name` LIKE ? ORDER BY `name` ASC, `id` ASC LIMIT 101"))
			Expect(args).To(HaveLen(1))
		})
	})
})
