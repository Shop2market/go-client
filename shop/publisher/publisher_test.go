package publisher_test

import (
	"io/ioutil"
	"net/http"

	. "github.com/Shop2market/go-client/shop/publisher"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Shop/Publisher", func() {
	It("Fetches publisher", func() {
		content, err := ioutil.ReadFile("fixtures/shop_publisher.json")
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

		publisher, err := Find(&Query{ShopId: 1, PublisherId: 2})
		Expect(err).NotTo(HaveOccurred())
		Expect(publisher).NotTo(BeNil())
		Expect(publisher.ID).To(Equal(17))
		Expect(publisher.Name).To(Equal("Kieskeurig.nl"))
		Expect(publisher.ProductApiEnabled).To(Equal(false))
		Expect(publisher.ExportMappedInAdcurve).To(Equal(false))
		Expect(publisher.TipTypes).To(HaveLen(1))
		Expect(publisher.TipTypes[0]).To(Equal("ROI"))

	})
})
