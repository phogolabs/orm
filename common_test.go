package orm_test

import (
	"io/ioutil"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/orm"
	"github.com/phogolabs/parcello"
)

var _ = Describe("Migrate", func() {
	It("passes the migrations to underlying migration executor", func() {
		dir, err := ioutil.TempDir("", "orm_generator")
		Expect(err).To(BeNil())

		url := filepath.Join(dir, "orm.db")
		db, err := orm.Open("sqlite3", url)
		Expect(err).To(BeNil())
		Expect(db.Migrate(parcello.Dir(dir))).To(Succeed())
	})
})
