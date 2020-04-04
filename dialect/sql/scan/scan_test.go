package scan_test

import (
	"database/sql"

	"github.com/phogolabs/orm/dialect/sql/scan"

	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scan", func() {
	var db *sql.DB

	type User struct {
		ID       int     `db:"id"`
		Name     *string `db:"name"`
		Password string  `db:"password"`
	}

	BeforeEach(func() {
		var err error

		db, err = sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
		Expect(err).To(BeNil())

		_, err = db.Exec("CREATE TABLE users (id int, name varchar(255), password varchar(10))")
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		_, err := db.Exec("DELETE FROM users")
		Expect(err).To(BeNil())

		Expect(db.Close()).To(Succeed())
	})

	Describe("Row", func() {
		It("scans the row successfully", func() {
			_, err := db.Exec("INSERT INTO users VALUES(1, 'root', 'swordfish')")
			Expect(err).To(BeNil())

			user := &User{}

			rows, err := db.Query("SELECT * FROM users")
			Expect(err).To(BeNil())

			Expect(scan.Row(rows, user)).To(Succeed())
			Expect(user.ID).To(Equal(1))
			Expect(*user.Name).To(Equal("root"))
			Expect(user.Password).To(Equal("swordfish"))
		})
	})

	Describe("Rows", func() {
		It("scans the row successfully", func() {
			_, err := db.Exec("INSERT INTO users VALUES(1, 'root', 'swordfish')")
			Expect(err).To(BeNil())

			_, err = db.Exec("INSERT INTO users VALUES(2, 'admin', 'qwerty')")
			Expect(err).To(BeNil())

			users := []*User{
				&User{ID: 1},
			}

			rows, err := db.Query("SELECT name,password FROM users")
			Expect(err).To(BeNil())

			Expect(scan.Rows(rows, &users)).To(Succeed())
			Expect(users).To(HaveLen(2))

			Expect(users[0].ID).To(Equal(1))
			Expect(*users[0].Name).To(Equal("root"))
			Expect(users[0].Password).To(Equal("swordfish"))

			Expect(users[1].ID).To(Equal(0))
			Expect(*users[1].Name).To(Equal("admin"))
			Expect(users[1].Password).To(Equal("qwerty"))
		})
	})
})
