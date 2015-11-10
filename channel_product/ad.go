package channel_product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Ad struct {
	*Product
	ProductName     string   `json:"product_name"`
	PictureLink     string   `json:"picture_link"`
	Deeplink        string   `json:"deeplink"`
	ShopCode        string   `json:"shop_code"`
	ProductEan      string   `json:"product_ean"`
	Category        string   `json:"category"`
	Brand           string   `json:"brand"`
	DeliveryPeriod  string   `json:"delivery_period"`
	ProductsInStock int      `json:"products_in_stock"`
	Price           int      `json:"price"`
	Traffic         int      `json:"traffic"`
	Profit          float64  `json:"profit"`
	Costs           float64  `json:"costs"`
	ROI             *float64 `json:"roi"`
	Revenue         float64  `json:"revenue"`
	Margin          float64  `json:"margin"`
}

type AdQuery struct {
	*ProductsQuery
	StartTimeId   string
	EndTimeId     string
	OnlyWithStats bool
}

func (adQuery *AdQuery) RawQuery() string {
	values, _ := url.ParseQuery(adQuery.ProductsQuery.RawQuery())
	values.Add("start", adQuery.StartTimeId)
	values.Add("end", adQuery.EndTimeId)
	if adQuery.OnlyWithStats {
		values.Add("only_with_stats", "true")
	}
	return values.Encode()
}

func FindAds(productsQuery *AdQuery) ([]*Ad, error) {
	productUrl, err := buildAdsQueryUrl(productsQuery)
	response, err := http.Get(productUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("Failed fetching ads: %s, code %d\n%s", productUrl, response.StatusCode, string(body))
	}

	ads := []*Ad{}
	err = json.NewDecoder(response.Body).Decode(&ads)
	if err != nil {
		return nil, err
	}
	return ads, nil
}
func buildAdsQueryUrl(adsQuery *AdQuery) (string, error) {
	productUrl, err := url.Parse(fmt.Sprintf("%s/shops/%d/publishers/%d/ads", Endpoint, adsQuery.ShopId, adsQuery.PublisherId))
	if err != nil {
		return "", err
	}
	productUrl.RawQuery = adsQuery.RawQuery()
	return productUrl.String(), nil
}
