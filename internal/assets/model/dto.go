package model

type AssetDTO struct {
    AssetCode string `json:"assetCode"`
    AssetName string `json:"assetName"`
    AssetType string `json:"assetType"`
    Status    string `json:"status"`
}

type AssetListDTO struct {
	ID       string `json:"id"`
	Code     string `json:"asset_code"`
	Name     string `json:"name"`
	SerialNo string `json:"serial_no"`
	Type     string `json:"type"`
	Category string `json:"category"`
	Status   string `json:"status"`
}

type MyAssetDTO struct {
	AssetCode             string `json:"assetCode"`
	AssetType             string `json:"assetType"`
	AssetName             string `json:"assetName"`
	AssignedOn            string `json:"assignedOn"`
	AcknowledgementStatus string `json:"acknowledgementStatus,omitempty"`
	ConditionAtAssignment string `json:"conditionAtAssignment,omitempty"`
}

type AssetReportDTO struct {
	ID        string `json:"id"`
	AssetCode string `json:"assetCode"`
	AssetName string `json:"assetName"`
	AssetType string `json:"assetType"`
	Category  string `json:"category"`
	Status    string `json:"status"`
	SerialNo  string `json:"serialNo"`
	Location  string `json:"location"`
	Vendor    string `json:"vendor"`
	CreatedAt string `json:"createdAt"`
}

type RowError struct {
	Row    int    `json:"row"`
	Reason string `json:"reason"`
}

type ImportResultDTO struct {
	Imported int        `json:"imported"`
	Skipped  int        `json:"skipped"`
	Errors   []RowError `json:"errors"`
}
