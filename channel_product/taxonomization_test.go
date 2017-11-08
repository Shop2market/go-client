package channel_product_test

import (
	. "github.com/Shop2market/go-client/channel_product"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Taxonomization", func() {
	Context("puts map form tip on channel products", func() {
		It("PUT's to map_from_tip", func() {
			taxonomizationParams := TaxonomyParams{ShopCodes: []string{"a", "b"}, Taxonomies: TaxonomyTypeParams{TypeID: "1", ID: 100}, UserID: 1}
			server := ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/shops/1/publishers/5/products/taxonomy_from_tip", ""),
					ghttp.VerifyJSONRepresenting(taxonomizationParams),
				),
			)
			Endpoint = server.URL()
			Expect(Taxonomize(1, 5, []string{"a", "b"}, 1, 100)).NotTo(HaveOccurred())
		})
	})
})
