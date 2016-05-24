package channel_product_test

import (
	"io/ioutil"
	"net/http"

	. "github.com/Shop2market/go-client/channel_product"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Channel product", func() {
	Context("parameters construction", func() {
		It("supports time ranges", func() {
			server := ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/shops/1/publishers/5/ads", "start=20150105&end=20150505"),
					ghttp.RespondWith(http.StatusOK, "[]"),
				),
			)
			Endpoint = server.URL()

			FindAds(&AdQuery{ProductsQuery: &ProductsQuery{ShopId: 1, PublisherId: 5}, StartTimeId: "20150105", EndTimeId: "20150505"})
		})

		It("supports only_with_stats", func() {
			server := ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/shops/1/publishers/5/ads", "start=20150105&end=20150505&only_with_stats=true"),
					ghttp.RespondWith(http.StatusOK, "[]"),
				),
			)
			Endpoint = server.URL()

			FindAds(&AdQuery{ProductsQuery: &ProductsQuery{ShopId: 1, PublisherId: 5}, StartTimeId: "20150105", EndTimeId: "20150505", OnlyWithStats: true})
		})
	})

	It("deserializes response", func() {
		content, err := ioutil.ReadFile("fixtures/ads_response.json")
		Expect(err).NotTo(HaveOccurred())

		server := ghttp.NewServer()
		server.AppendHandlers(
			ghttp.RespondWith(http.StatusOK, string(content)),
		)

		Endpoint = server.URL()
		ads, _ := FindAds(&AdQuery{ProductsQuery: &ProductsQuery{ShopId: 1, PublisherId: 5}})
		Expect(ads).To((HaveLen(2)))
		timestr := "2015-11-04 14:16:58.52 +0000 UTC"
		Expect(ads[0].QuarantinedAt.String()).To(Equal(timestr))
	})
})
