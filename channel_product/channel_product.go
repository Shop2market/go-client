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
	Id                  *ProductId `json:"product_id"`
	Active              bool       `json:"active"`
	ManuallyDeactivated bool       `json:"manually_deactivated"`
}
type ProductId struct {
	ShopCode    string `json:"shop_code"`
	ShopId      int    `json:"shop_id"`
	PublisherId int    `json:"publisher_id"`
}

type ProductsQuery struct {
	ShopId              int
	PublisherId         int
	Skip                *int
	Limit               *int
	Enabled             *bool
	ManuallyDeactivated *bool
	ShopCodes           *[]string
}

func (productsQuery *ProductsQuery) RawQuery() string {
	query := url.Values{}
	if productsQuery.Enabled != nil {
		if *productsQuery.Enabled {
			query.Add("enabled", "true")
		} else {
			query.Add("enabled", "false")
		}
	}

	if productsQuery.ManuallyDeactivated != nil {
		if *productsQuery.ManuallyDeactivated {
			query.Add("manually_deactivated", "true")
		} else {
			query.Add("manually_deactivated", "false")
		}
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

func Find(productsQuery *ProductsQuery) ([]*Product, error) {
	productUrl, err := url.Parse(fmt.Sprintf("%s/shops/%d/publishers/%d/products", Endpoint, productsQuery.ShopId, productsQuery.PublisherId))
	if err != nil {
		return nil, err
	}
	productUrl.RawQuery = productsQuery.RawQuery()

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
