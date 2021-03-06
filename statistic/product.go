package statistic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var Endpoint string
var Username string
var Password string

type StatisticProduct struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Ean      string `json:"product_ean"`
	Brand    string `json:"product_brand"`
	ShopCode string `json:"shop_code"`
	MaxCPO   int    `json:"max_cpo"`
	Category string `json:"product_category_name"`
	Price    int    `json:"product_price"`
	Statistics
}

func (stats Statistics) TotalCosts() float64 {
	var total float64
	for _, stat := range stats {
		total += stat.Costs
	}
	return total
}

// Totals up the traffic from the daily stats record and
// converts to int cause traffic would be a whole number
func (stats Statistics) TotalTraffic() int {
	var total int
	for _, stat := range stats {
		total += int(stat.Traffic)
	}
	return total
}

func (stats Statistics) TotalProfit() float64 {
	return stats.TotalCMargin() - stats.TotalCosts()
}

func (stats Statistics) TotalCMargin() float64 {
	var total float64
	for _, stat := range stats {
		total += stat.CMargin
	}
	return total
}
func (stats Statistics) TotalCEXAmount() float64 {
	var total float64
	for _, stat := range stats {
		total += stat.CEXAmount
	}
	return total
}

type Statistics []*Statistic

type Statistic struct {
	BounceRate              float64 `json:"bounce_rate"`
	CCPO                    float64 `json:"ccpo"`
	CEXAmount               float64 `json:"cex_amount"`
	CMargin                 float64 `json:"cmargin"`
	Contribution            float64 `json:"contribution"`
	Conversion              float64 `json:"conversion"`
	Costs                   float64 `json:"costs"`
	CPO                     float64 `json:"cpo"`
	CROAS                   float64 `json:"croas"`
	CROI                    float64 `json:"croi"`
	ECPC                    float64 `json:"ecpc"`
	OrderAmountExcludingTax float64 `json:"order_amount_excluding_tax"`
	OrderAmountIncludingTax float64 `json:"order_amount_including_tax"`
	Orders                  float64 `json:"orders"`
	Quantity                float64 `json:"quantity"`
	Roas                    float64 `json:"roas"`
	Traffic                 float64 `json:"traffic"`
	Tos                     float64 `json:"tos"`
	Contributed             float64 `json:"contributed"`
	Views                   float64 `json:"views"`
	Assists                 float64 `json:"assists"`
	AssistRatio             float64 `json:"assist_ratio"`
	TimeId                  string  `json:"time_id"`
}

type DailyProductsQuery struct {
	ShopId      int
	PublisherId int
	StartDate   time.Time
	StopDate    time.Time
	ShopCodes   *[]string
	Limit       *int
	Skip        *int
}

func (productsQuery *DailyProductsQuery) RawQuery() string {
	query := url.Values{}

	startDateKey, stopDateKey := DailyTimeId(productsQuery.StartDate), DailyTimeId(productsQuery.StopDate)
	if startDateKey == stopDateKey {
		query.Add("time_id", fmt.Sprintf("%s", startDateKey))
	} else {
		query.Add("time_id", fmt.Sprintf("%s:%s", startDateKey, stopDateKey))
	}

	if productsQuery.ShopCodes != nil {
		for _, shopCode := range *productsQuery.ShopCodes {
			query.Add("shop_codes[]", shopCode)
		}
	}

	if productsQuery.Limit != nil {
		query.Add("limit", strconv.Itoa(*productsQuery.Limit))
	}
	if productsQuery.Skip != nil {
		query.Add("skip", strconv.Itoa(*productsQuery.Skip))
	}

	return query.Encode()
}

type DailyProductFinder func(*DailyProductsQuery) ([]*StatisticProduct, error)

func FindDailyProduct(query *DailyProductsQuery) ([]*StatisticProduct, error) {
	url, err := getStatsUrl(query)
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
	stats := []*StatisticProduct{}
	err = json.NewDecoder(response.Body).Decode(&stats)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func getStatsUrl(query *DailyProductsQuery) (string, error) {
	statsUrl, err := url.Parse(fmt.Sprintf("%s/api/v1/shops/%d/publishers/%d/shop_products/statistics.json", Endpoint, query.ShopId, query.PublisherId))
	if err != nil {
		return "", err
	}
	statsUrl.RawQuery = query.RawQuery()
	return statsUrl.String(), nil
}
