package publisher_max_cpc_range

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type MaxCpcRange struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	ChannelCategoryID int     `json:"channel_category_id"`
	MaxCpcMin         float64 `json:"max_cpc_min"`
	MaxCpcMax         float64 `json:"max_cpc_max"`
}

var Endpoint string
var Username string
var Password string

type Query struct {
	ShopId      int
	PublisherId int
}

func Find(query *Query) ([]*MaxCpcRange, error) {
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
	maxCpcRanges := []*MaxCpcRange{}

	if err := json.NewDecoder(response.Body).Decode(&maxCpcRanges); err != nil {
		return nil, err
	}
	return maxCpcRanges, nil
}

func apiUrl(query *Query) (string, error) {
	url, err := url.Parse(fmt.Sprintf("%s/api/v1/shops/%d/publishers/%d/max_cpc_ranges.json", Endpoint, query.ShopId, query.PublisherId))
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

type Finder func(*Query) ([]*MaxCpcRange, error)

// Use for tests to stub calls to API
func DummyFinder(query *Query) ([]*MaxCpcRange, error) {
	return []*MaxCpcRange{}, nil
}
