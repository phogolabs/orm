package sql_test

import (
	"github.com/phogolabs/orm/dialect/sql"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/phogolabs/orm/dialect/sql/scan/mock"
)

var _ = Describe("Mutation", func() {
	var entity *User

	BeforeEach(func() {
		email := "jack@example.com"
		// mock it
		entity = &User{
			ID:    "007",
			Name:  "Jack",
			Email: &email,
			Group: &Group{
				ID:   "555",
				Name: "guest",
			},
		}
	})

	Describe("InsertMutation", func() {
		It("creates new insert mutation", func() {
			query, params := sql.NewInsert("users").Entity(entity).Query()

			Expect(query).To(Equal("INSERT INTO `users` (`id`, `name`, `email`, `group_id`) VALUES (?, ?, ?, ?)"))
			Expect(params).To(HaveLen(4))
			Expect(params[0]).To(Equal(entity.ID))
			Expect(params[1]).To(Equal(entity.Name))
			Expect(params[2]).To(Equal(entity.Email))
			Expect(params[3]).To(Equal(entity.Group.ID))
		})
	})

	Describe("UpdateMutation", func() {
		It("creates new update mutation", func() {
			query, params := sql.NewUpdate("users").Entity(entity, "id", "name", "email", "group_id").Query()

			Expect(query).To(Equal("UPDATE `users` SET `name` = ?, `email` = ?, `group_id` = ? WHERE `id` = ?"))
			Expect(params).To(HaveLen(4))
			Expect(params[3]).To(Equal(entity.ID))
			Expect(params[0]).To(Equal(entity.Name))
			Expect(params[1]).To(Equal(entity.Email))
			Expect(params[2]).To(Equal(entity.Group.ID))
		})

		Context("when the columns are not provided", func() {
			It("updates all columns", func() {
				query, params := sql.NewUpdate("users").Entity(entity).Query()

				Expect(query).To(Equal("UPDATE `users` SET `name` = ?, `email` = ?, `group_id` = ? WHERE `id` = ?"))
				Expect(params).To(HaveLen(4))
				Expect(params[3]).To(Equal(entity.ID))
				Expect(params[0]).To(Equal(entity.Name))
				Expect(params[1]).To(Equal(entity.Email))
				Expect(params[2]).To(Equal(entity.Group.ID))
			})
		})

		Context("when the value is nil", func() {
			BeforeEach(func() {
				entity.Email = nil
			})

			It("sets it to null", func() {
				query, params := sql.NewUpdate("users").Entity(entity, "email").Query()

				Expect(query).To(Equal("UPDATE `users` SET `email` = NULL WHERE `id` = ?"))
				Expect(params).To(HaveLen(1))
				Expect(params[0]).To(Equal(entity.ID))
			})
		})
	})

	Describe("DeleteMutation", func() {
		It("creates new delete mutation", func() {
			query, params := sql.NewDelete("users").Entity(entity).Query()

			Expect(query).To(Equal("DELETE FROM `users` WHERE `id` = ?"))
			Expect(params).To(HaveLen(1))
			Expect(params[0]).To(Equal(entity.ID))
		})
	})
})
