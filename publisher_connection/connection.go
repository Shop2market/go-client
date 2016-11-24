package publisher_connection

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type ConnectionDetails struct {
	ShopId              *string `json:"shop_id,omitempty"`
	ProductUpdateApiKey *string `json:"product_update_api_key,omitempty"`
	PublicKey           *string `json:"public_key,omitempty"`
	Signature           *string `json:"signature,omitempty"`
	SignatureDate       *string `json:"signature_date,omitempty"`
}

type Connection struct {
	ID                 int    `json:"id"`
	Imported           bool   `json:"imported"`
	Connected          bool   `json:"connected"`
	ConnectionType     string `json:"connection_type"`
	ConnectionProvider string `json:"connection_provider"`
	ConnectToLive      bool   `json:"connect_to_live"`
	ProductApiEnabled  bool   `json:"product_api_enabled"`
	*ConnectionDetails `json:"connection"`
}

var Endpoint string
var Username string
var Password string

type Query struct {
	ShopId      int
	PublisherId int
}

func Find(query *Query) ([]*Connection, error) {
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
	connections := []*Connection{}

	if err := json.NewDecoder(response.Body).Decode(&connections); err != nil {
		return nil, err
	}
	return connections, nil
}

func apiUrl(query *Query) (string, error) {
	url, err := url.Parse(fmt.Sprintf("%s/api/v1/shops/%d/publishers/%d/costs.json", Endpoint, query.ShopId, query.PublisherId))
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

type Finder func(*Query) ([]*Connection, error)

// Use for tests to stub calls to API
func DummyFinder(query *Query) ([]*Connection, error) {
	return []*Connection{}, nil
}
