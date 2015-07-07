package shop_product

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

var Endpoint string
var Username string
var Password string

type ProductsQuery struct {
	ShopId    int
	ShopCodes *[]string
	Limit     *int
	Skip      *int
}

func (productsQuery *ProductsQuery) RawQuery() string {
	query := url.Values{}

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

type Product struct {
	ProductName    string `json:"product_name"`
	PictureLink    string `json:"picture_link"`
	Deeplink       string `json:"deeplink"`
	ShopCode       string `json:"shop_code"`
	ProductEan     string `json:"product_ean"`
	Enabled        bool   `json:"enabled"`
	ProductBrand   string `json:"product_brand"`
	DeliveryPeriod string `json:"delivery_period"`
	ProductInStock int    `json:"product_in_stock"`
	SellingPrice   int    `json:"selling_price"`
	ShopCategory   string `json:"shop_category"`
}

type Finder func(productQuery *ProductsQuery) ([]*Product, error)

func Find(productQuery *ProductsQuery) ([]*Product, error) {
	url, err := shopUrl(productQuery)
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
	products := []*Product{}
	err = json.NewDecoder(response.Body).Decode(&products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func shopUrl(productQuery *ProductsQuery) (string, error) {
	s2mUrl, err := url.Parse(fmt.Sprintf("%s/api/v1/shops/%d/shop_products.json", Endpoint, productQuery.ShopId))
	if err != nil {
		return "", err
	}
	s2mUrl.RawQuery = productQuery.RawQuery()
	return s2mUrl.String(), nil
}
