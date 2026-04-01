package model

type MaintenanceDTO struct {
	ID                   string  `json:"id"`
	AssetID              string  `json:"assetId"`
	MaintenanceType      string  `json:"maintenanceType"`
	Description          string  `json:"description"`
	SentForRepairAt      string  `json:"sentForRepairAt"`
	Vendor               string  `json:"vendor"`
	Notes                string  `json:"notes"`
	MaintStatus          string  `json:"maintStatus"`
	ReturnedFromRepairAt string  `json:"returnedFromRepairAt"`
	RepairCostINR        float64 `json:"repairCostINR"`
	CreatedAt            string  `json:"createdAt"`
	UpdatedAt            string  `json:"updatedAt"`
}

type MaintenanceRequest struct {
	MaintenanceType string `json:"maintenanceType"`
	Description     string `json:"description"`
	SentForRepairAt string `json:"sentForRepairAt"`
	Vendor          string `json:"vendor"`
}

type UpdateMaintenanceRequest struct {
	Status               string  `json:"status"`
	ReturnedFromRepairAt string  `json:"returnedFromRepairAt"`
	RepairCostINR        float64 `json:"repairCostINR"`
	Vendor               string  `json:"vendor"`
	Notes                string  `json:"notes"`
}
