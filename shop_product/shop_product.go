package shop_product

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var Endpoint string
var Username string
var Password string

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

func Find(shopId int, shopCodes []string) ([]*Product, error) {
	url, err := shopUrl(shopId, shopCodes)
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

func shopUrl(shopId int, shopCodes []string) (string, error) {
	s2mUrl, err := url.Parse(fmt.Sprintf("%s/api/v1/shops/%d/shop_products.json", Endpoint, shopId))
	if err != nil {
		return "", err
	}
	query := s2mUrl.Query()

	for _, shopCode := range shopCodes {
		query.Add("shop_codes[]", shopCode)
	}
	s2mUrl.RawQuery = query.Encode()
	return s2mUrl.String(), nil
}
