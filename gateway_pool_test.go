package orm_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/orm"
	"github.com/phogolabs/parcello"
)

var _ = Describe("GatewayPool", func() {
	var pool *orm.GatewayPool

	BeforeEach(func() {
		pool = &orm.GatewayPool{
			URL:        os.Getenv("TEST_PSQL_URL"),
			Migrations: parcello.ManagerAt("migration"),
			Routines:   parcello.ManagerAt("routine"),
		}

		if pool.URL == "" {
			Skip("export TEST_PSQL_URL environment variable")
		}
	})

	AfterEach(func() {
		Expect(pool.Close()).To(Succeed())
	})

	Describe("Get", func() {
		It("returns a gateway successfully", func() {
			gateway, err := pool.Get("phogo")
			Expect(err).NotTo(HaveOccurred())
			Expect(gateway).NotTo(BeNil())

			users := []User{}
			Expect(gateway.Select(&users, orm.Routine("select-users"))).To(Succeed())

			gateway2, err := pool.Get("phogo")
			Expect(err).NotTo(HaveOccurred())
			Expect(gateway2).NotTo(BeNil())

			Expect(gateway).To(Equal(gateway2))
		})

		Context("when the name is empty", func() {
			It("returns an error", func() {
				gateway, err := pool.Get("")
				Expect(err).To(MatchError("orm: the provided key cannot be empty"))
				Expect(gateway).To(BeNil())
			})
		})

		Context("when the pool is isolated", func() {
			BeforeEach(func() {
				pool.Isolated = true
			})

			It("returns a gateway successfully", func() {
				gateway, err := pool.Get("phogo")
				Expect(err).NotTo(HaveOccurred())
				Expect(gateway).NotTo(BeNil())

				users := []User{}
				Expect(gateway.Select(&users, orm.SQL("SELECT * FROM phogo.users"))).To(Succeed())
			})

			Context("when cannot connect to the db", func() {
				BeforeEach(func() {
					pool.URL = "mongo://127.0.0.1:5432/foobar?sslmode=disable"
				})

				It("returns an error", func() {
					gateway, err := pool.Get("phogo")
					Expect(err).To(MatchError(`orm: name: phogo operation: parse_url error: not supported driver "mongo"`))
					Expect(gateway).To(BeNil())
				})
			})
		})

		Context("when the URL is invalid", func() {
			BeforeEach(func() {
				pool.URL = "://127.0.0.1:5432/foobar?sslmode=disable"
			})

			It("returns an error", func() {
				gateway, err := pool.Get("phogo")
				Expect(err).To(MatchError(`orm: name: phogo operation: connect error: parse ://127.0.0.1:5432/foobar?sslmode=disable: missing protocol scheme`))
				Expect(gateway).To(BeNil())
			})
		})
	})
})
