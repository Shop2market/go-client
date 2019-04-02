package channel_product_test

import (
	"io/ioutil"
	"net/http"

	. "github.com/Shop2market/go-client/channel_product"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Ads", func() {
	Context("parameters construction", func() {
		It("supports time ranges", func() {
			server := ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/shops/1/publishers/5/ads", "start=20150105&stop=20150505"),
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
					ghttp.VerifyRequest("GET", "/shops/1/publishers/5/ads", "start=20150105&stop=20150505&only_with_stats=true"),
					ghttp.RespondWith(http.StatusOK, "[]"),
				),
			)
			Endpoint = server.URL()

			FindAds(&AdQuery{ProductsQuery: &ProductsQuery{ShopId: 1, PublisherId: 5}, StartTimeId: "20150105", EndTimeId: "20150505", OnlyWithStats: true})
		})
	})

	It("Ads deserializes response", func() {
		content, err := ioutil.ReadFile("fixtures/ads_response.json")
		Expect(err).NotTo(HaveOccurred())

		server := ghttp.NewServer()
		server.AppendHandlers(
			ghttp.RespondWith(http.StatusOK, string(content)),
		)

		Endpoint = server.URL()
		ads, _ := FindAds(&AdQuery{ProductsQuery: &ProductsQuery{ShopId: 1, PublisherId: 5}})
		Expect(ads).To((HaveLen(4)))
		timestr := "2015-11-04 14:16:58.52 +0000 UTC"
		Expect(ads[0].QuarantinedAt.String()).To(Equal(timestr))
		Expect(ads[0].StockStatus).To(Equal("yes"))
		Expect(ads[0].ProductsInStock).To(Equal(0))
		Expect(ads[0].DeactivationReason).To(Equal("content_rules"))
		Expect(ads[0].IsMappedToTaxonomy("297", false)).To(BeFalse())

		Expect(ads[2].OrderAmountExcludingTax).To(Equal(33.33))
		Expect(ads[2].ChannelCategoryIDs).To(Equal([]int{414341}))
		Expect(ads[2].Taxonomies).To(Equal(map[string][]int{"297": []int{401856}}))
		Expect(ads[2].RulesTaxonomies).To(Equal(map[string][]int{"326": []int{414381}}))
		Expect(ads[2].IsMappedToTaxonomy("297", false)).To(BeTrue())
		Expect(ads[2].IsMappedToTaxonomy("326", true)).To(BeTrue())
		Expect(ads[2].AllTaxonomies("326")).To(Equal([]int{414341, 401856}))

		Expect(ads[3].RulesTaxonomies).To(Equal(map[string][]int{"326": []int{414381}}))
		Expect(ads[3].AllTaxonomies("326")).To(Equal([]int{414381}))
		Expect(ads[3].IsMappedToTaxonomy("297", false)).To(BeFalse())
		Expect(ads[3].IsMappedToTaxonomy("326", true)).To(BeTrue())

	})
})
