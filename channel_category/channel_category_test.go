package channel_category_test

import (
	"io/ioutil"
	"net/http"

	. "github.com/Shop2market/go-client/channel_category"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("ChannelCategory", func() {
	It("Retrives categories", func() {
		content, err := ioutil.ReadFile("fixtures/channel_categories_index_response.json")
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

		categories, err := Find(&Query{ShopId: 1, PublisherId: 2})
		categories.BuildPaths()
		Expect(err).NotTo(HaveOccurred())
		Expect(categories).To(HaveLen(5))
		Expect(categories).To(ContainElement(&ChannelCategory{
			Name:       "Consumer Electronics",
			ParentID:   0,
			ExternalID: "12345_1",
			Path:       "Consumer Electronics",
			ID:         1,
		}))
		Expect(categories).To(ContainElement(&ChannelCategory{
			Name:       "Home and Living",
			ParentID:   0,
			ExternalID: "12345_2",
			Path:       "Home and Living",
			ID:         2,
		}))
		Expect(categories).To(ContainElement(&ChannelCategory{
			Name:       "Lamps",
			ParentID:   2,
			ExternalID: "12345_5",
			Path:       "Home and Living -> Lamps",
			ID:         5,
		}))
	})
})
