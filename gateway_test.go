package orm_test

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/phogolabs/orm"
	"github.com/phogolabs/orm/dialect/sql"
	"github.com/phogolabs/orm/fake"
	"github.com/phogolabs/parcello"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Open", func() {
	Context("when cannot open the database", func() {
		It("returns an error", func() {
			gateway, err := orm.Open("sqlite4", "/tmp/orm.db")
			Expect(gateway).To(BeNil())
			Expect(err).To(MatchError(`orm: unsupported driver: "sqlite4"`))
		})
	})

	Context("when the dns is not wrong", func() {
		It("returns an error", func() {
			gateway, err := orm.Open("mysql", "localhost")
			Expect(gateway).To(BeNil())
			Expect(err).To(MatchError(`invalid DSN: missing the slash separating the database name`))
		})
	})
})

var _ = Describe("Connect", func() {
	It("opens the URL successfully", func() {
		gateway, err := orm.Connect("sqlite3://orm.db")
		Expect(err).To(BeNil())
		Expect(gateway).NotTo(BeNil())
		Expect(gateway.Close()).To(Succeed())
	})

	Context("when cannot open the database", func() {
		It("returns an error", func() {
			gateway, err := orm.Connect(":::::")
			Expect(err).To(MatchError("parse \":::::\": missing protocol scheme"))
			Expect(gateway).To(BeNil())
		})
	})

	Context("when the dns is not wrong", func() {
		It("returns an error", func() {
			gateway, err := orm.Connect("sqlite3://localhost:5430/db")
			Expect(gateway).To(BeNil())
			Expect(err).To(MatchError("unable to open database file"))
		})
	})
})

var _ = Describe("Gateway", func() {
	var (
		ctx     context.Context
		gateway *orm.Gateway
	)

	BeforeEach(func() {
		var err error

		gateway, err = orm.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
		Expect(err).To(BeNil())

		ctx = context.TODO()

		query :=
			sql.CreateTable("users").IfNotExists().
				Columns(
					sql.Column("id").Type("int"),
					sql.Column("first_name").Type("varchar(255)").Attr("NOT NULL"),
					sql.Column("last_name").Type("varchar(255)").Attr("NOT NULL"),
					sql.Column("email").Type("varchar(255)").Attr("NULL"),
					sql.Column("created_at").Type("timestamp").Attr("NULL"),
				).
				PrimaryKey("id")

		_, err = gateway.Exec(ctx, query)
		Expect(err).To(Succeed())

		for i := 0; i < 10; i++ {
			query :=
				sql.Insert("users").
					Columns("id", "first_name", "last_name", "email").
					Values(i, faker.FirstName(), faker.LastName(), faker.Email()).
					Returning("id")

			result, err := gateway.Exec(ctx, query)
			Expect(err).To(Succeed())

			affected, err := result.RowsAffected()
			Expect(err).To(Succeed())
			Expect(affected).To(BeNumerically("==", 1))
		}
	})

	AfterEach(func() {
		_, err := gateway.Exec(ctx, sql.Raw("DELETE FROM users"))
		Expect(err).To(Succeed())

		Expect(gateway.Close()).To(Succeed())
	})

	Describe("ReadDir", func() {
		var fileSystem *fake.FileSystem

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

			fileSystem = &fake.FileSystem{}
			fileSystem.OpenFileReturns(parcello.NewResourceFile(node), nil)

			fileSystem.WalkStub = func(dir string, fn filepath.WalkFunc) error {
				return fn(node.Name, &parcello.ResourceFileInfo{Node: node}, nil)
			}
		})

		It("reads the directory", func() {
			Expect(gateway.ReadDir(fileSystem)).To(Succeed())
		})

		Context("when the routine is executed", func() {
			It("execs the routine", func() {
				Expect(gateway.ReadDir(fileSystem)).To(Succeed())
				_, err := gateway.Query(ctx, orm.Routine("cmd"))
				Expect(err).To(Succeed())
			})
		})
	})

	Describe("ReadFrom", func() {
		It("reads the routines from file", func() {
			script := fmt.Sprintf("%v", time.Now().UnixNano())
			buffer := bytes.NewBufferString(fmt.Sprintf("-- name: %v", script))
			fmt.Fprintln(buffer)
			fmt.Fprintln(buffer, "SELECT * FROM sqlite_master")
			_, err := gateway.ReadFrom(buffer)
			Expect(err).To(Succeed())
		})
	})

	Describe("Migrate", func() {
		var fileSystem *fake.FileSystem

		BeforeEach(func() {
			buffer := bytes.NewBufferString("-- name: up")

			fmt.Fprintln(buffer)
			fmt.Fprintln(buffer, "SELECT * FROM sqlite_master")

			content := buffer.Bytes()

			node := &parcello.Node{
				Name:    "00060524000000_setup.sql",
				Content: &content,
				Mutex:   &sync.RWMutex{},
			}

			fileSystem = &fake.FileSystem{}
			fileSystem.OpenFileReturns(parcello.NewResourceFile(node), nil)

			fileSystem.WalkStub = func(dir string, fn filepath.WalkFunc) error {
				return fn(node.Name, &parcello.ResourceFileInfo{Node: node}, nil)
			}

			query :=
				sql.CreateTable("migrations").IfNotExists().
					Columns(
						sql.Column("id").Type("varchar(15)"),
						sql.Column("description").Type("text").Attr("NOT NULL"),
						sql.Column("created_at").Type("timestamp").Attr("NOT NULL"),
					).
					PrimaryKey("id")

			_, err := gateway.Exec(ctx, query)
			Expect(err).To(Succeed())
		})

		It("executes the migration successfully", func() {
			Expect(gateway.Migrate(fileSystem)).To(Succeed())
		})
	})

	Describe("Begin", func() {
		It("starts new transaction", func() {
			tx, err := gateway.Begin(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(tx).NotTo(BeNil())
		})

		Context("when the gateway is closed", func() {
			It("returns an error", func() {
				gateway, err := orm.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
				Expect(err).To(BeNil())
				Expect(gateway.Close()).To(Succeed())

				tx, err := gateway.Begin(ctx)
				Expect(err).To(MatchError("sql: database is closed"))
				Expect(tx).To(BeNil())
			})
		})
	})

	Describe("RunInTx", func() {
		It("starts new transaction", func() {
			err := gateway.RunInTx(ctx, func(tx *orm.TxGateway) error {
				entities := []*User{}
				Expect(tx.All(ctx, sql.Raw("SELECT * FROM users"), &entities)).To(Succeed())

				entity := &User{}
				Expect(tx.Only(ctx, sql.Raw("SELECT * FROM users WHERE id = 0"), entity)).To(Succeed())
				Expect(tx.First(ctx, sql.Raw("SELECT * FROM users WHERE id = 0"), entity)).To(Succeed())

				_, err := tx.Exec(ctx, sql.Raw("UPDATE users SET created_at = strftime()"))
				Expect(err).To(Succeed())

				rows, err := tx.Query(ctx, sql.Raw("UPDATE users SET created_at = strftime()"))
				Expect(err).To(Succeed())
				Expect(rows.Close()).To(Succeed())

				return nil
			})

			Expect(err).NotTo(HaveOccurred())
		})

		Context("when the gateway is closed", func() {
			It("returns an error", func() {
				gateway, err := orm.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
				Expect(err).To(BeNil())
				Expect(gateway.Close()).To(Succeed())

				err = gateway.RunInTx(ctx, func(tx *orm.TxGateway) error {
					return nil
				})

				Expect(err).To(MatchError("sql: database is closed"))
			})
		})

		Context("when the transaction fails", func() {
			It("returns an error", func() {
				err := gateway.RunInTx(ctx, func(tx *orm.TxGateway) error {
					Expect(tx).NotTo(BeNil())
					return fmt.Errorf("oh no")
				})

				Expect(err).To(MatchError("oh no"))
			})
		})
	})

	Describe("Dialect", func() {
		It("returns the dialect", func() {
			Expect(gateway.Dialect()).To(Equal("sqlite3"))
		})
	})

	Describe("All", func() {
		It("returns all entities", func() {
			entities := []*User{}

			Expect(gateway.All(ctx, sql.Raw("SELECT * FROM users"), &entities)).To(Succeed())
			Expect(entities).To(HaveLen(10))

			Expect(entities[0].ID).To(Equal(0))
			Expect(entities[0].Email).NotTo(BeNil())

			Expect(entities[1].ID).To(Equal(1))
			Expect(entities[1].Email).NotTo(BeNil())
		})

		Context("when the database operation fail", func() {
			It("returns an error", func() {
				entities := []*User{}

				Expect(gateway.All(ctx, sql.Raw("SELECT * FROM unknown.users"), &entities)).To(MatchError("no such table: unknown.users"))
				Expect(entities).To(HaveLen(0))
			})
		})

		Context("when the routine is unknown", func() {
			It("returns an error", func() {
				entities := []*User{}

				err := gateway.All(ctx, sql.Routine("my-unknown-routine"), entities)
				Expect(err).To(MatchError("query 'my-unknown-routine' not found"))
			})
		})
	})

	Describe("Only", func() {
		It("returns the first entity", func() {
			entity := &User{}
			Expect(gateway.Only(ctx, sql.Raw("SELECT * FROM users WHERE id = 0"), entity)).To(Succeed())

			Expect(entity.ID).To(Equal(0))
			Expect(entity.Email).NotTo(BeNil())
		})

		Context("when the provided type is not compatible", func() {
			It("returns an error", func() {
				entity := "root"
				Expect(gateway.Only(ctx, sql.Raw("SELECT * FROM users WHERE id = 0"), &entity)).To(MatchError("sql/scan: columns do not match (5 > 1)"))
			})
		})

		Context("when there are more than one entities", func() {
			It("returns an error", func() {
				entity := &User{}
				Expect(gateway.Only(ctx, sql.Raw("SELECT * FROM users"), entity)).To(MatchError("orm: user not singular"))
			})
		})

		Context("when there are not entities", func() {
			It("returns an error", func() {
				entity := &User{}
				Expect(gateway.Only(ctx, sql.Raw("SELECT * FROM users WHERE id > 1000"), entity)).To(MatchError("orm: user not found"))
			})
		})

		Context("when the database operation fail", func() {
			It("returns an error", func() {
				entity := &User{}
				Expect(gateway.Only(ctx, sql.Raw("SELECT * FROM unknown.users"), entity)).To(MatchError("no such table: unknown.users"))
			})
		})

		Context("when the routine is unknown", func() {
			It("returns an error", func() {
				entity := &User{}
				err := gateway.Only(ctx, sql.Routine("my-unknown-routine"), entity)
				Expect(err).To(MatchError("query 'my-unknown-routine' not found"))
			})
		})
	})

	Describe("First", func() {
		It("returns the first entity", func() {
			entity := &User{}
			Expect(gateway.First(ctx, sql.Raw("SELECT * FROM users"), entity)).To(Succeed())

			Expect(entity.ID).To(Equal(0))
			Expect(entity.Email).NotTo(BeNil())
		})

		Context("when the provided type is not compatible", func() {
			It("returns an error", func() {
				entity := "root"
				Expect(gateway.First(ctx, sql.Raw("SELECT * FROM users WHERE id = 0"), &entity)).To(MatchError("sql/scan: columns do not match (5 > 1)"))
			})
		})

		Context("when there are not entities", func() {
			It("returns an error", func() {
				entity := &User{}
				Expect(gateway.First(ctx, sql.Raw("SELECT * FROM users WHERE id > 1000"), entity)).To(MatchError("orm: user not found"))
			})
		})

		Context("when the database operation fail", func() {
			It("returns an error", func() {
				entity := &User{}
				Expect(gateway.First(ctx, sql.Raw("SELECT * FROM unknown.users"), entity)).To(MatchError("no such table: unknown.users"))
			})
		})

		Context("when the routine is unknown", func() {
			It("returns an error", func() {
				entity := &User{}
				err := gateway.First(ctx, sql.Routine("my-unknown-routine"), entity)
				Expect(err).To(MatchError("query 'my-unknown-routine' not found"))
			})
		})
	})

	Describe("Exec", func() {
		Context("when the query has wrong syntax", func() {
			It("returns an error", func() {
				_, err := gateway.Exec(ctx, sql.Raw("SELECT * FROM unknown.users"))
				Expect(err).To(MatchError("no such table: unknown.users"))
			})
		})

		Context("when the routine is unknown", func() {
			It("returns an error", func() {
				_, err := gateway.Exec(ctx, sql.Routine("my-unknown-routine"))
				Expect(err).To(MatchError("query 'my-unknown-routine' not found"))
			})
		})
	})
})
