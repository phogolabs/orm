package orm_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/orm"
)

var _ = Describe("Middleware", func() {
	var (
		r *http.Request
		g *orm.Gateway
	)

	BeforeEach(func() {
		g = &orm.Gateway{}
		r = httptest.NewRequest("GET", "http://example.com", nil)
	})

	It("sets the gateway successfully", func() {
		ctx := orm.SetContext(r.Context(), g)
		r = r.WithContext(ctx)

		gw := orm.GetContext(r.Context())
		Expect(gw).To(Equal(g))
	})

	Context("when the gateway is not presented", func() {
		It("returns an error", func() {
			gw := orm.GetContext(r.Context())
			Expect(gw).To(BeNil())
		})
	})

	It("formats the key correctly", func() {
		Expect(orm.GatewayCtxKey.String()).To(Equal("orm/middleware context value Gateway"))
	})
})
