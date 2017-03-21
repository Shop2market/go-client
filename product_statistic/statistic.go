package product_statistic

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var Endpoint string

type Statistic struct {
	ShopCode                string   `json:"shop_code"`
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

type TimeId string

func NewTimeId(date time.Time) TimeId {
	return TimeId(fmt.Sprintf("%d%02d%02d", date.Year(), date.Month(), date.Day()))
}

func FindForTimePeriod(shopId, publisherId int, timePeriod int) *StatisticIterator {
	startTime := NewTimeId(time.Now().AddDate(0, 0, (timePeriod+1)*-1))
	stopTime := NewTimeId(time.Now().AddDate(0, 0, -1))
	iterator, _ := Find(shopId, publisherId, startTime, stopTime)
	return iterator
}

func Find(shopId, publisherId int, startTimeId, stopTimeId TimeId) (*StatisticIterator, error) {
	url, err := buildUrl(shopId, publisherId, startTimeId, stopTimeId, nil)
	if err != nil {
		return nil, err
	}
	return find(url)
}

func FindForShopCodes(shopId, publisherId int, startTimeId, stopTimeId TimeId, shopCodes []string) (*StatisticIterator, error) {
	url, err := buildUrl(shopId, publisherId, startTimeId, stopTimeId, shopCodes)
	if err != nil {
		return nil, err
	}
	return find(url)
}

func find(url string) (*StatisticIterator, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		response.Body.Close()
		return nil, fmt.Errorf("Response not OK\nServer Body:\n%s", string(body))
	}
	return &StatisticIterator{Decoder: json.NewDecoder(response.Body), Closer: response.Body}, nil
}
func buildUrl(shopId, publisherId int, startTimeId, stopTimeId TimeId, shopCodes []string) (string, error) {
	resourceUrl, err := url.Parse(fmt.Sprintf("%s/shops/%d/publishers/%d/statistics", Endpoint, shopId, publisherId))
	if err != nil {
		return "", err
	}
	query := url.Values{}
	query.Add("start", string(startTimeId))
	query.Add("stop", string(stopTimeId))
	if shopCodes != nil {
		for _, shopCode := range shopCodes {
			query.Add("shop_code[]", shopCode)
		}
	}

	resourceUrl.RawQuery = query.Encode()
	return resourceUrl.String(), nil
}

type StatisticIterator struct {
	io.Closer
	*json.Decoder
}

func (iterator *StatisticIterator) More() bool {
	return iterator.Decoder.More()
}

func (iterator *StatisticIterator) Next() (*Statistic, error) {
	stats := &Statistic{}
	err := iterator.Decoder.Decode(stats)
	if err == io.EOF {
		err = nil
		iterator.Close()
		return nil, nil
	}
	return stats, err
}

func (iterator *StatisticIterator) All() ([]*Statistic, error) {
	stats := []*Statistic{}
	for iterator.More() {
		stat, err := iterator.Next()
		if err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}
	return stats, nil
}
