package channel_product_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	. "github.com/Shop2market/go-client/channel_product"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Store", func() {

	It("deserialization of channel product", func() {
		products := []*Product{}
		body, _ := ioutil.ReadFile("fixtures/channel_product_response.json")
		json.Unmarshal(body, &products)

		Expect(products).To(Equal([]*Product{
			&Product{
				Id: &ProductId{
					ShopCode:    "151656",
					ShopId:      1,
					PublisherId: 5,
				},
				Active: true,
			},
			&Product{
				Id: &ProductId{
					ShopCode:    "149350",
					ShopId:      1,
					PublisherId: 5,
				},
				Active: true,
			},
		}))
	})

	Context("requests channel products", func() {
		It("sends correct parameters", func() {
			server := ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/shops/1/publishers/5/products", "enabled=true&limit=10&skip=0"),
					ghttp.RespondWith(http.StatusOK, "[]"),
				),
			)
			Endpoint = server.URL()

			Find(1, 5, 0, 10)
		})

		It("deserializes response", func() {
			content, err := ioutil.ReadFile("fixtures/channel_product_response.json")
			Expect(err).NotTo(HaveOccurred())

			server := ghttp.NewServer()
			server.AppendHandlers(
				ghttp.RespondWith(http.StatusOK, string(content)),
			)

			Endpoint = server.URL()

			Expect(Find(1, 5, 0, 10)).To(Equal([]*Product{
				&Product{
					Id: &ProductId{
						ShopCode:    "151656",
						ShopId:      1,
						PublisherId: 5,
					},
					Active: true,
				},
				&Product{
					Id: &ProductId{
						ShopCode:    "149350",
						ShopId:      1,
						PublisherId: 5,
					},
					Active: true,
				},
			}))

		})
	})

})
