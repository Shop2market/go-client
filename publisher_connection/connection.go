package publisher_connection

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type GoogleConnection struct {
	GoogleMerchantId *string `json:"google_merchant_id,omitempty"`
}
type MiabConnection struct {
	FeedID     *string `json:"feedID,omitempty"`
	CampaignID *string `json:"campaignID,omitempty"`
}
type ConnectionDetails struct {
	ShopId              *string `json:"shop_id,omitempty"`
	ProductUpdateApiKey *string `json:"product_update_api_key,omitempty"`
	PublicKey           *string `json:"public_key,omitempty"`
	PrivateKey          *string `json:"private_key,omitempty"`
	SellerID            *string `json:"seller_id,omitempty"`
	MarketPlaceID       *string `json:"marketplace_id,omitempty"`
	MWSToken            *string `json:"mws_token,omitempty"`
	MiabConnection
	GoogleConnection
}

type Connection struct {
	ID                   int    `json:"id"`
	Imported             bool   `json:"imported"`
	Connected            bool   `json:"connected"`
	ConnectionType       string `json:"connection_type"`
	ConnectionProvider   string `json:"connection_provider"`
	ConnectToLive        bool   `json:"connect_to_live"`
	ProductApiEnabled    bool   `json:"product_api_enabled"`
	PatchOnProductUpdate bool   `json:"patch_on_product_update"`
	*ConnectionDetails   `json:"connection"`
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
	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("Error fetching connections: %s", response.Status)
	}
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
