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
	Id          *ProductId `json:"product_id"`
	Active      bool       `json:"active"`
	Enabled     bool       `json:"enabled"`
	ManuallySet *bool      `json:"manually_set,omitempty"`
}

type ProductId struct {
	ShopCode    string `json:"shop_code"`
	ShopId      int    `json:"shop_id"`
	PublisherId int    `json:"publisher_id"`
}

type ProductsQuery struct {
	ShopId      int
	PublisherId int
	Skip        *int
	Limit       *int
	Active      *bool
	Enabled     *bool
	ManuallySet *bool
	ShopCodes   *[]string
}

func (productsQuery *ProductsQuery) RawQuery() string {
	query := url.Values{}
	if productsQuery.Active != nil {
		if *productsQuery.Active {
			query.Add("active", "true")
		} else {
			query.Add("active", "false")
		}
	}
	if productsQuery.Enabled != nil {
		if *productsQuery.Enabled {
			query.Add("enabled", "true")
		} else {
			query.Add("enabled", "false")
		}
	}

	if productsQuery.ManuallySet != nil {
		if *productsQuery.ManuallySet {
			query.Add("manually_set", "true")
		} else {
			query.Add("manually_set", "false")
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

type Finder func(*ProductsQuery) ([]*Product, error)

func Find(productsQuery *ProductsQuery) ([]*Product, error) {
	productUrl, err := buildQueryUrl(productsQuery)
	response, err := http.Get(productUrl)
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

func buildQueryUrl(productsQuery *ProductsQuery) (string, error) {
	productUrl, err := url.Parse(fmt.Sprintf("%s/shops/%d/publishers/%d/products", Endpoint, productsQuery.ShopId, productsQuery.PublisherId))
	if err != nil {
		return "", err
	}
	productUrl.RawQuery = productsQuery.RawQuery()
	return productUrl.String(), nil
}
