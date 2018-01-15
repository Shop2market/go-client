package product_statistic_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/Shop2market/go-client/product_statistic"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("AggregatedStatistic", func() {
	It("Finds statistics with timperiod", func() {
		startTime := NewTimeId(time.Now().AddDate(0, 0, -30))
		stopTime := NewTimeId(time.Now().AddDate(0, 0, -1))
		queryString := fmt.Sprintf("start=%s&stop=%s", startTime, stopTime)

		content, err := ioutil.ReadFile("fixtures/product_statistic.jsonl")
		Expect(err).NotTo(HaveOccurred())

		server := ghttp.NewServer()
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/shops/155/publishers/4/statistics", queryString),
				ghttp.RespondWith(http.StatusOK, string(content)),
			),
		)
		Endpoint = server.URL()

		iterator := FindForTimePeriod(155, 4, 30)
		Expect(err).NotTo(HaveOccurred())

		Expect(iterator.More()).To(BeTrue())
		Expect(iterator.Next()).To(Equal(&Statistic{
			ShopCode: "100337",
			Quantity: 1,
		}))

		Expect(iterator.More()).To(BeTrue())
		roi := -1.0
		maxCPC := 0.02
		eCPC := 0.1
		Expect(iterator.Next()).To(Equal(&Statistic{
			ShopCode: "100964",
			Traffic:  2,
			Costs:    0.22,
			Profit:   -0.22,
			MaxCPC:   &maxCPC,
			ECPC:     &eCPC,
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
	It("PD-4406: Finds statistics with timperiod - until today", func() {
		startTime := NewTimeId(time.Now().AddDate(0, 0, -29))
		stopTime := NewTimeId(time.Now())
		queryString := fmt.Sprintf("start=%s&stop=%s", startTime, stopTime)

		content, err := ioutil.ReadFile("fixtures/product_statistic.jsonl")
		Expect(err).NotTo(HaveOccurred())

		server := ghttp.NewServer()
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/shops/155/publishers/4/statistics", queryString),
				ghttp.RespondWith(http.StatusOK, string(content)),
			),
		)
		Endpoint = server.URL()

		iterator := FindForTimePeriod(155, 4, -30)
		Expect(err).NotTo(HaveOccurred())

		Expect(iterator.More()).To(BeTrue())
		Expect(iterator.Next()).To(Equal(&Statistic{
			ShopCode: "100337",
			Quantity: 1,
		}))

		Expect(iterator.More()).To(BeTrue())
		roi := -1.0
		maxCPC := 0.02
		eCPC := 0.1
		Expect(iterator.Next()).To(Equal(&Statistic{
			ShopCode: "100964",
			Traffic:  2,
			Costs:    0.22,
			Profit:   -0.22,
			MaxCPC:   &maxCPC,
			ECPC:     &eCPC,
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
		maxCPC := 0.02
		eCPC := 0.1
		Expect(iterator.Next()).To(Equal(&Statistic{
			ShopCode: "100964",
			Traffic:  2,
			Costs:    0.22,
			Profit:   -0.22,
			MaxCPC:   &maxCPC,
			ECPC:     &eCPC,
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
