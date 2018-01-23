package ad_test

import (
	"fmt"
	"net/http"

	"github.com/cthulhu/go-steun/fixture"
	"github.com/cthulhu/go-steun/time_id"

	. "github.com/Shop2market/go-client/shop/publisher/ad"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Ad", func() {
	Context("FindAggregationsWithTimePeriod", func() {
		It("returns nothing if no aggregation keys", func() {
			agg, err := FindAggregationsWithTimePeriod(1, 1, []string{}, 30)
			Expect(err).NotTo(HaveOccurred())
			Expect(agg).To(BeNil())

			agg, err = FindAggregationsWithTimePeriod(1, 1, nil, 30)
			Expect(err).NotTo(HaveOccurred())
			Expect(agg).To(BeNil())

		})
		It("returns aggregated by user1", func() {

			startTime := time_id.NewByDays(-30)
			stopTime := time_id.NewByDays(-1)
			queryString := fmt.Sprintf("start=%s&stop=%s", startTime, stopTime)

			server := ghttp.NewServer()
			defer server.Close()

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/shops/1/publishers/1/ads", queryString),
					ghttp.RespondWith(http.StatusOK, fixture.Read("ads.json")),
				),
			)
			Endpoint = server.URL()

			agg, err := FindAggregationsWithTimePeriod(1, 1, []string{"User1"}, 30)
			Expect(err).NotTo(HaveOccurred())
			Expect(agg).NotTo(BeNil())
			user1Aggregations := (*agg)["User1"]
			Expect(len(user1Aggregations)).To(Equal(52))
			Expect(user1Aggregations["25761"].Traffic).To(Equal(125.0))
			Expect(user1Aggregations["25761"].Costs).To(BeNumerically("~", 194.1799, 0.01))

			Expect(user1Aggregations["27298"].Traffic).To(Equal(100.0))
			Expect(user1Aggregations["27298"].Costs).To(BeNumerically("~", 51.69, 0.01))

		})
		It("returns aggregated by user1 with negative time", func() {

			startTime := time_id.NewByDays(-29)
			stopTime := time_id.NewByDays(0)
			queryString := fmt.Sprintf("start=%s&stop=%s", startTime, stopTime)

			server := ghttp.NewServer()
			defer server.Close()

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/shops/1/publishers/1/ads", queryString),
					ghttp.RespondWith(http.StatusOK, fixture.Read("ads.json")),
				),
			)
			Endpoint = server.URL()

			agg, err := FindAggregationsWithTimePeriod(1, 1, []string{"User1"}, -30)
			Expect(err).NotTo(HaveOccurred())
			Expect(agg).NotTo(BeNil())
			user1Aggregations := (*agg)["User1"]
			Expect(len(user1Aggregations)).To(Equal(52))
			Expect(user1Aggregations["25761"].Traffic).To(Equal(125.0))
			Expect(user1Aggregations["25761"].Costs).To(BeNumerically("~", 194.1799, 0.01))

			Expect(user1Aggregations["27298"].Traffic).To(Equal(100.0))
			Expect(user1Aggregations["27298"].Costs).To(BeNumerically("~", 51.69, 0.01))

		})
	})
})
