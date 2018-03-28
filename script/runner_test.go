package script_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/svett/gom"
	"github.com/svett/gom/script"
)

var _ = Describe("Runner", func() {
	var runner *script.Runner

	BeforeEach(func() {
		dir, err := ioutil.TempDir("", "gom_runner")
		Expect(err).To(BeNil())

		db := filepath.Join(dir, "gom.db")
		gateway, err := gom.Open("sqlite3", db)
		Expect(err).To(BeNil())

		runner = &script.Runner{
			Dir:     dir,
			Gateway: gateway,
		}
	})

	JustBeforeEach(func() {
		command := &bytes.Buffer{}
		fmt.Fprintln(command, "-- name: system-tables")
		fmt.Fprintln(command, "SELECT * FROM sqlite_master")

		path := filepath.Join(runner.Dir, "commands.sql")
		Expect(ioutil.WriteFile(path, command.Bytes(), 0700)).To(Succeed())
	})

	AfterEach(func() {
		runner.Gateway.Close()
	})

	It("runs the command successfully", func() {
		rows, err := runner.Run("system-tables")
		Expect(err).To(Succeed())

		columns, err := rows.Columns()
		Expect(err).To(Succeed())
		Expect(columns).To(ContainElement("type"))
		Expect(columns).To(ContainElement("name"))
		Expect(columns).To(ContainElement("tbl_name"))
		Expect(columns).To(ContainElement("rootpage"))
		Expect(columns).To(ContainElement("sql"))
	})

	Context("when the command does not exist", func() {
		JustBeforeEach(func() {
			path := filepath.Join(runner.Dir, "commands.sql")
			Expect(os.Remove(path)).To(Succeed())
		})

		It("returns an error", func() {
			_, err := runner.Run("system-tables")
			Expect(err).To(MatchError("Command 'system-tables' not found"))
		})
	})

	Context("when the database is not available", func() {
		JustBeforeEach(func() {
			Expect(runner.Gateway.Close()).To(Succeed())
		})

		It("return an error", func() {
			_, err := runner.Run("system-tables")
			Expect(err).To(MatchError("sql: database is closed"))
		})
	})
})