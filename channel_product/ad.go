package channel_product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Ad struct {
	*Product
	RulesTaxonomies         map[string][]int `json:"taxonomies"`
	ProductName             string           `json:"product_name"`
	PictureLink             string           `json:"picture_link"`
	Deeplink                string           `json:"deeplink"`
	ShopCode                string           `json:"shop_code"`
	ProductEan              string           `json:"product_ean"`
	Category                string           `json:"category"`
	Brand                   string           `json:"brand"`
	DeliveryPeriod          string           `json:"delivery_period"`
	ProductsInStock         int              `json:"products_in_stock"`
	Price                   int              `json:"price"`
	Traffic                 int              `json:"traffic"`
	Profit                  float64          `json:"profit"`
	Costs                   float64          `json:"costs"`
	ROI                     *float64         `json:"roi"`
	Revenue                 float64          `json:"revenue"`
	Margin                  float64          `json:"margin"`
	OrderAmountExcludingTax float64          `json:"order_amount_excluding_tax"`
	OrderAmountIncludingTax float64          `json:"order_amount_including_tax"`
	Assists                 float64          `json:"assists"`
	Orders                  float64          `json:"orders"`
	Contributed             float64          `json:"contribution"`
	Quantity                float64          `json:"quantity"`
	QuarantinedAt           *time.Time       `json:"quarantined_at"`
	StockStatus             string           `json:"stock_status"`
	DisabledAt              *time.Time       `json:"disabled_at"`
	DeactivationReason      string           `json:"deactivation_reason"`
}

type AdQuery struct {
	*ProductsQuery
	StartTimeId   string
	EndTimeId     string
	OnlyWithStats bool
}

type AdsFinder func(productsQuery *AdQuery) ([]*Ad, error)

// Useful for tests
func DummyAdsFinder(productsQuery *AdQuery) ([]*Ad, error) {
	return []*Ad{}, nil
}

func (ad *Ad) AllTaxonomies(channelCategoryTaxonomyId string) []int {
	taxonomies := map[string][]int{}
	for id, ruleTaxonomy := range ad.RulesTaxonomies {
		taxonomies[id] = ruleTaxonomy
	}
	for id, taxonomy := range ad.Taxonomies {
		taxonomies[id] = taxonomy
	}
	if len(ad.ChannelCategoryIDs) != 0 {
		// dirty hack for a while PD-3912 is not done
		taxonomies[channelCategoryTaxonomyId] = ad.ChannelCategoryIDs
		// dirty hack for a while PD-3912 is not done
	}
	flattenedTaxonomies := []int{}
	for _, taxonomy := range taxonomies {
		flattenedTaxonomies = append(flattenedTaxonomies, taxonomy...)
	}
	return flattenedTaxonomies
}

func (ad *Ad) IsMappedToTaxonomy(taxonomyId string, isCategory bool) bool {
	taxonomies := map[string][]int{}
	for id, ruleTaxonomy := range ad.RulesTaxonomies {
		taxonomies[id] = ruleTaxonomy
	}
	for id, taxonomy := range ad.Taxonomies {
		taxonomies[id] = taxonomy
	}
	if isCategory && len(ad.ChannelCategoryIDs) != 0 {
		// dirty hack for a while PD-3912 is not done
		taxonomies[taxonomyId] = ad.ChannelCategoryIDs
		// dirty hack for a while PD-3912 is not done
	}
	return len(taxonomies[taxonomyId]) > 0
}

func (adQuery *AdQuery) RawQuery() string {
	values, _ := url.ParseQuery(adQuery.ProductsQuery.RawQuery())
	values.Add("start", adQuery.StartTimeId)
	values.Add("stop", adQuery.EndTimeId)
	if adQuery.OnlyWithStats {
		values.Add("only_with_stats", "true")
	}
	return values.Encode()
}

func FindAds(productsQuery *AdQuery) ([]*Ad, error) {
	productUrl, err := buildAdsQueryUrl(productsQuery)
	if err != nil {
		return nil, err
	}
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
