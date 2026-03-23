package assets

type AssetListDTO struct {
	ID       string `json:"id"`
	Code     string `json:"asset_code"`
	Name     string `json:"name"`
	SerialNo string `json:"serial_no"`
	Type     string `json:"type"`
	Category string `json:"category"`
	Status   string `json:"status"`
}