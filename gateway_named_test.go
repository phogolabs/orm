package oak_test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	lk "github.com/ulule/loukoum"

	_ "github.com/mattn/go-sqlite3"
	"github.com/phogolabs/oak"
	"github.com/phogolabs/oak/fake"
	"github.com/phogolabs/parcello"
)

var _ = Describe("Gateway", func() {
	var db *oak.Gateway

	Describe("Open", func() {
		Context("when cannot open the database", func() {
			It("returns an error", func() {
				g, err := oak.Open("sqlite4", "/tmp/oak.db")
				Expect(g).To(BeNil())
				Expect(err).To(MatchError(`sql: unknown driver "sqlite4" (forgotten import?)`))
			})
		})
	})

	Describe("API", func() {
		type Person struct {
			FirstName string `db:"first_name"`
			LastName  string `db:"last_name"`
			Email     string `db:"email"`
		}

		BeforeEach(func() {
			var err error
			db, err = oak.Open("sqlite3", "/tmp/oak.db")
			Expect(err).To(BeNil())
			Expect(db.DriverName()).To(Equal("sqlite3"))

			buffer := &bytes.Buffer{}
			fmt.Fprintln(buffer, "CREATE TABLE users (")
			fmt.Fprintln(buffer, "  first_name text,")
			fmt.Fprintln(buffer, "  last_name text,")
			fmt.Fprintln(buffer, "  email text")
			fmt.Fprintln(buffer, ");")
			fmt.Fprintln(buffer)

			_, err = db.Exec(oak.SQL(buffer.String()))
			Expect(err).To(BeNil())

			param := oak.P{"first_name": "John", "last_name": "Doe", "email": "john@example.com"}
			_, err = db.NamedExec(oak.NamedSQL("INSERT INTO users VALUES(:first_name, :last_name, :email)", param))
			Expect(err).To(Succeed())
		})

		AfterEach(func() {
			_, err := db.Exec(oak.SQL("DROP TABLE users"))
			Expect(err).To(BeNil())
			Expect(db.Close()).To(Succeed())
		})

		Describe("Routine", func() {
			var script string

			BeforeEach(func() {
				script = fmt.Sprintf("%v", time.Now().UnixNano())
				buffer := bytes.NewBufferString(fmt.Sprintf("-- name: %v", script))
				fmt.Fprintln(buffer)
				fmt.Fprintln(buffer, "SELECT * FROM users")
				Expect(db.LoadRoutinesFromReader(buffer)).To(Succeed())
			})

			It("returns a named command", func() {
				stmt, err := db.NamedRoutine(script, oak.P{})
				Expect(stmt).NotTo(BeNil())
				Expect(err).To(BeNil())

				query, params := stmt.NamedQuery()
				Expect(query).To(Equal("SELECT * FROM users"))
				Expect(params).To(BeEmpty())
			})

			Context("when loading a whole directory", func() {
				BeforeEach(func() {
					buffer := bytes.NewBufferString(fmt.Sprintf("-- name: %v", "cmd"))
					fmt.Fprintln(buffer)
					fmt.Fprintln(buffer, "SELECT * FROM categories")

					content := buffer.Bytes()

					node := &parcello.Node{
						Name:    "script.sql",
						Content: &content,
						Mutex:   &sync.RWMutex{},
					}

					fileSystem := &fake.FileSystem{}
					fileSystem.OpenFileReturns(parcello.NewResourceFile(node), nil)

					fileSystem.WalkStub = func(dir string, fn filepath.WalkFunc) error {
						return fn(node.Name, &parcello.ResourceFileInfo{Node: node}, nil)
					}

					Expect(db.LoadRoutinesFromFileSystem(fileSystem)).To(Succeed())
				})

				It("returns a command", func() {
					stmt, err := db.NamedRoutine("cmd")
					Expect(err).To(BeNil())
					Expect(stmt).NotTo(BeNil())

					query, params := stmt.NamedQuery()
					Expect(query).To(Equal("SELECT * FROM categories"))
					Expect(params).To(BeEmpty())
				})
			})

			Context("when the named statement does not exits", func() {
				It("does not return a statement", func() {
					_, err := db.NamedRoutine("down", oak.P{})
					Expect(err).To(MatchError("query 'down' not found"))
				})
			})
		})

		Describe("NamedSelect", func() {
			It("executes a query successfully", func() {
				query := lk.Select("first_name", "last_name", "email").From("users")

				persons := []Person{}
				Expect(db.NamedSelect(&persons, query)).To(Succeed())
				Expect(persons).To(HaveLen(1))
				Expect(persons[0].FirstName).To(Equal("John"))
				Expect(persons[0].LastName).To(Equal("Doe"))
				Expect(persons[0].Email).To(Equal("john@example.com"))
			})

			Context("with context", func() {
				It("executes a query successfully", func() {
					query := lk.Select("first_name", "last_name", "email").From("users")
					persons := []Person{}
					Expect(db.NamedSelectContext(context.Background(), &persons, query)).To(Succeed())
					Expect(persons).To(HaveLen(1))
					Expect(persons[0].FirstName).To(Equal("John"))
					Expect(persons[0].LastName).To(Equal("Doe"))
					Expect(persons[0].Email).To(Equal("john@example.com"))
				})
			})

			Context("when the query fails", func() {
				It("returns an error", func() {
					query := lk.Select("name").From("categories")

					persons := []Person{}
					Expect(db.NamedSelect(&persons, query)).To(MatchError("no such table: categories"))
					Expect(persons).To(BeEmpty())
				})
			})

			Context("when an embedded statement is used", func() {
				It("executes a query successfully", func() {
					query := oak.NamedSQL("SELECT * FROM users WHERE first_name = :name", sql.Named("name", "John"))

					persons := []Person{}
					Expect(db.NamedSelect(&persons, query)).To(Succeed())
					Expect(persons).To(HaveLen(1))
					Expect(persons[0].FirstName).To(Equal("John"))
					Expect(persons[0].LastName).To(Equal("Doe"))
					Expect(persons[0].Email).To(Equal("john@example.com"))
				})

				Context("when the query does not exist", func() {
					It("returns an error", func() {
						query := oak.NamedSQL("SELECT * FROM categories")

						persons := []Person{}
						Expect(db.NamedSelect(&persons, query)).To(MatchError("no such table: categories"))
						Expect(persons).To(BeEmpty())
					})
				})
			})
		})

		Describe("NamedSelectOne", func() {
			It("executes a query successfully", func() {
				query := lk.Select("first_name", "last_name", "email").From("users")

				person := Person{}
				Expect(db.NamedSelectOne(&person, query)).To(Succeed())
			})

			Context("with context", func() {
				It("executes a query successfully", func() {
					query := lk.Select("first_name", "last_name", "email").From("users")

					person := Person{}
					Expect(db.NamedSelectOneContext(context.Background(), &person, query)).To(Succeed())
				})
			})

			Context("when the query fails", func() {
				It("returns an error", func() {
					query := lk.Select("name").From("categories")

					person := Person{}
					Expect(db.NamedSelectOne(&person, query)).To(MatchError("no such table: categories"))
				})
			})
		})

		Describe("Query", func() {
			It("executes a query successfully", func() {
				query := lk.Select("first_name", "last_name", "email").From("users")

				var (
					firstName string
					lastName  string
					email     string
				)

				rows, err := db.NamedQuery(query)
				Expect(err).To(BeNil())
				Expect(rows.Next()).To(BeTrue())

				Expect(rows.Scan(&firstName, &lastName, &email)).To(Succeed())
				Expect(firstName).To(Equal("John"))
				Expect(lastName).To(Equal("Doe"))
				Expect(email).To(Equal("john@example.com"))

				Expect(rows.Next()).To(BeFalse())
				Expect(rows.Close()).To(Succeed())
			})

			Context("with context", func() {
				It("executes a query successfully", func() {
					query := lk.Select("first_name", "last_name", "email").From("users")

					var (
						firstName string
						lastName  string
						email     string
					)

					rows, err := db.NamedQueryContext(context.Background(), query)
					Expect(err).To(BeNil())
					Expect(rows.Next()).To(BeTrue())

					Expect(rows.Scan(&firstName, &lastName, &email)).To(Succeed())
					Expect(firstName).To(Equal("John"))
					Expect(lastName).To(Equal("Doe"))
					Expect(email).To(Equal("john@example.com"))

					Expect(rows.Next()).To(BeFalse())
					Expect(rows.Close()).To(Succeed())
				})
			})

			Context("when the query fails", func() {
				It("returns an error", func() {
					query := lk.Select("name").From("categories")

					rows, err := db.NamedQuery(query)
					Expect(err).To(MatchError("no such table: categories"))
					Expect(rows).To(BeNil())
				})
			})
		})

		Describe("QueryRow", func() {
			It("executes a query successfully", func() {
				query := lk.Select("first_name", "last_name", "email").From("users")

				row, err := db.NamedQueryRow(query)
				Expect(err).To(BeNil())
				Expect(row).NotTo(BeNil())

				var (
					firstName string
					lastName  string
					email     string
				)

				Expect(row.Scan(&firstName, &lastName, &email)).To(Succeed())
				Expect(firstName).To(Equal("John"))
				Expect(lastName).To(Equal("Doe"))
				Expect(email).To(Equal("john@example.com"))
			})

			Context("with context", func() {
				It("executes a query successfully", func() {
					query := lk.Select("first_name", "last_name", "email").From("users")

					row, err := db.NamedQueryRowContext(context.Background(), query)
					Expect(err).To(BeNil())
					Expect(row).NotTo(BeNil())

					var (
						firstName string
						lastName  string
						email     string
					)

					Expect(row.Scan(&firstName, &lastName, &email)).To(Succeed())
					Expect(firstName).To(Equal("John"))
					Expect(lastName).To(Equal("Doe"))
					Expect(email).To(Equal("john@example.com"))
				})
			})

			Context("when the query fails", func() {
				It("returns an error", func() {
					query := lk.Select("name").From("categories")

					row, err := db.NamedQueryRow(query)
					Expect(err).To(MatchError("no such table: categories"))
					Expect(row).To(BeNil())
				})
			})
		})

		Describe("Exec", func() {
			It("executes a query successfully", func() {
				query := lk.Delete("users")

				_, err := db.Exec(query)
				Expect(err).To(Succeed())

				rows, err := db.NamedQuery(oak.NamedSQL("SELECT * FROM users"))
				Expect(err).To(BeNil())
				Expect(rows).NotTo(BeNil())
				Expect(rows.Next()).To(BeFalse())
				Expect(rows.Close()).To(Succeed())
			})

			Context("when context", func() {
				It("executes a query successfully", func() {
					query := lk.Delete("users")

					_, err := db.NamedExecContext(context.Background(), query)
					Expect(err).To(Succeed())

					rows, err := db.Query(oak.SQL("SELECT * FROM users"))
					Expect(err).To(BeNil())
					Expect(rows).NotTo(BeNil())
					Expect(rows.Next()).To(BeFalse())
					Expect(rows.Close()).To(Succeed())
				})
			})

			Context("when the query fails", func() {
				It("returns an error", func() {
					query := lk.Delete("categories")
					_, err := db.NamedExec(query)
					Expect(err).To(MatchError("no such table: categories"))
				})
			})
		})

		Describe("Tx", func() {
			var tx *oak.Tx

			BeforeEach(func() {
				var err error
				tx, err = db.Begin()
				Expect(err).To(Succeed())
			})

			Describe("Select", func() {
				It("executes a query successfully", func() {
					query := lk.Select("first_name", "last_name", "email").From("users")

					persons := []Person{}
					Expect(tx.NamedSelect(&persons, query)).To(Succeed())
					Expect(persons).To(HaveLen(1))
					Expect(persons[0].FirstName).To(Equal("John"))
					Expect(persons[0].LastName).To(Equal("Doe"))
					Expect(persons[0].Email).To(Equal("john@example.com"))
					Expect(tx.Commit()).To(Succeed())
				})

				Context("with context", func() {
					It("executes a query successfully", func() {
						query := lk.Select("first_name", "last_name", "email").From("users")

						persons := []Person{}
						Expect(tx.NamedSelectContext(context.Background(), &persons, query)).To(Succeed())
						Expect(persons).To(HaveLen(1))
						Expect(persons[0].FirstName).To(Equal("John"))
						Expect(persons[0].LastName).To(Equal("Doe"))
						Expect(persons[0].Email).To(Equal("john@example.com"))
						Expect(tx.Commit()).To(Succeed())
					})
				})
			})

			Describe("SelectOne", func() {
				It("executes a query successfully", func() {
					query := lk.Select("first_name", "last_name", "email").From("users")

					person := Person{}
					Expect(tx.NamedSelectOne(&person, query)).To(Succeed())
					Expect(tx.Commit()).To(Succeed())
				})

				Context("with context", func() {
					It("executes a query successfully", func() {
						query := lk.Select("first_name", "last_name", "email").From("users")

						person := Person{}
						Expect(tx.NamedSelectOneContext(context.Background(), &person, query)).To(Succeed())
						Expect(tx.Commit()).To(Succeed())
					})
				})
			})

			Describe("QueryRow", func() {
				It("executes a query successfully", func() {
					query := lk.Select("first_name", "last_name", "email").From("users")
					_, err := tx.NamedQueryRow(query)
					Expect(err).To(Succeed())
				})

				Context("with context", func() {
					It("executes a query successfully", func() {
						query := lk.Select("first_name", "last_name", "email").From("users")
						_, err := tx.NamedQueryRowContext(context.Background(), query)
						Expect(err).To(Succeed())
					})
				})
			})

			Describe("Exec", func() {
				It("executes a query successfully", func() {
					query := lk.Delete("users")

					_, err := tx.Exec(query)
					Expect(err).To(Succeed())

					rows, err := tx.NamedQuery(oak.NamedSQL("SELECT * FROM users"))
					Expect(err).To(BeNil())
					Expect(rows).NotTo(BeNil())
					Expect(rows.Next()).To(BeFalse())
					Expect(rows.Close()).To(Succeed())
					Expect(tx.Commit()).To(Succeed())
				})

				Context("with context", func() {
					It("executes a query successfully", func() {
						query := lk.Delete("users")

						_, err := tx.NamedExecContext(context.Background(), query)
						Expect(err).To(Succeed())

						rows, err := tx.NamedQueryContext(context.Background(), oak.NamedSQL("SELECT * FROM users"))
						Expect(err).To(BeNil())
						Expect(rows).NotTo(BeNil())
						Expect(rows.Next()).To(BeFalse())
						Expect(rows.Close()).To(Succeed())
						Expect(tx.Commit()).To(Succeed())
					})
				})
			})

			Describe("Rollback", func() {
				It("rollbacks the transaction successfully", func() {
					query := lk.Delete("users")

					_, err := tx.NamedExec(query)
					Expect(err).To(Succeed())
					Expect(tx.Rollback()).To(Succeed())
				})
			})
		})
	})
})
