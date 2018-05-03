package oak_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/oak"
)

var _ = Describe("Middleware", func() {
	var (
		r *http.Request
		g *oak.Gateway
	)

	BeforeEach(func() {
		g = &oak.Gateway{}
		r = httptest.NewRequest("GET", "http://example.com", nil)
	})

	It("sets the gateway successfully", func() {
		r = oak.WithGateway(r, g)
		Expect(oak.GetGateway(r)).To(Equal(g))
	})

	It("sets the middleware successfully", func() {
		wr := httptest.NewRecorder()
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
			Expect(oak.GetGateway(r)).To(Equal(g))
		})

		router := oak.GatewayHandler(g)(h)
		router.ServeHTTP(wr, r)
		Expect(wr.Code).To(Equal(http.StatusNoContent))
	})

	It("formats the key correctly", func() {
		Expect(oak.GatewayCtxKey.String()).To(Equal("oak/middleware context value Gateway"))
	})
})
