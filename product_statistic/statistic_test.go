package product_statistic_test

import (
	"io/ioutil"
	"net/http"

	. "github.com/Shop2market/go-client/product_statistic"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("AggregatedStatistic", func() {
	It("Finds statistics", func() {
		content, err := ioutil.ReadFile("fixtures/product_statistic.jsonl")
		Expect(err).NotTo(HaveOccurred())

		server := ghttp.NewServer()
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/shops/155/publishers/4/statistics", "start=20150901&stop=20150912"),
				ghttp.RespondWith(http.StatusOK, string(content)),
			),
		)
		Endpoint = server.URL()
		iterator, err := Find(155, 4, "20150901", "20150912")
		Expect(err).NotTo(HaveOccurred())

		Expect(iterator.More()).To(BeTrue())
		Expect(iterator.Next()).To(Equal(&Statistic{
			ShopCode: "100337",
			Quantity: 1,
		}))

		Expect(iterator.More()).To(BeTrue())
		roi := -1.0
		Expect(iterator.Next()).To(Equal(&Statistic{
			ShopCode: "100964",
			Traffic:  2,
			Costs:    0.22,
			Profit:   -0.22,
			ROI:      &roi,
		}))

		Expect(iterator.More()).To(BeTrue())
		Expect(iterator.Next()).NotTo(BeNil())
		Expect(iterator.More()).To(BeTrue())
		Expect(iterator.Next()).NotTo(BeNil())
		Expect(iterator.More()).To(BeTrue())
		Expect(iterator.Next()).NotTo(BeNil())
		Expect(iterator.More()).To(BeTrue())
		Expect(iterator.Next()).NotTo(BeNil())

		Expect(iterator.More()).To(BeFalse())
		Expect(iterator.Next()).To(BeNil())
		Expect(iterator.More()).To(BeFalse())
		Expect(iterator.Next()).To(BeNil())
	})
	It("Finds statistics, with shopCodes", func() {
		content, err := ioutil.ReadFile("fixtures/product_statistic.jsonl")
		Expect(err).NotTo(HaveOccurred())

		server := ghttp.NewServer()
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/shops/155/publishers/4/statistics", "start=20150901&stop=20150912&shop_code[]=001&shop_code[]=002"),
				ghttp.RespondWith(http.StatusOK, string(content)),
			),
		)
		Endpoint = server.URL()
		iterator, err := FindForShopCodes(155, 4, "20150901", "20150912", []string{"001", "002"})
		Expect(err).NotTo(HaveOccurred())

		Expect(iterator.More()).To(BeTrue())
		Expect(iterator.Next()).To(Equal(&Statistic{
			ShopCode: "100337",
			Quantity: 1,
		}))
	})
})
