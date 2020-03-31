package orm_test

import (
	"github.com/phogolabs/orm"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cursor", func() {
	Describe("DecodeCursor", func() {
		It("decodes the cursor successfully", func() {
			cursor, err := orm.DecodeCursor("W3siY29sdW1uIjoibmFtZSIsIm9yZGVyIjoiKyIsInZhbHVlIjpudWxsfSx7ImNvbHVtbiI6ImlkIiwib3JkZXIiOiIrIiwidmFsdWUiOm51bGx9XQ")
			Expect(err).To(BeNil())
			Expect(cursor).NotTo(BeNil())

			positions := *cursor
			Expect(positions).To(HaveLen(2))

			Expect(positions[0].Column).To(Equal("name"))
			Expect(positions[0].Order).To(Equal("+"))

			Expect(positions[1].Column).To(Equal("id"))
			Expect(positions[1].Order).To(Equal("+"))
		})

		Context("when the token is empty", func() {
			It("decodes the cursor successfully", func() {
				cursor, err := orm.DecodeCursor("")
				Expect(err).To(BeNil())
				Expect(cursor).NotTo(BeNil())

				positions := *cursor
				Expect(positions).To(BeEmpty())
			})
		})

		Context("when the token is invalid", func() {
			It("returns an error", func() {
				cursor, err := orm.DecodeCursor("wrong")
				Expect(err).To(MatchError("illegal base64 data at input byte 5"))
				Expect(cursor).To(BeNil())
			})
		})
	})

	Describe("String", func() {
		It("returns the token", func() {
			cursor := &orm.Cursor{
				&orm.Position{Column: "name", Order: "+"},
				&orm.Position{Column: "id", Order: "+"},
			}

			Expect(cursor.String()).To(Equal("W3siY29sdW1uIjoibmFtZSIsIm9yZGVyIjoiKyIsInZhbHVlIjpudWxsfSx7ImNvbHVtbiI6ImlkIiwib3JkZXIiOiIrIiwidmFsdWUiOm51bGx9XQ"))
		})

		Context("when the cursor is empty", func() {
			It("returns the token", func() {
				cursor := &orm.Cursor{}
				Expect(cursor.String()).To(BeEmpty())
			})
		})
	})
})
