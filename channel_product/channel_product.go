package channel_product

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

var Endpoint string

type Product struct {
	Id     *ProductId `json:"product_id"`
	Active bool       `json:"active"`
}
type ProductId struct {
	ShopCode    string `json:"shop_code"`
	ShopId      int    `json:"shop_id"`
	PublisherId int    `json:"publisher_id"`
}

func Find(shopId, publisherId, skip, limit int) ([]*Product, error) {
	productUrl, err := url.Parse(fmt.Sprintf("%s/shops/%d/publishers/%d/products", Endpoint, shopId, publisherId))
	if err != nil {
		return nil, err
	}
	query := productUrl.Query()
	query.Add("enabled", "true")
	query.Add("limit", strconv.Itoa(limit))
	query.Add("skip", strconv.Itoa(skip))

	productUrl.RawQuery = query.Encode()

	response, err := http.Get(productUrl.String())
	if err != nil {
		return nil, err
	}
	products := []*Product{}
	err = json.NewDecoder(response.Body).Decode(&products)
	if err != nil {
		return nil, err
	}
	return products, nil
}
