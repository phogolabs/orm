package sql_test

import (
	"encoding/json"

	"github.com/phogolabs/orm/dialect/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Paginator", func() {
	var selector *sql.Selector

	BeforeEach(func() {
		selector = sql.Select().From(sql.Table("users")).
			Where(sql.Like("name", "john")).
			OrderOf(sql.OrderBy{
				{Column: "name", Direction: "asc"},
			}).
			Limit(100)
	})

	It("creates new routine successfully", func() {
		paginator, err := selector.Clone().PaginateBy(sql.Asc("id")).Seek(nil)
		Expect(err).NotTo(HaveOccurred())

		query, _ := paginator.Query()
		Expect(query).To(Equal("SELECT * FROM `users` WHERE `name` LIKE ? ORDER BY `name` ASC, `id` ASC LIMIT ?"))

		params := make(map[string]interface{})
		params["id"] = 1
		params["name"] = "John"

		cursor, err := paginator.Cursor(&params)
		Expect(err).NotTo(HaveOccurred())

		paginator, err = selector.Clone().PaginateBy(sql.Asc("id")).Seek(cursor)
		Expect(err).NotTo(HaveOccurred())

		query, _ = paginator.Query()
		Expect(query).To(Equal("SELECT * FROM `users` WHERE `name` LIKE ? AND ((`name` > ?) OR ((`name` = ?) AND (`id` > ?))) ORDER BY `name` ASC, `id` ASC LIMIT ?"))

		data, err := json.Marshal(cursor)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(Equal(`[{"column":"name","order":"asc","value":"John"},{"column":"id","order":"asc","value":1}]`))
	})

	Describe("Seek", func() {
		Context("when the pagination column not provided", func() {
			It("returns an error", func() {
				paginator, err := selector.Clone().PaginateBy("").Seek(&sql.Cursor{})
				Expect(err).To(MatchError("sql: pagination column not provided"))
				Expect(paginator).To(BeNil())
			})
		})

		Context("when the cursor is mismatched", func() {
			It("returns an error", func() {
				cursor := &sql.Cursor{
					&sql.Vector{Column: "category", Order: "asc", Value: 1},
				}

				paginator, err := selector.Clone().PaginateBy(sql.Asc("id")).Seek(cursor)
				Expect(err).To(MatchError("sql: pagination cursor position mismatch"))
				Expect(paginator).To(BeNil())
			})

			It("returns an error", func() {
				cursor := &sql.Cursor{
					&sql.Vector{Column: "name", Order: "asc", Value: "John"},
				}

				paginator, err := selector.Clone().PaginateBy("id ASC").Seek(cursor)
				Expect(err).To(MatchError("sql: pagination column should be placed at the end"))
				Expect(paginator).To(BeNil())
			})
		})
	})

	Describe("Cursor", func() {
		It("returns the next cursor", func() {
			type User struct {
				ID   int    `db:"id"`
				Name string `db:"name"`
			}

			user := &User{ID: 1, Name: "John"}

			paginator, err := selector.Clone().PaginateBy(sql.Desc("id")).Seek(&sql.Cursor{})
			Expect(err).NotTo(HaveOccurred())

			cursor, err := paginator.Cursor(user)
			Expect(err).NotTo(HaveOccurred())

			positions := *cursor

			Expect(positions).To(HaveLen(2))
			Expect(positions[0].Column).To(Equal("name"))
			Expect(positions[0].Order).To(Equal("asc"))
			Expect(positions[0].Value).To(Equal("John"))

			Expect(positions[1].Column).To(Equal("id"))
			Expect(positions[1].Order).To(Equal("desc"))
			Expect(positions[1].Value).To(Equal(1))
		})

		Context("when the input is slice", func() {
			It("returns the next cursor", func() {
				type User struct {
					ID   int    `db:"id"`
					Name string `db:"name"`
				}

				u := []*User{
					{ID: 2, Name: "Brown"},
					{ID: 1, Name: "John"},
				}

				paginator, err := selector.Clone().PaginateBy(sql.Desc("id")).Seek(&sql.Cursor{})
				Expect(err).NotTo(HaveOccurred())

				cursor, err := paginator.Cursor(&u)
				Expect(err).NotTo(HaveOccurred())
				positions := *cursor

				Expect(positions).To(HaveLen(2))
				Expect(positions[0].Column).To(Equal("name"))
				Expect(positions[0].Order).To(Equal("asc"))
				Expect(positions[0].Value).To(Equal("John"))

				Expect(positions[1].Column).To(Equal("id"))
				Expect(positions[1].Order).To(Equal("desc"))
				Expect(positions[1].Value).To(Equal(1))
			})
		})

		Context("when the provided argument is not valid", func() {
			It("panics", func() {
				paginator, err := selector.Clone().PaginateBy(sql.Asc("id")).Seek(&sql.Cursor{})
				Expect(err).NotTo(HaveOccurred())
				Expect(func() { paginator.Cursor(nil) }).To(Panic())
			})
		})
	})

	Describe("SetDialect", func() {
		It("sets the dialect", func() {
			paginator, err := selector.Clone().PaginateBy(sql.Asc("id")).Seek(&sql.Cursor{})
			Expect(err).NotTo(HaveOccurred())
			paginator.SetDialect("postgres")
		})
	})

	Describe("Query", func() {
		It("returns the actual query", func() {
			paginator, err := selector.Clone().PaginateBy(sql.Asc("id")).Seek(&sql.Cursor{})
			Expect(err).NotTo(HaveOccurred())

			query, _ := paginator.Query()
			Expect(query).To(Equal("SELECT * FROM `users` WHERE `name` LIKE ? ORDER BY `name` ASC, `id` ASC LIMIT ?"))
		})
	})
})

var _ = Describe("Cursor", func() {
	Describe("DecodeCursor", func() {
		It("decodes a cursor successfully", func() {
			cursor, err := sql.DecodeCursor("W3siY29sdW1uIjoiaWQiLCJvcmRlciI6ImFzYyIsInZhbHVlIjoxfV0")
			Expect(err).NotTo(HaveOccurred())
			Expect(cursor).NotTo(BeNil())

			positions := *cursor
			Expect(positions).To(HaveLen(1))
			Expect(positions[0].Column).To(Equal("id"))
			Expect(positions[0].Order).To(Equal("asc"))
			Expect(positions[0].Value).To(BeNumerically("==", 1))
		})

		Describe("when the string is malformed", func() {
			It("returns an error", func() {
				cursor, err := sql.DecodeCursor("wrong")
				Expect(err).To(MatchError("illegal base64 data at input byte 5"))
				Expect(cursor).To(BeNil())
			})
		})

		Context("when the string is empty", func() {
			It("decodes a cursor successfully", func() {
				cursor, err := sql.DecodeCursor("")
				Expect(err).NotTo(HaveOccurred())
				Expect(cursor).NotTo(BeNil())
			})
		})
	})

	Describe("String", func() {
		It("returns a string", func() {
			cursor := &sql.Cursor{
				&sql.Vector{Column: "id", Order: "asc", Value: 1},
			}

			Expect(cursor.String()).To(Equal("W3siY29sdW1uIjoiaWQiLCJvcmRlciI6ImFzYyIsInZhbHVlIjoxfV0"))
		})

		Context("when the cursor is empty", func() {
			It("returns an empty string", func() {
				cursor := &sql.Cursor{}
				Expect(cursor.String()).To(BeEmpty())
			})
		})
	})
})
