package publisher_connection

type ConnectionDetails struct {
	Enabled             bool    `json:"enabled"`
	ShopId              *string `json:"shop_id,omitempty"`
	ProductUpdateApiKey *string `json:"product_update_api_key"`
}

type Connection struct {
	ID                 int    `json:"id"`
	Imported           bool   `json:"imported"`
	Connected          bool   `json:"connected"`
	ConnectionType     string `json:"connection_type"`
	ConnectionProvider string `json:"connection_provider"`
	*ConnectionDetails `json:"connection"`
}
