package catalog_test

import (
	"io/ioutil"
	"net/http"

	. "github.com/Shop2market/go-client/catalog"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("ShopProduct", func() {
	It("should read products the catalog, one at a time", func() {
		content, err := ioutil.ReadFile("fixtures/1.jsonl")
		Expect(err).NotTo(HaveOccurred())

		server := ghttp.NewServer()
		Endpoint = server.URL()
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, string(content)),
			),
		)
		productChan, errorChan := Find(1)
		product := <-productChan
		Expect(product["Shop Code"]).To(Equal("20201"))
		Expect(product.SubCategory()).To(Equal(""))
		product = <-productChan
		Expect(product["Shop Code"]).To(Equal("20301"))
		Expect(product.SubCategory()).To(Equal("Sub category name"))
		product = <-productChan
		Expect(errorChan).Should(BeClosed())
		Expect(productChan).Should(BeClosed())
	})

	It("should return error if json error", func() {
		server := ghttp.NewServer()
		Endpoint = server.URL()
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, string("non json")),
			),
		)
		productChan, errorChan := Find(1)

		err := <-errorChan
		Expect(err).To(HaveOccurred())

		Expect(productChan).Should(BeClosed())
	})

	It("requests the correct url", func() {
		server := ghttp.NewServer()
		Endpoint = server.URL()
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/sorted_shops/1.jsonl"),
				ghttp.RespondWith(http.StatusOK, ""),
			),
		)

		Find(1)
	})
})
