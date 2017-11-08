package channel_product

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type MappingParams struct {
	ChannelCategoryIds []int    `json:"channel_category_ids"`
	ShopCodes          []string `json:"shop_codes"`
	UserID             int      `json:"user_id"`
}

func Map(shopId, publisherId int, shopCodes []string, taxonomyTypeID, taxonomyID int) error {
	url, err := buildMapUrl(shopId, publisherId)
	if err != nil {
		return err
	}
	mappingParams := MappingParams{ShopCodes: shopCodes, ChannelCategoryIds: []int{taxonomyID}, UserID: 1}
	jsonData, err := json.Marshal(mappingParams)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	return err
}
func buildMapUrl(shopId, publisherId int) (string, error) {
	uri, err := url.Parse(fmt.Sprintf("%s/shops/%d/publishers/%d/products/map_from_tip", Endpoint, shopId, publisherId))
	return uri.String(), err
}
