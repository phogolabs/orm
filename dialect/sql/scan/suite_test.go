package scan_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestScan(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scan Suite")
}
