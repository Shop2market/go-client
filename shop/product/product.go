package product

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var Endpoint string

type BonoboProduct map[string]interface{}
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
	DeliveryCostsKey    = "Delivery Costs"
	ProductInStockKey   = "Product in stock"
	StockStatusKey      = "Stock Status"
	EnabledKey          = "Enabled"
	DisabledAtKey       = "Disabled At"
	SellingPriceExclKey = "Selling Price Ex"
	SellingPriceInclKey = "Selling Price"
	VendorCodeKey       = "Vendor Code"
	MarkedPriceKey      = "Marked Price"
)

func (s ShopProduct) MarkedPrice() string {
	return s[MarkedPriceKey]
}
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
func (s ShopProduct) DeliveryCosts() string {
	return s[DeliveryCostsKey]
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
func (s ShopProduct) VendorCode() *string {
	v, exists := s[VendorCodeKey]
	if !exists {
		return nil
	}
	return &v
}
func (s ShopProduct) UserField(field int) *string {
	key := fmt.Sprintf("User%d", field)
	v, exists := s[key]
	if !exists {
		return nil
	}
	return &v
}

type Finder func(int) (<-chan ShopProduct, <-chan error)

func fetchValue(hash map[string]interface{}, key string) string {
	value, ok := hash[key].(string)
	if ok {
		return value
	}
	return ""
}
func fetchNumberValue(hash map[string]interface{}, key string) string {
	value, ok := hash[key].(float64)
	if ok {
		return fmt.Sprintf("%.0f", value)
	}
	return ""
}
func fetchBool(hash map[string]interface{}, key string) string {
	value, ok := hash[key].(bool)
	if ok && value {
		return "true"
	}
	return "false"
}
func (bp BonoboProduct) toShopProducts() []ShopProduct {
	shopProducts := []ShopProduct{}
	variants := bp["variants"].([]interface{})
	key := ""
	for _, variantInterface := range variants {
		shopProduct := ShopProduct{}
		shopProduct[ShopCodeKey] = fetchValue(bp, "shop_code")
		shopProduct[ProductNameKey] = fetchValue(bp, "product_name")
		shopProduct[DescriptionKey] = fetchValue(bp, "description")
		shopProduct[ProductBrandKey] = fetchValue(bp, "brand")
		shopProduct[DeeplinkKey] = fetchValue(bp, "deeplink")
		shopProduct[CategoryPathKey] = fetchValue(bp, "category_path")
		variant := variantInterface.(map[string]interface{})
		shopProduct[VariantIDKey] = variant["variant_id"].(string)
		shopProduct[PictureLinkKey] = fetchValue(variant, "picture_link")
		shopProduct[ProductEanKey] = fetchValue(variant, "ean")
		shopProduct[FloorPriceKey] = fetchNumberValue(variant, "cost_price_excl")
		shopProduct[ProductInStockKey] = fetchNumberValue(variant, "product_in_stock")
		shopProduct[StockStatusKey] = fetchValue(variant, "stock_status")
		shopProduct[DeliveryPeriodKey] = fetchValue(variant, "delivery_period")
		shopProduct[DeliveryCostsKey] = fetchNumberValue(variant, "delivery_cost")
		shopProduct[SellingPriceInclKey] = fetchNumberValue(variant, "price_incl")
		shopProduct[SellingPriceExclKey] = fetchNumberValue(variant, "price_excl")
		shopProduct[MarkedPriceKey] = fetchNumberValue(variant, "marked_price")
		shopProduct[EnabledKey] = fetchBool(variant, "enabled")
		shopProduct[DisabledAtKey] = fetchValue(variant, "disabled_at")
		shopProduct[VendorCodeKey] = fetchValue(variant, "vendor_code")
		for i := 1; i <= 50; i++ {
			key = fmt.Sprintf("User%d", i)
			shopProduct[key] = fetchValue(variant, strings.ToLower(key))
		}
		shopProducts = append(shopProducts, shopProduct)
	}
	return shopProducts
}

func FindAll(shopId int) (<-chan ShopProduct, <-chan error) {
	shopProductChannel := make(chan ShopProduct)
	errorChannel := make(chan error, 5)
	request, err := http.NewRequest("GET", catalogUrl(shopId), nil)
	if err != nil {
		errorChannel <- err
		return shopProductChannel, errorChannel
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		errorChannel <- err
		return shopProductChannel, errorChannel
	}
	go func() {
		defer close(shopProductChannel)
		defer close(errorChannel)
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		var product BonoboProduct
		var shopProduct ShopProduct
		for {
			err := decoder.Decode(&product)
			if err == io.EOF {
				return
			}
			if err != nil {
				errorChannel <- err
				return
			}
			for _, shopProduct = range product.toShopProducts() {
				shopProductChannel <- shopProduct
			}

		}
	}()
	return shopProductChannel, errorChannel
}

func catalogUrl(shopId int) string {
	return fmt.Sprintf("%s/shops/%d/mongo/products?sorted=true", Endpoint, shopId)
}
