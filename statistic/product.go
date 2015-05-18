package statistic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var Endpoint string
var Username string
var Password string

type StatisticProduct struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Ean        string `json:"product_ean"`
	Brand      string `json:"product_brand"`
	ShopCode   string `json:"shop_code"`
	MaxCPO     int    `json:"max_cpo"`
	Category   string `json:"product_category_name"`
	Price      int    `json:"product_price"`
	Statistics []*DailyStatistic
}

func (statsProduct *StatisticProduct) TotalRoi() float64 {
	var total float64
	for _, stat := range statsProduct.Statistics {
		total += stat.Roi
	}
	return total
}
func (statsProduct *StatisticProduct) TotalCosts() float64 {
	var total float64
	for _, stat := range statsProduct.Statistics {
		total += stat.Costs
	}
	return total
}

// Totals up the traffic from the daily stats record and
// converts to int cause traffic would be a whole number
func (statsProduct *StatisticProduct) TotalTraffic() int {
	var total int
	for _, stat := range statsProduct.Statistics {
		total += int(stat.Traffic)
	}
	return total
}
func (statsProduct *StatisticProduct) TotalContributedProfit() float64 {
	var total float64
	for _, stat := range statsProduct.Statistics {
		total += stat.ContributedProfit
	}
	return total
}

type DailyStatistic struct {
	BounceRate              float64 `json:"bounce_rate"`
	CCPO                    float64 `json:"ccpo"`
	CEXAmount               float64 `json:"cex_amount"`
	ContributedProfit       float64 `json:"contributed_profit"`
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
	Roi                     float64 `json:"roi"`
	Traffic                 float64 `json:"traffic"`
	Tos                     float64 `json:"tos"`
	Contributed             float64 `json:"contributed"`
	Views                   float64 `json:"views"`
	Assists                 float64 `json:"assists"`
	AssistRatio             float64 `json:"assist_ratio"`
	TimeId                  string  `json:"time_id"`
}

func FindDailyProduct(shopId, publisherId int, startDate, stopDate string, shopCodes []string) ([]*StatisticProduct, error) {
	url, err := getStatsUrl(shopId, publisherId, startDate, stopDate, shopCodes)
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

func getStatsUrl(shopId, publisherId int, startDate, stopDate string, shopCodes []string) (string, error) {
	statsUrl, err := url.Parse(fmt.Sprintf("%s/api/v1/shops/%d/publishers/%d/shop_products/statistics.json", Endpoint, shopId, publisherId))
	if err != nil {
		return "", err
	}
	query := statsUrl.Query()
	query.Add("time_id", fmt.Sprintf("%s:%s", startDate, stopDate))

	for _, shopCode := range shopCodes {
		query.Add("shop_codes[]", shopCode)
	}
	statsUrl.RawQuery = query.Encode()
	return statsUrl.String(), nil
}
