package scan_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestScan(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scan Suite")
}

type User struct {
	ID        string     `db:"id,primary_key"`
	Name      string     `db:"name"`
	Email     *string    `db:"email"`
	DeletedAt *time.Time `db:"deleted_at"`
	CreatedAt time.Time  `db:"created_at,read_only"`
}

func StringToPtr(v string) *string {
	return &v
}
