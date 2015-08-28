package statistic_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/Shop2market/go-client/statistic"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Publisher statistics", func() {
	Context("Deserialize Publisher statisitcs", func() {
		It("can read from json", func() {
			fileContent, _ := ioutil.ReadFile("fixtures/publisher_response.json")

			publisher := &Publisher{}
			Expect(json.Unmarshal(fileContent, &publisher)).ToNot(HaveOccurred())
			Expect(publisher.Statistics).To(HaveLen(31))

			Expect(publisher.Id).To(Equal(114))
			Expect(publisher.Name).To(Equal("PublisherName"))

			Expect(publisher.Statistics[0].Quantity).To(Equal(33.0))
			Expect(publisher.Statistics[0].Traffic).To(Equal(414.0))
			Expect(publisher.Statistics[0].Orders).To(Equal(27.0))
		})
	})

	Context("Find statistics for publisher", func() {
		Username = "username"
		Password = "password"

		It("Fetches statistic", func() {
			startDate := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(2015, 1, 31, 0, 0, 0, 0, time.UTC)

			content, err := ioutil.ReadFile("fixtures/publisher_response.json")
			Expect(err).NotTo(HaveOccurred())

			server := ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/shops/1/publishers/2/statistics.json", "time_id=20150101%3A20150131"),
					ghttp.VerifyBasicAuth(Username, Password),
					ghttp.RespondWith(http.StatusOK, string(content)),
				),
			)
			Endpoint = server.URL()

			stats, err := FindPublisherStatistic(&PublisherStatisticQuery{ShopId: 1, PublisherId: 2, StartDate: startDate, StopDate: endDate})
			Expect(err).NotTo(HaveOccurred())
			Expect(stats.Id).To(Equal(114))

		})
	})
})
