package ad

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cthulhu/go-steun/time_id"
)

var Endpoint string

// TODO: move to its own package
type Statistic struct {
	Traffic                 float64  `json:"traffic"`
	Orders                  float64  `json:"orders"`
	Costs                   float64  `json:"costs"`
	OrderAmountExcludingTax float64  `json:"order_amount_excluding_tax"`
	OrderAmountIncludingTax float64  `json:"order_amount_including_tax"`
	Assists                 float64  `json:"assists"`
	Quantity                float64  `json:"quantity"`
	CEXAmount               float64  `json:"cex_amount"`
	CMargin                 float64  `json:"cmargin"`
	Profit                  float64  `json:"profit"`
	ROI                     *float64 `json:"roi"`
	ECPC                    *float64 `json:"ecpc"`
	MaxCPC                  *float64 `json:"max_cpc"`
}

type UserFields struct {
	User1  *string `json:"User1"`
	User2  *string `json:"User2"`
	User3  *string `json:"User3"`
	User4  *string `json:"User4"`
	User5  *string `json:"User5"`
	User6  *string `json:"User6"`
	User7  *string `json:"User7"`
	User8  *string `json:"User8"`
	User9  *string `json:"User9"`
	User10 *string `json:"User10"`
	User11 *string `json:"User11"`
	User12 *string `json:"User12"`
	User13 *string `json:"User13"`
	User14 *string `json:"User14"`
	User15 *string `json:"User15"`
	User16 *string `json:"User16"`
	User17 *string `json:"User17"`
	User18 *string `json:"User18"`
	User19 *string `json:"User19"`
	User21 *string `json:"User21"`
	User22 *string `json:"User22"`
	User23 *string `json:"User23"`
	User24 *string `json:"User24"`
	User25 *string `json:"User25"`
	User26 *string `json:"User26"`
	User27 *string `json:"User27"`
	User28 *string `json:"User28"`
	User29 *string `json:"User29"`
	User31 *string `json:"User31"`
	User32 *string `json:"User32"`
	User33 *string `json:"User33"`
	User34 *string `json:"User34"`
	User35 *string `json:"User35"`
	User36 *string `json:"User36"`
	User37 *string `json:"User37"`
	User38 *string `json:"User38"`
	User39 *string `json:"User39"`
	User41 *string `json:"User41"`
	User42 *string `json:"User42"`
	User43 *string `json:"User43"`
	User44 *string `json:"User44"`
	User45 *string `json:"User45"`
	User46 *string `json:"User46"`
	User47 *string `json:"User47"`
	User48 *string `json:"User48"`
	User49 *string `json:"User49"`
	User50 *string `json:"User50"`
}

type Aggregation struct {
	Key string
	*Statistic
}

func (agg *Aggregation) Aggregate(ad *Ad) {
	if ad.Statistic != nil {
		agg.Traffic += ad.Traffic
		agg.Orders += ad.Orders
		agg.Costs += ad.Costs
		agg.Assists += ad.Assists
		agg.OrderAmountExcludingTax += ad.OrderAmountExcludingTax
		agg.OrderAmountIncludingTax += ad.OrderAmountIncludingTax
		agg.Quantity += ad.Quantity
		agg.CEXAmount += ad.CEXAmount
		agg.CMargin += ad.CMargin
	}
}

// TODO: extend with all the ads properties
type Ad struct {
	ShopCode string `json:"shop_code"`
	*Statistic
	*UserFields
}

func (ad *Ad) GetAttribute(attributeName string) (string, error) {
	if attributeName == "User1" {
		if ad.UserFields != nil && ad.User1 != nil {
			return *ad.User1, nil
		}
		return "", nil
	}
	return "", fmt.Errorf("Unsupported attribute %s", attributeName)
}

type AggregationsByKey map[string]map[string]*Aggregation

func FindAggregationsWithTimePeriod(shopID, publisherID int, aggregationKeys []string, timePeriod int) (*AggregationsByKey, error) {
	if len(aggregationKeys) == 0 || Endpoint == "" {
		return nil, nil
	}
	aggregationsByKey := make(AggregationsByKey)
	startTime := time_id.NewByDays(timePeriod * -1)
	stopTime := time_id.NewByDays(-1)

	url := buildURL(shopID, publisherID, startTime, stopTime)

	response, err := http.Get(url)
	fmt.Println(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("[FAIL]: Status: %d, Server Body: %s", response.StatusCode, string(body))
	}
	var ads []*Ad
	decoder := json.NewDecoder(response.Body)

	if err = decoder.Decode(&ads); err != nil {
		return nil, err
	}
	for _, aggregationKey := range aggregationKeys {
		aggregationsByKey[aggregationKey] = make(map[string]*Aggregation)
	}

	for _, ad := range ads {
		for _, aggregationKey := range aggregationKeys {
			key, err := ad.GetAttribute(aggregationKey)
			if err != nil {
				return nil, err
			}
			agg, exist := aggregationsByKey[aggregationKey][key]
			if !exist {
				agg = &Aggregation{Statistic: &Statistic{}}
				aggregationsByKey[aggregationKey][key] = agg
			}
			agg.Aggregate(ad)
		}
	}
	return &aggregationsByKey, nil
}

func buildURL(shopID, publisherID int, startTime, stopTime time_id.TimeId) string {
	// TODO: Impelemt via url parse
	return fmt.Sprintf("%s/shops/%d/publishers/%d/ads?start=%s&stop=%s", Endpoint, shopID, publisherID, startTime.ToStr(), stopTime.ToStr())
}
