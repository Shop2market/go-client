package publisher_connection_test

import (
	"encoding/json"
	"io/ioutil"

	. "github.com/Shop2market/go-client/publisher_connection"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Connection", func() {
	It("unmarshals of publisher connection", func() {
		connections := []Connection{}
		body, _ := ioutil.ReadFile("fixtures/connection_response.json")
		json.Unmarshal(body, &connections)
		Expect(connections).To(Equal([]Connection{
			Connection{
				ID:                 11107,
				Imported:           true,
				Connected:          true,
				ConnectionType:     "custom",
				ConnectionProvider: "beslist_marketplace",
				ConnectionDetails: &ConnectionDetails{
					Enabled:             true,
					ShopId:              stringPtr("46834"),
					ProductUpdateApiKey: stringPtr("6754856"),
				},
			},
		}))
	})
})

func stringPtr(str string) *string {
	return &str
}
