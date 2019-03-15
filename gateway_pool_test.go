package orm_test

import (
	"io/ioutil"
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
			URL:       os.Getenv("TEST_PSQL_URL"),
			Isolation: true,
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
		})

		Context("when cannot connect to the db", func() {
			BeforeEach(func() {
				pool.URL = "mongo://127.0.0.1:5432/foobar?sslmode=disable"
			})

			It("returns an error", func() {
				gateway, err := pool.Get("phogo")
				Expect(err).To(MatchError(`name: phogo error: not supported driver "mongo"`))
				Expect(gateway).To(BeNil())
			})
		})
	})

	Describe("Migrates", func() {
		It("migrates successfully", func() {
			dir, err := ioutil.TempDir("", "orm_generator")
			Expect(err).To(BeNil())

			Expect(pool.Migrate(parcello.Dir(dir), "phogo")).To(Succeed())
		})
	})
})
