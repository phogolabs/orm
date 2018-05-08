package oak_test

import (
	"io/ioutil"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/oak"
	"github.com/phogolabs/parcello"
)

var _ = Describe("Migrate", func() {
	It("passes the migrations to underlying migration executor", func() {
		dir, err := ioutil.TempDir("", "oak_generator")
		Expect(err).To(BeNil())

		url := filepath.Join(dir, "oak.db")
		db, err := oak.Open("sqlite3", url)
		Expect(err).To(BeNil())
		Expect(oak.Migrate(db, parcello.Dir(dir))).To(Succeed())
	})
})
