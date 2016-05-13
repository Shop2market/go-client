package publisher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var Endpoint string
var Username string
var Password string

type Query struct {
	ShopId      int
	PublisherId int
}

type Publisher struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	ProductApiEnabled bool   `json:"product_api_enabled"`
}

// Finder - Main find functor, can be overloaded for stubs or assigned with package Find function
type Finder func(shopId, publisherId int) (*Publisher, error)

func Find(query *Query) (*Publisher, error) {
	url, err := apiUrl(query)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(Username, Password)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	publisher := Publisher{}

	if err := json.NewDecoder(response.Body).Decode(&publisher); err != nil {
		return nil, err
	}
	return &publisher, nil
}

func apiUrl(query *Query) (string, error) {
	url, err := url.Parse(fmt.Sprintf("%s/api/v1/shops/%d/publishers/%d.json", Endpoint, query.ShopId, query.PublisherId))
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
