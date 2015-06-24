package shop_product_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	. "github.com/Shop2market/go-client/shop_product"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("ShopProduct", func() {
	It("deserializes a shop product", func() {
		products := []*Product{}
		body, _ := ioutil.ReadFile("fixtures/shop_product_index_response.json")
		json.Unmarshal(body, &products)

		Expect(products).To(Equal([]*Product{
			&Product{
				ProductName:    "Pastamachine",
				PictureLink:    "http://shopurl.com/data/product/groot-nieuw/KS100133z.png",
				Deeplink:       "http://shopurl.com/product/6367/4999-5011-5016/keuken/funcooking/overige-funcooking/pastamachine/",
				ShopCode:       "100133",
				ProductEan:     "8717253001334",
				Enabled:        false,
				ProductBrand:   "",
				DeliveryPeriod: "Voor 20.00 uur besteld, morgen in huis!",
				ProductInStock: 306,
				SellingPrice:   1999,
				ShopCategory:   "Keuken -> Funcooking -> Overige funcooking",
			},
			&Product{
				ProductName:    "BESTRON Stofzakken",
				PictureLink:    "http://shopurl.com/data/product/440x440/-100337-1.jpg",
				Deeplink:       "http://shopurl.com/product/100337/bestron-stofzakken-d0013s/",
				ShopCode:       "100337",
				ProductEan:     "8712184010486",
				Enabled:        true,
				ProductBrand:   "BESTRON",
				DeliveryPeriod: "Op werkdagen voor 21.00 besteld, morgen in huis",
				ProductInStock: 47,
				SellingPrice:   1250,
				ShopCategory:   "Huishouden -> Stofzuigen -> Stofzakken",
			},
		}))
	})
	It("Retrives the product", func() {
		content, err := ioutil.ReadFile("fixtures/shop_product_index_response.json")
		Expect(err).NotTo(HaveOccurred())

		server := ghttp.NewServer()
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyBasicAuth("username", "password"),
				ghttp.RespondWith(http.StatusOK, string(content)),
			),
		)
		Username = "username"
		Password = "password"
		Endpoint = server.URL()

		products, err := Find(1, []string{"001", "002"})
		Expect(err).NotTo(HaveOccurred())
		Expect(products).To(HaveLen(2))
		Expect(products).To(ContainElement(&Product{
			ProductName:    "Pastamachine",
			PictureLink:    "http://shopurl.com/data/product/groot-nieuw/KS100133z.png",
			Deeplink:       "http://shopurl.com/product/6367/4999-5011-5016/keuken/funcooking/overige-funcooking/pastamachine/",
			ShopCode:       "100133",
			ProductEan:     "8717253001334",
			Enabled:        false,
			ProductBrand:   "",
			DeliveryPeriod: "Voor 20.00 uur besteld, morgen in huis!",
			ProductInStock: 306,
			SellingPrice:   1999,
			ShopCategory:   "Keuken -> Funcooking -> Overige funcooking",
		}))

		Expect(products).To(ContainElement(&Product{
			ProductName:    "BESTRON Stofzakken",
			PictureLink:    "http://shopurl.com/data/product/440x440/-100337-1.jpg",
			Deeplink:       "http://shopurl.com/product/100337/bestron-stofzakken-d0013s/",
			ShopCode:       "100337",
			ProductEan:     "8712184010486",
			Enabled:        true,
			ProductBrand:   "BESTRON",
			DeliveryPeriod: "Op werkdagen voor 21.00 besteld, morgen in huis",
			ProductInStock: 47,
			SellingPrice:   1250,
			ShopCategory:   "Huishouden -> Stofzuigen -> Stofzakken",
		}))

	})
})
