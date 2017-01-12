package channel_product_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

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
		manuallySet := true

		Expect(products).To(Equal([]*Product{
			&Product{
				Id: &ProductId{
					ShopCode:    "151656",
					ShopId:      1,
					PublisherId: 5,
				},
				Active:             true,
				Enabled:            true,
				ManuallySet:        &manuallySet,
				ChannelCategoryIDs: []int{1, 2},
				Taxonomies:         map[string][]int{"1": []int{2, 3}},
			},
			&Product{
				Id: &ProductId{
					ShopCode:    "149350",
					ShopId:      1,
					PublisherId: 5,
				},
				Active:             true,
				Enabled:            true,
				ChannelCategoryIDs: []int{1, 2},
			},
		}))
	})
	Context("posts touch on channel products", func() {
		It("PUT's to touch", func() {
			shopCodes := []string{"a", "b"}
			server := ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/shops/1/publishers/5/products/touch", ""),
					ghttp.VerifyJSONRepresenting(shopCodes),
				),
			)
			Endpoint = server.URL()

			Touch(1, 5, shopCodes)
		})
	})
	Context("posts webhook on channel products", func() {
		It("PUT's to webhook", func() {
			shopCodes := []string{"a", "b"}
			server := ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/shops/1/products/webhook", ""),
					ghttp.VerifyJSONRepresenting(shopCodes),
				),
			)
			Endpoint = server.URL()

			Webhook(1, shopCodes)
		})
		It("doesn't PUT to webhook", func() {
			shopCodes := []string{}
			server := ghttp.NewServer()
			Endpoint = server.URL()
			Webhook(1, shopCodes)
		})
	})
	Context("requests channel products", func() {
		Context("parameters construction", func() {
			It("supports active", func() {
				server := ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/shops/1/publishers/5/products", "active=true"),
						ghttp.RespondWith(http.StatusOK, "[]"),
					),
				)
				Endpoint = server.URL()

				active := true
				Find(&ProductsQuery{ShopId: 1, PublisherId: 5, Active: &active})
			})
			It("supports active, false", func() {
				server := ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/shops/1/publishers/5/products", "active=false"),
						ghttp.RespondWith(http.StatusOK, "[]"),
					),
				)
				Endpoint = server.URL()

				active := false
				Find(&ProductsQuery{ShopId: 1, PublisherId: 5, Active: &active})
			})

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
						ghttp.VerifyRequest("GET", "/shops/1/publishers/5/products", "manually_set=true"),
						ghttp.RespondWith(http.StatusOK, "[]"),
					),
				)
				Endpoint = server.URL()

				manuallyDeactivated := true
				Find(&ProductsQuery{ShopId: 1, PublisherId: 5, ManuallySet: &manuallyDeactivated})
			})
			It("supports manually_set_tip false", func() {
				server := ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/shops/1/publishers/5/products", "manually_set_tip=false"),
						ghttp.RespondWith(http.StatusOK, "[]"),
					),
				)
				Endpoint = server.URL()

				manuallySetTip := false
				Find(&ProductsQuery{ShopId: 1, PublisherId: 5, ManuallySetTip: &manuallySetTip})
			})
			It("supports manually_set_tip true", func() {
				server := ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/shops/1/publishers/5/products", "manually_set_tip=true"),
						ghttp.RespondWith(http.StatusOK, "[]"),
					),
				)
				Endpoint = server.URL()

				manuallySetTip := true
				Find(&ProductsQuery{ShopId: 1, PublisherId: 5, ManuallySetTip: &manuallySetTip})
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
						ghttp.VerifyRequest("GET", "/shops/1/publishers/5/products", "manually_set=false"),
						ghttp.RespondWith(http.StatusOK, "[]"),
					),
				)
				Endpoint = server.URL()

				manuallyDeactivated := false
				Find(&ProductsQuery{ShopId: 1, PublisherId: 5, ManuallySet: &manuallyDeactivated})
			})
			It("supports last_updated_before, false", func() {
				server := ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/shops/1/publishers/5/products", "last_updated_before=2014-09-12"),
						ghttp.RespondWith(http.StatusOK, "[]"),
					),
				)
				Endpoint = server.URL()

				lastUpdatedBefore := time.Date(2014, 9, 12, 0, 0, 0, 0, time.UTC)

				Find(&ProductsQuery{ShopId: 1, PublisherId: 5, LastUpdatedBefore: &lastUpdatedBefore})
			})
		})
		It("sends correct parameters", func() {
			server := ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/shops/1/publishers/5/products", "active=true&limit=10&skip=0"),
					ghttp.RespondWith(http.StatusOK, "[]"),
				),
			)
			Endpoint = server.URL()

			active := true
			limit := 10
			skip := 0
			Find(&ProductsQuery{ShopId: 1, PublisherId: 5, Active: &active, Limit: &limit, Skip: &skip})
		})
	})
})
