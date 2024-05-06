package components

import (
	"github.com/jarcoal/httpmock"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("BnPriceCli", func() {
	var (
		client BnPriceCli
	)

	ginkgo.BeforeEach(func() {
		client = NewBnPriceCLi()
		httpmock.Activate()
	})

	ginkgo.AfterEach(func() {
		httpmock.DeactivateAndReset()
	})

	ginkgo.It("should return the average price correctly", func() {
		httpmock.RegisterResponder("GET", "https://api.binance.com/api/v3/klines",
			httpmock.NewStringResponder(200, `[["", "100.5", "", "", "101.5"]]`))

		price, err := client.QueryETHPrice(1609459200, 1609545600, INTERVAL_1MIN)
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(price.String()).To(gomega.Equal("101"))
	})

	ginkgo.It("should handle an error response from the API", func() {
		httpmock.RegisterResponder("GET", "https://api.binance.com/api/v3/klines",
			httpmock.NewStringResponder(500, ""))

		_, err := client.QueryETHPrice(1609459200, 1609545600, INTERVAL_1MIN)
		gomega.Expect(err).ShouldNot(gomega.BeNil())
	})

	ginkgo.It("should handle no data being returned", func() {
		httpmock.RegisterResponder("GET", "https://api.binance.com/api/v3/klines",
			httpmock.NewStringResponder(200, `[]`))

		_, err := client.QueryETHPrice(1609459200, 1609545600, INTERVAL_1MIN)
		gomega.Expect(err).ShouldNot(gomega.BeNil())
		gomega.Expect(err.Error()).To(gomega.Equal("price not found"))
	})
})
