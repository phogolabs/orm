package orm_test

import (
	"fmt"

	"github.com/phogolabs/orm"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotFoundError", func() {
	It("returns a not found error", func() {
		errx := &orm.NotFoundError{}
		Expect(errx.Error()).To(Equal("orm:  not found"))
	})

	Describe("IsNotFound", func() {
		It("returns true", func() {
			errx := &orm.NotFoundError{}
			Expect(orm.IsNotFound(errx)).To(BeTrue())
		})

		Context("when the error is nil", func() {
			It("returns false", func() {
				Expect(orm.IsNotFound(nil)).To(BeFalse())
			})
		})
	})

	Describe("MaskNotFound", func() {
		It("masks the error", func() {
			Expect(orm.MaskNotFound(&orm.NotFoundError{})).To(BeNil())
		})

		It("masks the error", func() {
			Expect(orm.MaskNotFound(fmt.Errorf("oh no"))).To(MatchError("oh no"))
		})
	})
})

var _ = Describe("NotSingularError", func() {
	It("returns a not singular error", func() {
		errx := &orm.NotSingularError{}
		Expect(errx.Error()).To(Equal("orm:  not singular"))
	})

	Describe("IsNotSingular", func() {
		It("returns true", func() {
			errx := &orm.NotSingularError{}
			Expect(orm.IsNotSingular(errx)).To(BeTrue())
		})

		Context("when the error is nil", func() {
			It("returns false", func() {
				Expect(orm.IsNotSingular(nil)).To(BeFalse())
			})
		})
	})
})

var _ = Describe("NotLoadedError", func() {
	It("returns a not loaded  error", func() {
		errx := &orm.NotLoadedError{}
		Expect(errx.Error()).To(Equal("orm:  edge was not loaded"))
	})

	Describe("IsNotLoaded", func() {
		It("returns true", func() {
			errx := &orm.NotLoadedError{}
			Expect(orm.IsNotLoaded(errx)).To(BeTrue())
		})

		Context("when the error is nil", func() {
			It("returns false", func() {
				Expect(orm.IsNotLoaded(nil)).To(BeFalse())
			})
		})
	})
})

var _ = Describe("ConstraintError", func() {
	It("returns a constraint error", func() {
		errx := &orm.ConstraintError{}
		Expect(errx.Error()).To(Equal("orm: constraint failed: "))
	})

	Describe("UnWrap", func() {
		It("returns the wrapped error", func() {
			errx := &orm.ConstraintError{}
			Expect(errx.Unwrap()).To(BeNil())
		})
	})

	Describe("IsConstraintError", func() {
		It("returns true", func() {
			errx := &orm.ConstraintError{}
			Expect(orm.IsConstraintError(errx)).To(BeTrue())
		})

		Context("when the error is nil", func() {
			It("returns false", func() {
				Expect(orm.IsConstraintError(nil)).To(BeFalse())
			})
		})
	})
})
