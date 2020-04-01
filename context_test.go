package orm_test

import (
	"context"

	"github.com/phogolabs/orm"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Context", func() {
	Describe("NewGatewayContext", func() {
		It("creates a new context", func() {
			gateway := &orm.Gateway{}
			ctx := orm.NewGatewayContext(context.TODO(), gateway)
			Expect(ctx.Value(orm.GatewayCtxKey)).To(Equal(gateway))
		})
	})

	Describe("GetewayContextFrom", func() {
		It("returns the gateway", func() {
			gateway := &orm.Gateway{}
			ctx := orm.NewGatewayContext(context.TODO(), gateway)

			injected := orm.GatewayFromContext(ctx)
			Expect(injected).To(Equal(gateway))
			Expect(injected).NotTo(BeNil())
		})

		Context("when the gateway is not set", func() {
			It("returns nil gateway", func() {
				gateway := orm.GatewayFromContext(context.TODO())
				Expect(gateway).To(BeNil())
			})
		})
	})

	It("formats the key correctly", func() {
		Expect(orm.GatewayCtxKey.String()).To(Equal("orm: context value gateway"))
	})
})
