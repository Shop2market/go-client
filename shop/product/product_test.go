package product_test

import (
	"io/ioutil"
	"net/http"

	"github.com/Shop2market/go-client/catalog"
	. "github.com/Shop2market/go-client/shop/product"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Product", func() {
	Context("FindAll", func() {
		It("returns all products if there is data", func() {
			content, err := ioutil.ReadFile("fixtures/bonobo_products.jsonl")
			Expect(err).NotTo(HaveOccurred())

			server := ghttp.NewServer()
			Endpoint = server.URL()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/shops/1792/mongo/products", "sorted=true&recent=true"),
					ghttp.RespondWith(http.StatusOK, string(content)),
				),
			)

			productChan, errorChan := FindAll(1792)
			var product catalog.ShopProduct
			product = <-productChan
			Expect(product.ShopCode()).To(Equal("13861613"))
			Expect(product.ProductName()).To(Equal("adidaS Black leather backpack"))
			Expect(product.VariantID()).To(Equal("26621028"))
			Expect(product.Description()).To(Equal("dolor.\r\n\r\nNieuwe regel test."))
			Expect(product.FloorPrice()).To(Equal("998"))
			Expect(product.ProductBrand()).To(Equal("Adidas"))
			Expect(product.Deeplink()).To(Equal("http://freek-en-dreesman.webshopapp.com/adidas-adidas-black-leather-backpack.html"))
			Expect(product.PictureLink()).To(Equal("http://static.webshopapp.com/shops/073661/files/028410639/file.jpg"))
			Expect(product.ProductEan()).To(Equal("234565"))
			Expect(product.CategoryPath()).To(Equal("Women»T-Shirts"))
			Expect(product.ProductInStock()).To(Equal("34"))
			Expect(product.StockStatus()).To(Equal("true"))
			Expect(product.SubCategory()).To(Equal(""))
			Expect(product.DeliveryPeriod()).To(Equal("1 week"))
			Expect(product.SellingPrice()).To(Equal("1000"))

			product = <-productChan
			Expect(product.ShopCode()).To(Equal("13861613"))
			Expect(product.VariantID()).To(Equal("26621030"))
			Expect(product.FloorPrice()).To(Equal(""))
			Expect(product.ProductEan()).To(Equal(""))
			Expect(product.ProductInStock()).To(Equal("0"))

			product = <-productChan
			Expect(product.ShopCode()).To(Equal("13861613"))
			Expect(product.VariantID()).To(Equal("26621032"))
			Expect(product.FloorPrice()).To(Equal("0"))

			product = <-productChan
			Expect(product.ShopCode()).To(Equal("13861613"))
			Expect(product.VariantID()).To(Equal("26621034"))

			product = <-productChan
			Expect(product.ShopCode()).To(Equal("13861613"))
			Expect(product.VariantID()).To(Equal("26621036"))

			product = <-productChan
			Expect(product.ShopCode()).To(Equal("13861617"))
			Expect(product.VariantID()).To(Equal("24435201"))

			product = <-productChan
			Expect(product.ShopCode()).To(Equal("13861617"))
			Expect(product.VariantID()).To(Equal("64505915"))

			product = <-productChan
			Expect(product.ShopCode()).To(Equal("13861617"))
			Expect(product.VariantID()).To(Equal("64500617"))
			Expect(<-errorChan).To(BeNil())
			Expect(errorChan).Should(BeClosed())
			Expect(productChan).Should(BeClosed())
		})
		It("returns nothing if there is no data", func() {
			server := ghttp.NewServer()
			Endpoint = server.URL()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/shops/1792/mongo/products", "sorted=true&recent=true"),
					ghttp.RespondWith(http.StatusOK, ""),
				),
			)
			productChan, errorChan := FindAll(1792)
			Expect(<-errorChan).To(BeNil())
			Expect(errorChan).Should(BeClosed())
			Expect(productChan).Should(BeClosed())
		})
		It("returns error if there is bad json", func() {
			server := ghttp.NewServer()
			Endpoint = server.URL()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/shops/1792/mongo/products", "sorted=true&recent=true"),
					ghttp.RespondWith(http.StatusOK, "{["),
				),
			)
			productChan, errorChan := FindAll(1792)
			Expect(<-errorChan).To(HaveOccurred())
			Expect(errorChan).Should(BeClosed())
			Expect(productChan).Should(BeClosed())
		})
	})
})
