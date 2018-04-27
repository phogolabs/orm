package oak_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/oak"
	"github.com/phogolabs/oak/fake"
	"github.com/phogolabs/parcello"
)

var _ = Describe("Command", func() {
	var script string

	BeforeEach(func() {
		script = fmt.Sprintf("%v", time.Now().UnixNano())
		buffer := bytes.NewBufferString(fmt.Sprintf("-- name: %v", script))
		fmt.Fprintln(buffer)
		fmt.Fprintln(buffer, "SELECT * FROM users")
		Expect(oak.LoadSQLCommandsFromReader(buffer)).To(Succeed())
	})

	It("returns a command", func() {
		stmt := oak.Command(script)
		Expect(stmt).NotTo(BeNil())

		query, params := stmt.Prepare()
		Expect(query).To(Equal("SELECT * FROM users"))
		Expect(params).To(BeEmpty())
	})

	It("returns a named command", func() {
		stmt := oak.NamedCommand(script, oak.P{})
		Expect(stmt).NotTo(BeNil())

		query, params := stmt.Prepare()
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
				return fn("command", &parcello.ResourceFileInfo{Node: node}, nil)
			}

			Expect(oak.LoadSQLCommandsFrom(fileSystem)).To(Succeed())
		})

		It("returns a command", func() {
			stmt := oak.Command("cmd")
			Expect(stmt).NotTo(BeNil())

			query, params := stmt.Prepare()
			Expect(query).To(Equal("SELECT * FROM categories"))
			Expect(params).To(BeEmpty())
		})
	})

	Context("when the statement does not exits", func() {
		It("does not return a statement", func() {
			Expect(func() { oak.Command("down") }).To(Panic())
		})
	})

	Context("when the named statement does not exits", func() {
		It("does not return a statement", func() {
			Expect(func() { oak.NamedCommand("down", oak.P{}) }).To(Panic())
		})
	})
})

var _ = Describe("ParseURL", func() {
	It("parses the SQLite connection string successfully", func() {
		driver, source, err := oak.ParseURL("sqlite3://./oak.db")
		Expect(err).NotTo(HaveOccurred())
		Expect(driver).To(Equal("sqlite3"))
		Expect(source).To(Equal("./oak.db"))
	})

	It("parses the MySQL connection string successfully", func() {
		driver, source, err := oak.ParseURL("mysql://root@/oak")
		Expect(err).NotTo(HaveOccurred())
		Expect(driver).To(Equal("mysql"))
		Expect(source).To(Equal("root@/oak"))
	})

	It("parses the PostgreSQL connection string successfully", func() {
		driver, source, err := oak.ParseURL("postgres://localhost/oak?sslmode=disable")
		Expect(err).NotTo(HaveOccurred())
		Expect(driver).To(Equal("postgres"))
		Expect(source).To(Equal("postgres://localhost/oak?sslmode=disable"))
	})

	Context("when the URL is invalid", func() {
		It("returns an error", func() {
			driver, source, err := oak.ParseURL("::")
			Expect(driver).To(BeEmpty())
			Expect(source).To(BeEmpty())
			Expect(err).To(MatchError("parse ::: missing protocol scheme"))
		})
	})
})

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

var _ = Describe("Setup", func() {
	var (
		manager *fake.FileSystemManager
		gateway *oak.Gateway
	)

	BeforeEach(func() {
		dir, err := ioutil.TempDir("", "oak_generator")
		Expect(err).To(BeNil())
		url := filepath.Join(dir, "oak.db")

		manager = &fake.FileSystemManager{}
		manager.RootReturns(manager, nil)

		manager.OpenFileStub = func(name string, flag int, perm os.FileMode) (parcello.File, error) {
			fileContent := &bytes.Buffer{}
			file := &fake.File{}

			file.WriteStub = func(data []byte) (int, error) {
				return fileContent.Write(data)
			}

			file.ReadStub = func(data []byte) (int, error) {
				content := bytes.NewBuffer(fileContent.Bytes())
				return content.Read(data)
			}

			return file, nil
		}

		gateway, err = oak.Open("sqlite3", url)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		gateway.Close()
	})

	It("setups the project successfully", func() {
		Expect(oak.Setup(gateway, parcello.Dir("./example/database"))).To(Succeed())
	})

	Context("when the resource script is not found", func() {
		BeforeEach(func() {
			manager.RootReturns(nil, fmt.Errorf("Oh no!"))
		})

		It("returns an error", func() {
			Expect(oak.Setup(gateway, manager)).To(MatchError("Oh no!"))
		})
	})

	Context("when the resource script cannot be load", func() {
		BeforeEach(func() {
			manager.WalkReturns(fmt.Errorf("Oh no!"))
		})

		It("returns an error", func() {
			Expect(oak.Setup(gateway, manager)).To(MatchError("Oh no!"))
		})
	})

	Context("when the resource migration is not found", func() {
		BeforeEach(func() {
			manager.RootStub = func(name string) (parcello.FileSystemManager, error) {
				if name == "script" {
					return manager, nil
				}
				return nil, fmt.Errorf("Oh no!")
			}
		})

		It("returns an error", func() {
			Expect(oak.Setup(gateway, manager)).To(MatchError("Oh no!"))
		})
	})

	Context("when the resource script cannot be load", func() {
		It("returns an error", func() {
			Expect(oak.Setup(gateway, manager)).To(HaveOccurred())
		})
	})

	Context("when the loading the migration fails", func() {
		It("returns an error", func() {
			Expect(oak.Setup(gateway, manager)).To(MatchError("Command 'up' not found for migration '00060524000000_setup.sql'"))
		})
	})
})
