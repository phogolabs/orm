package orm_test

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	lk "github.com/ulule/loukoum"

	_ "github.com/mattn/go-sqlite3"
	"github.com/phogolabs/orm"
	"github.com/phogolabs/orm/fake"
	"github.com/phogolabs/parcello"
)

var _ = Describe("Gateway", func() {
	var db *orm.Gateway

	Describe("Open", func() {
		Context("when cannot open the database", func() {
			It("returns an error", func() {
				g, err := orm.Open("sqlite4", "/tmp/orm.db")
				Expect(g).To(BeNil())
				Expect(err).To(MatchError(`sql: unknown driver "sqlite4" (forgotten import?)`))
			})
		})
	})

	Describe("Connect", func() {
		It("opens the URL successfully", func() {
			g, err := orm.Connect("sqlite3://tmp/orm.db")
			Expect(g).NotTo(BeNil())
			Expect(err).To(BeNil())
			Expect(g.Close()).To(Succeed())
		})

		Context("when cannot open the database", func() {
			It("returns an error", func() {
				g, err := orm.Connect("://www")
				Expect(g).To(BeNil())
				Expect(err).To(MatchError("parse ://www: missing protocol scheme"))
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
			db, err = orm.Open("sqlite3", "/tmp/orm.db")
			Expect(err).To(BeNil())
			Expect(db.DriverName()).To(Equal("sqlite3"))

			buffer := &bytes.Buffer{}
			fmt.Fprintln(buffer, "CREATE TABLE users (")
			fmt.Fprintln(buffer, "  first_name text,")
			fmt.Fprintln(buffer, "  last_name text,")
			fmt.Fprintln(buffer, "  email text")
			fmt.Fprintln(buffer, ");")
			fmt.Fprintln(buffer)

			_, err = db.Exec(orm.SQL(buffer.String()))
			Expect(err).To(BeNil())

			_, err = db.Exec(orm.SQL("INSERT INTO users VALUES(?, ?, ?)", "John", "Doe", "john@example.com"))
			Expect(err).To(Succeed())
		})

		AfterEach(func() {
			_, err := db.Exec(orm.SQL("DROP TABLE users"))
			Expect(err).To(BeNil())
			Expect(db.Close()).To(Succeed())
		})

		Describe("Select", func() {
			It("executes a query successfully", func() {
				query := lk.Select("first_name", "last_name", "email").From("users")

				persons := []Person{}
				Expect(db.Select(&persons, query)).To(Succeed())
				Expect(persons).To(HaveLen(1))
				Expect(persons[0].FirstName).To(Equal("John"))
				Expect(persons[0].LastName).To(Equal("Doe"))
				Expect(persons[0].Email).To(Equal("john@example.com"))
			})

			Context("with context", func() {

				It("executes a query successfully", func() {
					query := lk.Select("first_name", "last_name", "email").From("users")
					persons := []Person{}
					Expect(db.SelectContext(context.Background(), &persons, query)).To(Succeed())
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
					Expect(db.Select(&persons, query)).To(MatchError("no such table: categories"))
					Expect(persons).To(BeEmpty())
				})
			})

			Context("when an embedded statement is used", func() {
				It("executes a query successfully", func() {
					query := orm.SQL("SELECT * FROM users WHERE first_name = ?", "John")

					persons := []Person{}
					Expect(db.Select(&persons, query)).To(Succeed())
					Expect(persons).To(HaveLen(1))
					Expect(persons[0].FirstName).To(Equal("John"))
					Expect(persons[0].LastName).To(Equal("Doe"))
					Expect(persons[0].Email).To(Equal("john@example.com"))
				})

				Context("when the query does not exist", func() {
					It("returns an error", func() {
						query := orm.SQL("SELECT * FROM categories")

						persons := []Person{}
						Expect(db.Select(&persons, query)).To(MatchError("no such table: categories"))
						Expect(persons).To(BeEmpty())
					})
				})
			})
		})

		Describe("SelectOne", func() {
			It("executes a query successfully", func() {
				query := lk.Select("first_name", "last_name", "email").From("users")

				person := Person{}
				Expect(db.SelectOne(&person, query)).To(Succeed())
			})

			Context("with context", func() {
				It("executes a query successfully", func() {
					query := lk.Select("first_name", "last_name", "email").From("users")

					person := Person{}
					Expect(db.SelectOneContext(context.Background(), &person, query)).To(Succeed())
				})
			})

			Context("when the query fails", func() {
				It("returns an error", func() {
					query := lk.Select("name").From("categories")

					person := Person{}
					Expect(db.SelectOne(&person, query)).To(MatchError("no such table: categories"))
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

				rows, err := db.Query(query)
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

					rows, err := db.QueryContext(context.Background(), query)
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

					rows, err := db.Query(query)
					Expect(err).To(MatchError("no such table: categories"))
					Expect(rows).To(BeNil())
				})
			})
		})

		Describe("QueryRow", func() {
			It("executes a query successfully", func() {
				query := lk.Select("first_name", "last_name", "email").From("users")

				row, err := db.QueryRow(query)
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

					row, err := db.QueryRowContext(context.Background(), query)
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

					row, err := db.QueryRow(query)
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

				rows, err := db.Query(orm.SQL("SELECT * FROM users"))
				Expect(err).To(BeNil())
				Expect(rows).NotTo(BeNil())
				Expect(rows.Next()).To(BeFalse())
				Expect(rows.Close()).To(Succeed())
			})

			Context("when context", func() {
				It("executes a query successfully", func() {
					query := lk.Delete("users")

					_, err := db.ExecContext(context.Background(), query)
					Expect(err).To(Succeed())

					rows, err := db.Query(orm.SQL("SELECT * FROM users"))
					Expect(err).To(BeNil())
					Expect(rows).NotTo(BeNil())
					Expect(rows.Next()).To(BeFalse())
					Expect(rows.Close()).To(Succeed())
				})
			})

			Context("when the query fails", func() {
				It("returns an error", func() {
					query := lk.Delete("categories")
					_, err := db.Exec(query)
					Expect(err).To(MatchError("no such table: categories"))
				})
			})

			Describe("Transaction", func() {
				Context("when the db is closed", func() {
					It("returns the error", func() {
						ddb, err := orm.Open("sqlite3", "/tmp/orm.db")
						Expect(err).To(BeNil())
						Expect(ddb.Close()).To(Succeed())

						terr := ddb.Transaction(func(tx *orm.Tx) error {
							return nil
						})

						Expect(terr).To(MatchError("sql: database is closed"))
					})
				})

				It("start the transaction successfully", func() {
					err := db.Transaction(func(tx *orm.Tx) error {
						_, err := tx.Exec(orm.SQL("DELETE FROM users"))
						Expect(err).NotTo(HaveOccurred())
						return nil

					})

					Expect(err).NotTo(HaveOccurred())

					rows, err := db.Query(orm.SQL("SELECT * FROM users"))
					Expect(err).To(BeNil())
					Expect(rows).NotTo(BeNil())
					Expect(rows.Next()).To(BeFalse())
					Expect(rows.Close()).To(Succeed())
				})

				Context("when the function returns an error", func() {
					It("rollbacks the transaction", func() {
						err := db.Transaction(func(tx *orm.Tx) error {
							return fmt.Errorf("Oh no!")
						})

						Expect(err).To(MatchError("Oh no!"))
					})
				})
			})
		})

		Describe("Tx", func() {
			var tx *orm.Tx

			BeforeEach(func() {
				var err error
				tx, err = db.Begin()
				Expect(err).To(Succeed())
			})

			Context("when the database is not available", func() {
				It("cannot start a transaction", func() {
					txDb, err := orm.Open("sqlite3", "/tmp/orm.db")
					Expect(err).To(BeNil())
					Expect(txDb.DriverName()).To(Equal("sqlite3"))
					Expect(txDb.Close()).To(Succeed())

					anotherTx, err := txDb.Begin()
					Expect(err).To(MatchError("sql: database is closed"))
					Expect(anotherTx).To(BeNil())
				})
			})

			Describe("Select", func() {
				It("executes a query successfully", func() {
					query := lk.Select("first_name", "last_name", "email").From("users")

					persons := []Person{}
					Expect(tx.Select(&persons, query)).To(Succeed())
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
						Expect(tx.SelectContext(context.Background(), &persons, query)).To(Succeed())
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
					Expect(tx.SelectOne(&person, query)).To(Succeed())
					Expect(tx.Commit()).To(Succeed())
				})

				Context("with context", func() {
					It("executes a query successfully", func() {
						query := lk.Select("first_name", "last_name", "email").From("users")

						person := Person{}
						Expect(tx.SelectOneContext(context.Background(), &person, query)).To(Succeed())
						Expect(tx.Commit()).To(Succeed())
					})
				})
			})

			Describe("QueryRow", func() {
				It("executes a query successfully", func() {
					query := lk.Select("first_name", "last_name", "email").From("users")
					_, err := tx.QueryRow(query)
					Expect(err).To(Succeed())
				})

				Context("with context", func() {
					It("executes a query successfully", func() {
						query := lk.Select("first_name", "last_name", "email").From("users")
						_, err := tx.QueryRowContext(context.Background(), query)
						Expect(err).To(Succeed())
					})
				})
			})

			Describe("Exec", func() {
				It("executes a query successfully", func() {
					query := lk.Delete("users")

					_, err := tx.Exec(query)
					Expect(err).To(Succeed())

					rows, err := tx.Query(orm.SQL("SELECT * FROM users"))
					Expect(err).To(BeNil())
					Expect(rows).NotTo(BeNil())
					Expect(rows.Next()).To(BeFalse())
					Expect(rows.Close()).To(Succeed())
					Expect(tx.Commit()).To(Succeed())
				})

				Context("with context", func() {
					It("executes a query successfully", func() {
						query := lk.Delete("users")

						_, err := tx.ExecContext(context.Background(), query)
						Expect(err).To(Succeed())

						rows, err := tx.QueryContext(context.Background(), orm.SQL("SELECT * FROM users"))
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

					_, err := tx.Exec(query)
					Expect(err).To(Succeed())
					Expect(tx.Rollback()).To(Succeed())
				})
			})
		})

		Context("when a routine is executed", func() {
			Context("when loading a singled file", func() {
				var script string

				BeforeEach(func() {
					script = fmt.Sprintf("%v", time.Now().UnixNano())
					buffer := bytes.NewBufferString(fmt.Sprintf("-- name: %v", script))
					fmt.Fprintln(buffer)
					fmt.Fprintln(buffer, "SELECT * FROM sqlite_master")
					_, err := db.ReadFrom(buffer)
					Expect(err).To(Succeed())
				})

				It("returns a command", func() {
					stmt := orm.Routine(script)
					query, params := stmt.NamedQuery()
					Expect(query).To(BeEmpty())
					Expect(params).To(BeEmpty())
				})

				It("executes the commands successfully", func() {
					stmt := orm.Routine(script)

					_, err := db.Exec(stmt)
					Expect(err).NotTo(HaveOccurred())

					query, params := stmt.NamedQuery()
					Expect(query).To(BeEmpty())
					Expect(params).To(BeEmpty())
				})
			})

			Context("when loading a whole directory", func() {
				BeforeEach(func() {
					buffer := bytes.NewBufferString(fmt.Sprintf("-- name: %v", "cmd"))
					fmt.Fprintln(buffer)
					fmt.Fprintln(buffer, "SELECT * FROM sqlite_master")

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

					Expect(db.ReadDir(fileSystem)).To(Succeed())
				})

				It("returns a command", func() {
					_, err := db.Exec(orm.Routine("cmd"))
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when the routine does not exits", func() {
				It("does not return a statement", func() {
					_, err := db.Exec(orm.Routine("down"))
					Expect(err).To(MatchError("query 'down' not found"))
				})
			})
		})
	})
})
