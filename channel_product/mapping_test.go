package channel_product_test

import (
	. "github.com/Shop2market/go-client/channel_product"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Mapping", func() {
	Context("posts map form tip on channel products", func() {
		It("PUT's to map_from_tip", func() {
			mappingParams := MappingParams{ShopCodes: []string{"a", "b"}, ChannelCategoryIds: []int{100}, UserID: 1}
			server := ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/shops/1/publishers/5/products/map_from_tip", ""),
					ghttp.VerifyJSONRepresenting(mappingParams),
				),
			)
			Endpoint = server.URL()
			Expect(Map(1, 5, []string{"a", "b"}, 1, 100)).NotTo(HaveOccurred())
		})
	})
})
