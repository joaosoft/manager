package mgr

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
)

var _ = Describe("Dummy Test", func() {
	Describe("when calling doing nothing", func() {
		BeforeEach(func() {

		})

		Context("when doing nothing but get an error", func() {
			It("should return an error", func() {
				_, err := decimal.NewFromString("xxx")

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
