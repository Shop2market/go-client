package channel_product

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

// Endpoint - Marmoset endpoint
var Endpoint string

// Product - Channel product from marmoset
type Product struct {
	Id                 *ProductId       `json:"product_id"`
	Active             bool             `json:"active"`
	Enabled            bool             `json:"enabled"`
	ManuallySet        *bool            `json:"manually_set,omitempty"`
	ChannelCategoryIDs []int            `json:"channel_category_ids"`
	Taxonomies         map[string][]int `json:"taxonomies"`
}

type ProductId struct {
	ShopCode    string `json:"shop_code"`
	ShopId      int    `json:"shop_id"`
	PublisherId int    `json:"publisher_id"`
}

type ProductsQuery struct {
	ShopId            int
	PublisherId       int
	Skip              *int
	Limit             *int
	Active            *bool
	Enabled           *bool
	ManuallySet       *bool
	LastUpdatedBefore *time.Time
	ShopCodes         *[]string
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
	if productsQuery.LastUpdatedBefore != nil {
		query.Add("last_updated_before", (*productsQuery.LastUpdatedBefore).Format("2006-01-02"))
	}

	return query.Encode()
}

type Finder func(*ProductsQuery) ([]*Product, error)

// Use for tests to stub calls to API
func DummyFinder(productsQuery *ProductsQuery) ([]*Product, error) {
	return []*Product{}, nil
}

type Products []*Product

type ProductsByShopCode struct{ Products }

func (s Products) Len() int      { return len(s) }
func (s Products) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ProductsByShopCode) Less(i, j int) bool {
	return s.Products[i].Id.ShopCode < s.Products[j].Id.ShopCode
}

func FindSorted(productsQuery *ProductsQuery) ([]*Product, error) {
	products, err := Find(productsQuery)
	sort.Sort(ProductsByShopCode{products})
	return products, err
}

func Find(productsQuery *ProductsQuery) ([]*Product, error) {
	productUrl, err := buildQueryUrl(productsQuery)
	response, err := http.Get(productUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	products := []*Product{}
	err = json.NewDecoder(response.Body).Decode(&products)
	if err != nil {
		return nil, err
	}
	return products, nil
}
func Touch(shopId, publisherId int, shopCodes []string) error {
	url, err := buildTouchUrl(shopId, publisherId)
	if err != nil {
		return err
	}
	jsonShopCodes, err := json.Marshal(shopCodes)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewReader(jsonShopCodes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(req)
	return err
}
func buildTouchUrl(shopId, publisherId int) (string, error) {
	uri, err := url.Parse(fmt.Sprintf("%s/shops/%d/publishers/%d/products/touch", Endpoint, shopId, publisherId))
	return uri.String(), err
}

func buildQueryUrl(productsQuery *ProductsQuery) (string, error) {
	productUrl, err := url.Parse(fmt.Sprintf("%s/shops/%d/publishers/%d/products", Endpoint, productsQuery.ShopId, productsQuery.PublisherId))
	if err != nil {
		return "", err
	}
	productUrl.RawQuery = productsQuery.RawQuery()
	return productUrl.String(), nil
}
