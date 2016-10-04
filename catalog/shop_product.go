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
const VariantIDKey string = "Variant ID"
const ProductNameKey string = "Product Name"
const PictureLinkKey string = "Picture Link"
const DeeplinkKey string = "Deeplink"
const ProductEanKey string = "Product Ean"
const ProductBrandKey string = "Product Brand"
const DeliveryPeriodKey string = "Delivery Period"
const ProductInStockKey string = "Product in stock"
const StockStatusKey string = "Stock Status"
const SellingPriceKey string = "Selling Price"

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

func (s ShopProduct) VariantID() string {
	return s[VariantIDKey]
}
func (s ShopProduct) ProductName() string {
	return s[ProductNameKey]
}
func (s ShopProduct) PictureLink() string {
	return s[PictureLinkKey]
}
func (s ShopProduct) Deeplink() string {
	return s[DeeplinkKey]
}
func (s ShopProduct) ProductEan() string {
	return s[ProductEanKey]
}
func (s ShopProduct) ProductBrand() string {
	return s[ProductBrandKey]
}
func (s ShopProduct) DeliveryPeriod() string {
	return s[DeliveryPeriodKey]
}
func (s ShopProduct) ProductInStock() string {
	return s[ProductInStockKey]
}
func (s ShopProduct) StockStatus() string {
	return s[StockStatusKey]
}

func (s ShopProduct) SellingPrice() string {
	return s[SellingPriceKey]
}

type Finder func(int) (<-chan ShopProduct, <-chan error)

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
	return fmt.Sprintf("%s/sorted_shops/%d.jsonl", Endpoint, shopId)
}
