package channel_product_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/Shop2market/go-client/channel_product"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Store", func() {

	It("deserialization of channel product", func() {
		products := []*Product{}
		body, _ := ioutil.ReadFile("fixtures/channel_product_response.json")
		json.Unmarshal(body, &products)

		Expect(products).To(Equal([]*Product{
			&Product{
				Id: &ProductId{
					ShopCode:    "151656",
					ShopId:      1,
					PublisherId: 5,
				},
				Active: true,
			},
			&Product{
				Id: &ProductId{
					ShopCode:    "149350",
					ShopId:      1,
					PublisherId: 5,
				},
				Active: true,
			},
		}))
	})

	Context("requests channel products", func() {
		It("sends correct parameters", func() {
			server := NewMockedServer("fixtures/channel_product_response.json")
			Endpoint = server.URL

			Find(1, 5, 0, 10)

			Expect(server.Requests).To(HaveLen(1))

			marmosetUrl, _ := url.Parse("/shops/1/publishers/5/products?enabled=true&limit=10&skip=0")
			Expect(server.Requests[0].URL).To(Equal(marmosetUrl))
		})

		It("deserializes response", func() {
			server := NewMockedServer("fixtures/channel_product_response.json")
			Endpoint = server.URL

			Expect(Find(1, 5, 0, 10)).To(Equal([]*Product{
				&Product{
					Id: &ProductId{
						ShopCode:    "151656",
						ShopId:      1,
						PublisherId: 5,
					},
					Active: true,
				},
				&Product{
					Id: &ProductId{
						ShopCode:    "149350",
						ShopId:      1,
						PublisherId: 5,
					},
					Active: true,
				},
			}))

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
