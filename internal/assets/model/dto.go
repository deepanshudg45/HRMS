package model

// request/response
type AssetDTO struct {
    AssetCode string `json:"assetCode"`
    AssetName string `json:"assetName"`
    AssetType string `json:"assetType"`
    Status    string `json:"status"`
}
