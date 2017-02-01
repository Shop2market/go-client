package publisher_connection_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	. "github.com/Shop2market/go-client/publisher_connection"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Connection", func() {
	It("unmarshals publisher connection", func() {
		connections := []*Connection{}
		body, _ := ioutil.ReadFile("fixtures/connection_response.json")
		json.Unmarshal(body, &connections)
		Expect(connections).To(Equal([]*Connection{
			&Connection{
				ID:                 11107,
				Imported:           true,
				Connected:          true,
				ConnectionType:     "custom",
				ConnectionProvider: "beslist_marketplace",
				ConnectToLive:      true,
				ProductApiEnabled:  true,
				ConnectionDetails: &ConnectionDetails{
					ShopId:              stringPtr("46834"),
					ProductUpdateApiKey: stringPtr("6754856"),
					SellerID:            stringPtr("sellthk"),
					MarketPlaceID:       stringPtr("abc"),
					MWSToken:            stringPtr("abs"),
					PublicKey:           stringPtr("756374"),
				},
			},
		}))
	})
	It("Retrives connection settings", func() {
		content, err := ioutil.ReadFile("fixtures/connection_response.json")
		Expect(err).NotTo(HaveOccurred())

		server := ghttp.NewServer()
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyBasicAuth("username", "password"),
				ghttp.VerifyRequest("GET", "/api/v1/shops/1/publishers/2/costs.json"),
				ghttp.RespondWith(http.StatusOK, string(content)),
			),
		)
		defer server.Close()
		Username = "username"
		Password = "password"
		Endpoint = server.URL()

		connections, err := Find(&Query{ShopId: 1, PublisherId: 2})
		Expect(err).NotTo(HaveOccurred())
		Expect(server.ReceivedRequests()).Should(HaveLen(1))
		Expect(connections).To(Equal([]*Connection{
			&Connection{
				ID:                 11107,
				Imported:           true,
				Connected:          true,
				ConnectionType:     "custom",
				ConnectionProvider: "beslist_marketplace",
				ConnectToLive:      true,
				ProductApiEnabled:  true,
				ConnectionDetails: &ConnectionDetails{
					ShopId:              stringPtr("46834"),
					ProductUpdateApiKey: stringPtr("6754856"),
					SellerID:            stringPtr("sellthk"),
					MarketPlaceID:       stringPtr("abc"),
					MWSToken:            stringPtr("abs"),
					PublicKey:           stringPtr("756374"),
				},
			},
		}))

	})
})

func stringPtr(str string) *string {
	return &str
}
