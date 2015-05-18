package statistic_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	. "github.com/Shop2market/go-client/statistic"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Product statistics", func() {
	Username = "username"
	Password = "password"
	Context("Request for daily statistics", func() {
		startDate := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2015, 1, 31, 0, 0, 0, 0, time.UTC)

		It("calls the statitics api with correct params", func() {
			server := NewMockedServer("fixtures/shop_code_statistics_response.json")
			Endpoint = server.URL
			FindDailyProduct(1, 2, startDate, endDate, []string{"a001", "b002", "c003"})
			Expect(server.Requests).To(HaveLen(1))

			expectedUrl, _ := url.ParseRequestURI("/api/v1/shops/1/publishers/2/shop_products/statistics.json?shop_codes%5B%5D=a001&shop_codes%5B%5D=b002&shop_codes%5B%5D=c003&time_id=20150101%3A20150131")

			Expect(expectedUrl).To(Equal(server.Requests[0].URL))

			username, password, _ := server.Requests[0].BasicAuth()
			Expect(username).To(Equal(Username))
			Expect(password).To(Equal(Password))
		})

		It("parses response", func() {
			server := NewMockedServer("fixtures/shop_code_statistics_response.json")
			Endpoint = server.URL
			stats, err := FindDailyProduct(1, 2, startDate, endDate, []string{"a001", "b002", "c003"})
			Expect(err).NotTo(HaveOccurred())
			Expect(stats).To(HaveLen(2))
			Expect(stats[0].Id).To(Equal(98769059))
			Expect(stats[0].Name).To(Equal("SAMSUNG Galaxy S4 mini"))
		})

		It("returns error if json broken", func() {
			server := NewMockedServer("fixtures/broken_stats_response.json")
			Endpoint = server.URL
			_, err := FindDailyProduct(1, 2, startDate, endDate, []string{"a001", "b002", "c003"})
			Expect(err).To(HaveOccurred())
		})
	})
	Context("Deserialize StatisticProduct", func() {
		It("can read from json", func() {
			fileContent, _ := ioutil.ReadFile("fixtures/shop_code_statistics_response.json")

			products := []*StatisticProduct{}
			Expect(json.Unmarshal(fileContent, &products)).ToNot(HaveOccurred())
			Expect(products).To(HaveLen(2))

			Expect(products[0].Id).To(Equal(98769059))
			Expect(products[0].Name).To(Equal("SAMSUNG Galaxy S4 mini"))
			Expect(products[0].Ean).To(Equal("8716406051240"))
			Expect(products[0].Brand).To(Equal("SAMSUNG"))
			Expect(products[0].ShopCode).To(Equal("151656"))
			Expect(products[0].MaxCPO).To(Equal(2105))
			Expect(products[0].Category).To(Equal("Telecom -> Mobiele telefoons & Smartphones -> Simlockvrije telefoons"))
			Expect(products[0].Price).To(Equal(22476))
			Expect(products[0].Statistics).To(HaveLen(29))

			Expect(products[0].Statistics[0].BounceRate).To(BeEquivalentTo(40))
			Expect(products[0].Statistics[0].CCPO).To(Equal(2.1999999999999997))
			Expect(products[0].Statistics[0].CEXAmount).To(Equal(86.36363636363636))
			Expect(products[0].Statistics[0].ContributedProfit).To(Equal(9.00548459804658))
			Expect(products[0].Statistics[0].Contribution).To(Equal(45.45454545454545))
			Expect(products[0].Statistics[0].Conversion).To(BeEquivalentTo(10))
			Expect(products[0].Statistics[0].Costs).To(Equal(13.12))
			Expect(products[0].Statistics[0].CPO).To(Equal(0.9999999999999999))
			Expect(products[0].Statistics[0].CROAS).To(Equal(85.36363636363637))
			Expect(products[0].Statistics[0].CROI).To(Equal(8.005484598046582))
			Expect(products[0].Statistics[0].OrderAmountExcludingTax).To(Equal(190.0))
			Expect(products[0].Statistics[0].OrderAmountIncludingTax).To(Equal(229.0))
			Expect(products[0].Statistics[0].Orders).To(Equal(1.0))
			Expect(products[0].Statistics[0].Quantity).To(Equal(1.0))
			Expect(products[0].Statistics[0].Roas).To(Equal(189.00000000000003))
			Expect(products[0].Statistics[0].Roi).To(Equal(8.00548459804658))
			Expect(products[0].Statistics[0].Traffic).To(Equal(10.0))
			Expect(products[0].Statistics[0].Tos).To(Equal(50.475))
			Expect(products[0].Statistics[0].Contributed).To(Equal(0.45454545454545453))
			Expect(products[0].Statistics[0].Views).To(Equal(0.0))
			Expect(products[0].Statistics[0].Assists).To(Equal(1.0))
			Expect(products[0].Statistics[0].AssistRatio).To(BeEquivalentTo(10.0))
			Expect(products[0].Statistics[0].TimeId).To(Equal("20150401"))
		})
	})

	Context("Statistic totals", func() {
		It("should sum up Roi", func() {
			product := &StatisticProduct{
				Statistics: []*DailyStatistic{
					&DailyStatistic{
						Roi: 1,
					},
					&DailyStatistic{
						Roi: 2.2,
					},
				},
			}
			Expect(product.TotalRoi()).To(Equal(3.2))
		})

		It("should sum up Costs", func() {
			product := &StatisticProduct{
				Statistics: []*DailyStatistic{
					&DailyStatistic{
						Costs: 1,
					},
					&DailyStatistic{
						Costs: 2.2,
					},
				},
			}
			Expect(product.TotalCosts()).To(Equal(3.2))
		})
		It("should sum up ContributedProfit", func() {
			product := &StatisticProduct{
				Statistics: []*DailyStatistic{
					&DailyStatistic{
						ContributedProfit: 1,
					},
					&DailyStatistic{
						ContributedProfit: 2.2,
					},
				},
			}
			Expect(product.TotalContributedProfit()).To(Equal(3.2))
		})

		It("should sum up Traffic", func() {
			product := &StatisticProduct{
				Statistics: []*DailyStatistic{
					&DailyStatistic{
						Traffic: 1.0,
					},
					&DailyStatistic{
						Traffic: 2.1,
					},
				},
			}
			Expect(product.TotalTraffic()).To(Equal(3))
		})
	})
})

type MockedServer struct {
	*httptest.Server
	Requests []*http.Request
	Response []byte
}

func NewMockedServer(responseFileName string) *MockedServer {
	ser := &MockedServer{}
	ser.Server = httptest.NewServer(ser)
	ser.Requests = []*http.Request{}
	data, _ := ioutil.ReadFile(responseFileName)
	ser.Response = data
	return ser
}
func (ser *MockedServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ser.Requests = append(ser.Requests, req)
	resp.WriteHeader(200)
	resp.Write(ser.Response)
}
