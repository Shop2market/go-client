package catalog

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var Endpoint string

type ShopProduct map[string]string

const DescriptionKey string = "Product Description"
const FloorPriceKey string = "Cost Price"
const StockKey string = "Product in stock"
const CategoryKey string = "Category"
const ShopCodeKey string = "Shop Code"

func (s ShopProduct) Description() string {
	return s[DescriptionKey]
}
func (s ShopProduct) FloorPrice() string {
	return s[FloorPriceKey]
}
func (s ShopProduct) Stock() (int, error) {
	return strconv.Atoi(s[StockKey])
}
func (s ShopProduct) Category() string {
	return s[CategoryKey]
}
func (s ShopProduct) ShopCode() string {
	return s[ShopCodeKey]
}

func Find(shopId int) (<-chan ShopProduct, <-chan error) {
	shopProductChannel := make(chan ShopProduct)
	errorChannel := make(chan error, 5)
	resp, err := http.Get(catalogUrl(shopId))
	if err != nil {
		errorChannel <- err
	}
	go func() {
		defer close(shopProductChannel)
		defer close(errorChannel)
		decoder := json.NewDecoder(resp.Body)
		defer resp.Body.Close()
		for {
			var shopProduct ShopProduct
			err := decoder.Decode(&shopProduct)
			if err == io.EOF {
				return
			}
			if err != nil {
				errorChannel <- err
				return
			}
			shopProductChannel <- shopProduct
		}
	}()
	return shopProductChannel, errorChannel
}

func catalogUrl(shopId int) string {
	return fmt.Sprintf("%s/shops/%d.jsonl", Endpoint, shopId)
}
