package statistic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Publisher struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Statistics
}

type PublisherStatisticQuery struct {
	ShopId      int
	PublisherId int
	StartDate   time.Time
	StopDate    time.Time
}

func (statsQuery *PublisherStatisticQuery) RawQuery() string {
	query := url.Values{}

	startDateKey, stopDateKey := DailyTimeId(statsQuery.StartDate), DailyTimeId(statsQuery.StopDate)
	if startDateKey == stopDateKey {
		query.Add("time_id", fmt.Sprintf("%s", startDateKey))
	} else {
		query.Add("time_id", fmt.Sprintf("%s:%s", startDateKey, stopDateKey))
	}

	return query.Encode()
}

func FindPublisherStatistic(query *PublisherStatisticQuery) (*Publisher, error) {
	url, err := getPublisherStatsUrl(query)
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
	stats := &Publisher{}
	err = json.NewDecoder(response.Body).Decode(&stats)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func getPublisherStatsUrl(query *PublisherStatisticQuery) (string, error) {
	statsUrl, err := url.Parse(fmt.Sprintf("%s/api/v1/shops/%d/publishers/%d/statistics.json", Endpoint, query.ShopId, query.PublisherId))
	if err != nil {
		return "", err
	}
	statsUrl.RawQuery = query.RawQuery()
	return statsUrl.String(), nil
}
