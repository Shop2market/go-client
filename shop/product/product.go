package product

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Shop2market/go-client/catalog"
)

var Endpoint string

type BonoboProduct map[string]interface{}

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
func (bp BonoboProduct) toShopProducts() []catalog.ShopProduct {
	shopProducts := []catalog.ShopProduct{}
	variants := bp["variants"].([]interface{})
	for _, variantInterface := range variants {
		shopProduct := catalog.ShopProduct{}
		shopProduct[catalog.ShopCodeKey] = fetchValue(bp, "shop_code")
		shopProduct[catalog.ProductNameKey] = fetchValue(bp, "product_name")
		shopProduct[catalog.DescriptionKey] = fetchValue(bp, "description")
		shopProduct[catalog.ProductBrandKey] = fetchValue(bp, "brand")
		shopProduct[catalog.DeeplinkKey] = fetchValue(bp, "deeplink")
		shopProduct[catalog.CategoryPathKey] = fetchValue(bp, "category_path")
		variant := variantInterface.(map[string]interface{})
		shopProduct[catalog.VariantIDKey] = variant["variant_id"].(string)
		shopProduct[catalog.PictureLinkKey] = fetchValue(variant, "picture_link")
		shopProduct[catalog.ProductEanKey] = fetchValue(variant, "ean")
		shopProduct[catalog.FloorPriceKey] = fetchNumberValue(variant, "cost_price_excl")
		shopProduct[catalog.ProductInStockKey] = fetchNumberValue(variant, "product_in_stock")
		shopProduct[catalog.StockStatusKey] = fetchValue(variant, "stock_status")
		shopProduct[catalog.DeliveryPeriodKey] = fetchValue(variant, "delivery_period")
		shopProduct[catalog.SellingPriceKey] = fetchNumberValue(variant, "price_incl")

		shopProducts = append(shopProducts, shopProduct)
	}
	return shopProducts
}

func FindAll(shopId int) (<-chan catalog.ShopProduct, <-chan error) {
	shopProductChannel := make(chan catalog.ShopProduct)
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
		var shopProduct catalog.ShopProduct
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
	return fmt.Sprintf("%s/shops/%d/mongo/products?sorted=true&recent=true", Endpoint, shopId)
}
