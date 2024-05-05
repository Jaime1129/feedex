package components

import (
	"testing"

	"github.com/jarcoal/httpmock"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func Test(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Running Test Suite")
}

var _ = ginkgo.Describe("EthScanCli", func() {
	var (
		client EthScanCli
		apiKey string
	)

	ginkgo.BeforeEach(func() {
		apiKey = "anykey"
		client = NewEthScanCli(apiKey)
		httpmock.Activate()
	})

	ginkgo.AfterEach(func() {
		httpmock.DeactivateAndReset()
	})

	ginkgo.Describe("QueryTrxFee", func() {
		ginkgo.It("should return transaction fee information correctly", func() {
			httpmock.RegisterResponder("GET", `=~https://api.etherscan.io/api`,
				httpmock.NewStringResponder(200, `{"result":{"effectiveGasPrice":"100","gasUsed":"21000"},"error":{"code":0}}`))

			trxHash := "some-trx-hash"
			trxResp, err := client.QueryTrxFee(trxHash)
			gomega.Expect(err).Should(gomega.BeNil())
			gomega.Expect(trxResp.Result.EffectiveGasPrice).To(gomega.Equal("100"))
			gomega.Expect(trxResp.Result.GasUsed).To(gomega.Equal("21000"))
		})

		ginkgo.It("should handle API error responses", func() {
			httpmock.RegisterResponder("GET", `=~https://api.etherscan.io/api`,
				httpmock.NewStringResponder(200, `{"error":{"code":1,"message":"Error in API"}}`))

			trxHash := "some-trx-hash"
			_, err := client.QueryTrxFee(trxHash)
			gomega.Expect(err).ShouldNot(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.Equal("ethscan api call returns error"))
		})
	})

	ginkgo.Describe("QueryBlock", func() {
		ginkgo.It("should return block information correctly", func() {
			httpmock.RegisterResponder("GET", `=~https://api.etherscan.io/api`,
				httpmock.NewStringResponder(200, `{"result":{"timestamp":"1609459200"},"error":{"code":0}}`))

			blockNumber := "123456"
			blockResp, err := client.QueryBlock(blockNumber)
			gomega.Expect(err).Should(gomega.BeNil())
			gomega.Expect(blockResp.Result.Timestamp).To(gomega.Equal("1609459200"))
		})

		ginkgo.It("should handle API error responses", func() {
			httpmock.RegisterResponder("GET", `=~https://api.etherscan.io/api`,
				httpmock.NewStringResponder(200, `{"error":{"code":1,"message":"Error in API"}}`))

			blockNumber := "123456"
			_, err := client.QueryBlock(blockNumber)
			gomega.Expect(err).ShouldNot(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.Equal("ethscan api call returns error"))
		})
	})
})
