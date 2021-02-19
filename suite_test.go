package orm_test

import (
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestORM(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ORM Suite")
}

// User represents a user
type User struct {
	ID        int        `db:"id,primary_key,not_null,read_only"`
	FirstName string     `db:"first_name,not_null"`
	LastName  string     `db:"last_name,not_null"`
	Email     *string    `db:"email"`
	CreatedAt *time.Time `db:"created_at"`
}
