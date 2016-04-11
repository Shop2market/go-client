package publisher_connection

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type ConnectionDetails struct {
	Enabled             bool    `json:"enabled"`
	ShopId              *string `json:"shop_id,omitempty"`
	ProductUpdateApiKey *string `json:"product_update_api_key"`
}

type Connection struct {
	ID                 int    `json:"id"`
	Imported           bool   `json:"imported"`
	Connected          bool   `json:"connected"`
	ConnectionType     string `json:"connection_type"`
	ConnectionProvider string `json:"connection_provider"`
	*ConnectionDetails `json:"connection"`
}

var Endpoint string
var Username string
var Password string

type Query struct {
	ShopId      int
	PublisherId int
}

func Find(query *Query) ([]Connection, error) {
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
	connections := []Connection{}

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
