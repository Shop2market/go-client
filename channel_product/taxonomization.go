package channel_product

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type TaxonomyTypeParams struct {
	TypeID string `json:"type"`
	ID     int    `json:"id"`
}
type TaxonomyParams struct {
	Taxonomies TaxonomyTypeParams `json:"taxonomies"`
	ShopCodes  []string           `json:"shop_codes"`
	UserID     int                `json:"user_id"`
}

func Taxonomize(shopId, publisherId int, shopCodes []string, taxonomyTypeID, taxonomyID int) (err error) {
	url, err := buildMapTaxonomyUrl(shopId, publisherId)
	if err != nil {
		return err
	}
	mappingParams := TaxonomyParams{ShopCodes: shopCodes, Taxonomies: TaxonomyTypeParams{TypeID: strconv.Itoa(taxonomyTypeID), ID: taxonomyID}, UserID: 1}
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
	if err == nil {
		res.Body.Close()
	}
	return err
}
func buildMapTaxonomyUrl(shopId, publisherId int) (string, error) {
	uri, err := url.Parse(fmt.Sprintf("%s/shops/%d/publishers/%d/products/taxonomy_from_tip", Endpoint, shopId, publisherId))
	return uri.String(), err
}
