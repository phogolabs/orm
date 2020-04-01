package sql_test

import (
	"github.com/phogolabs/orm/dialect/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeleteMutation", func() {
	It("creates new delete mutation", func() {
		type User struct {
			ID      string `db:"id,primary_key"`
			GroupID string `db:"group_id,primary_key"`
			Name    string `db:"name"`
		}

		user := &User{
			ID:      "007",
			GroupID: "guest",
			Name:    "Jack",
		}

		mutation := sql.NewDelete("users", user)

		query, params := mutation.Query()
		Expect(query).To(Equal("DELETE FROM `users` WHERE `id` = ? AND `group_id` = ?"))
		Expect(params).To(HaveLen(2))
		Expect(params[0]).To(Equal("007"))
		Expect(params[1]).To(Equal("guest"))
	})

	Context("when the struct does not have a primary key", func() {
		It("returns nil", func() {
			type User struct {
				ID      string `db:"id"`
				GroupID string `db:"group_id"`
				Name    string `db:"name"`
			}

			user := &User{
				ID:      "007",
				GroupID: "guest",
				Name:    "Jack",
			}

			mutation := sql.NewDelete("users", user)
			Expect(mutation).To(BeNil())
		})
	})
})

var _ = Describe("InsertMutation", func() {
	It("creates new insert mutation", func() {
		type User struct {
			ID      string `db:"id,primary_key"`
			GroupID string `db:"group_id,primary_key"`
			Name    string `db:"name"`
		}

		user := &User{
			ID:      "007",
			GroupID: "guest",
			Name:    "Jack",
		}

		mutation := sql.NewInsert("users", user)

		query, params := mutation.Query()
		Expect(query).To(Equal("INSERT INTO `users` (`id`, `group_id`, `name`) VALUES (?, ?, ?)"))
		Expect(params).To(HaveLen(3))
		Expect(params[0]).To(Equal("007"))
		Expect(params[1]).To(Equal("guest"))
		Expect(params[2]).To(Equal("Jack"))
	})
})

var _ = Describe("UpdateMutation", func() {
	It("creates new update mutation", func() {
		type User struct {
			ID      string `db:"id,primary_key,read_only"`
			GroupID string `db:"group_id,primary_key,read_only"`
			Name    string `db:"name"`
		}

		user := &User{
			ID:      "007",
			GroupID: "guest",
			Name:    "Jack",
		}

		mutation := sql.NewUpdate("users", user, "id", "group_id", "name")

		query, params := mutation.Query()
		Expect(query).To(Equal("UPDATE `users` SET `name` = ? WHERE `id` = ? AND `group_id` = ?"))
		Expect(params).To(HaveLen(3))
		Expect(params[0]).To(Equal("Jack"))
		Expect(params[1]).To(Equal("007"))
		Expect(params[2]).To(Equal("guest"))
	})

	Context("when the columns are not provided", func() {
		It("updates all columns", func() {
			type Student struct {
				ID      string `db:"id,primary_key,read_only"`
				GroupID string `db:"group_id"`
				Name    string `db:"name"`
			}

			student := &Student{
				ID:      "007",
				GroupID: "guest",
				Name:    "Jack",
			}

			mutation := sql.NewUpdate("students", student)

			query, params := mutation.Query()
			Expect(query).To(Equal("UPDATE `students` SET `group_id` = ?, `name` = ? WHERE `id` = ?"))
			Expect(params).To(HaveLen(3))
			Expect(params[0]).To(Equal("guest"))
			Expect(params[1]).To(Equal("Jack"))
			Expect(params[2]).To(Equal("007"))
		})
	})

	Context("when the value is nil", func() {
		It("sets it to null", func() {
			type Student struct {
				ID   string  `db:"id,primary_key,read_only"`
				Name *string `db:"name"`
			}

			student := &Student{
				ID: "007",
			}

			mutation := sql.NewUpdate("students", student)

			query, params := mutation.Query()
			Expect(query).To(Equal("UPDATE `students` SET `name` = NULL WHERE `id` = ?"))
			Expect(params).To(HaveLen(1))
			Expect(params[0]).To(Equal("007"))
		})
	})
})
