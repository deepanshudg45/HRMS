// DB model

package model

type Asset struct {
    AssetCode       string  `json:"assetCode"`
    AssetType       string  `json:"assetType"`
    AssetName       string  `json:"assetName"`
    Brand           string  `json:"brand"`
    Model           string  `json:"model"`
    AssetCategory   string  `json:"assetCategory"`
    SerialNo        string  `json:"serialNo"`
    PurchaseDate    string  `json:"purchaseDate"`
    PurchaseCostINR float64 `json:"purchaseCostINR"`
    Vendor          string  `json:"vendor"`
    WarrantyExpiry  string  `json:"warrantyExpiry"`
    Location        string  `json:"location"`
    Notes           string  `json:"notes"`
}

type CreateAssetRequest struct {
    AssetType        string  `json:"assetType" binding:"required"`
    AssetName        string  `json:"assetName" binding:"required"`
    Brand            string  `json:"brand"`
    Model            string  `json:"model"`
    AssetCategory    string  `json:"assetCategory" binding:"required"`
    SerialNo         string  `json:"serialNo"`
    PurchaseDate     string  `json:"purchaseDate"`
    PurchaseCostINR  float64 `json:"purchaseCostINR"`
    Vendor           string  `json:"vendor"`
    WarrantyExpiry   string  `json:"warrantyExpiry"`
    Location         string  `json:"location"`
    Notes            string  `json:"notes"`
}
