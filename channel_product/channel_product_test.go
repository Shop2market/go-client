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
				Active:              true,
				ManuallyDeactivated: true,
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
		Context("parameters construction", func() {
			It("supports enabled", func() {
				server := ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/shops/1/publishers/5/products", "enabled=true"),
						ghttp.RespondWith(http.StatusOK, "[]"),
					),
				)
				Endpoint = server.URL()

				enabled := true
				Find(&ProductsQuery{ShopId: 1, PublisherId: 5, Enabled: &enabled})
			})
			It("supports enabled, false", func() {
				server := ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/shops/1/publishers/5/products", "enabled=false"),
						ghttp.RespondWith(http.StatusOK, "[]"),
					),
				)
				Endpoint = server.URL()

				enabled := false
				Find(&ProductsQuery{ShopId: 1, PublisherId: 5, Enabled: &enabled})
			})
			It("supports limit", func() {
				server := ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/shops/1/publishers/5/products", "limit=1"),
						ghttp.RespondWith(http.StatusOK, "[]"),
					),
				)
				Endpoint = server.URL()

				limit := 1
				Find(&ProductsQuery{ShopId: 1, PublisherId: 5, Limit: &limit})
			})
			It("supports skip", func() {
				server := ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/shops/1/publishers/5/products", "skip=1"),
						ghttp.RespondWith(http.StatusOK, "[]"),
					),
				)
				Endpoint = server.URL()

				skip := 1
				Find(&ProductsQuery{ShopId: 1, PublisherId: 5, Skip: &skip})
			})
			It("supports manually_deactivated", func() {
				server := ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/shops/1/publishers/5/products", "manually_deactivated=true"),
						ghttp.RespondWith(http.StatusOK, "[]"),
					),
				)
				Endpoint = server.URL()

				manuallyDeactivated := true
				Find(&ProductsQuery{ShopId: 1, PublisherId: 5, ManuallyDeactivated: &manuallyDeactivated})
			})
			It("supports shop_codes", func() {
				server := ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/shops/1/publishers/5/products", "shop_codes%5B%5D=ED1&shop_codes%5B%5D=ED2"),
						ghttp.RespondWith(http.StatusOK, "[]"),
					),
				)
				Endpoint = server.URL()

				shopCodes := []string{"ED1", "ED2"}
				Find(&ProductsQuery{ShopId: 1, PublisherId: 5, ShopCodes: &shopCodes})
			})
			It("supports manually_deactivated, false", func() {
				server := ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/shops/1/publishers/5/products", "manually_deactivated=false"),
						ghttp.RespondWith(http.StatusOK, "[]"),
					),
				)
				Endpoint = server.URL()

				manuallyDeactivated := false
				Find(&ProductsQuery{ShopId: 1, PublisherId: 5, ManuallyDeactivated: &manuallyDeactivated})
			})
		})
		It("sends correct parameters", func() {
			server := ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/shops/1/publishers/5/products", "enabled=true&limit=10&skip=0"),
					ghttp.RespondWith(http.StatusOK, "[]"),
				),
			)
			Endpoint = server.URL()

			enabled := true
			limit := 10
			skip := 0
			Find(&ProductsQuery{ShopId: 1, PublisherId: 5, Enabled: &enabled, Limit: &limit, Skip: &skip})
		})

		It("deserializes response", func() {
			content, err := ioutil.ReadFile("fixtures/channel_product_response.json")
			Expect(err).NotTo(HaveOccurred())

			server := ghttp.NewServer()
			server.AppendHandlers(
				ghttp.RespondWith(http.StatusOK, string(content)),
			)

			Endpoint = server.URL()

			Expect(Find(&ProductsQuery{ShopId: 1, PublisherId: 5})).To(Equal([]*Product{
				&Product{
					Id: &ProductId{
						ShopCode:    "151656",
						ShopId:      1,
						PublisherId: 5,
					},
					Active:              true,
					ManuallyDeactivated: true,
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
