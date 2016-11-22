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

const (
	DescriptionKey      = "Product Description"
	FloorPriceKey       = "Cost Price"
	StockKey            = "Product in stock"
	CategoryPathKey     = "Category Path"
	CategoryKey         = "Category"
	SubCategoryKey      = "Sub category"
	ShopCodeKey         = "Shop Code"
	VariantIDKey        = "Variant ID"
	ProductNameKey      = "Product Name"
	PictureLinkKey      = "Picture Link"
	DeeplinkKey         = "Deeplink"
	ProductEanKey       = "Product Ean"
	ProductBrandKey     = "Product Brand"
	DeliveryPeriodKey   = "Delivery Period"
	ProductInStockKey   = "Product in stock"
	StockStatusKey      = "Stock Status"
	EnabledKey          = "Enabled"
	DisabledAtKey       = "Disabled At"
	SellingPriceExclKey = "Selling Price Ex"
	SellingPriceInclKey = "Selling Price"
)

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
func (s ShopProduct) CategoryPath() string {
	return s[CategoryPathKey]
}
func (s ShopProduct) SubCategory() string {
	return s[SubCategoryKey]
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

func (s ShopProduct) SellingPriceIncl() string {
	return s[SellingPriceInclKey]
}
func (s ShopProduct) SellingPriceExcl() string {
	return s[SellingPriceExclKey]
}
func (s ShopProduct) Enabled() string {
	return s[EnabledKey]
}
func (s ShopProduct) DisabledAt() string {
	return s[DisabledAtKey]
}

type Finder func(int) (<-chan ShopProduct, <-chan error)

func Find(shopId int) (<-chan ShopProduct, <-chan error) {
	shopProductChannel := make(chan ShopProduct)
	errorChannel := make(chan error, 5)
	resp, err := http.Get(catalogUrl(shopId))
	if err != nil {
		errorChannel <- err
		return shopProductChannel, errorChannel
	}
	fmt.Printf("-- Shop ID %d stream opened\n", shopId)
	go func() {
		defer close(shopProductChannel)
		defer close(errorChannel)
		decoder := json.NewDecoder(resp.Body)
		defer resp.Body.Close()
		defer fmt.Printf("-- Shop ID %d stream closed\n", shopId)
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
