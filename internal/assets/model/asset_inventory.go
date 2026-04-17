package model

import "github.com/google/uuid"

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

type AppError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

type CreateAssetRequest struct {
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

type UpdateAssetRequest struct {
	Brand          string `json:"brand"`
	Model          string `json:"model"`
	WarrantyExpiry string `json:"warrantyExpiry"`
	Location       string `json:"location"`
	Notes          string `json:"notes"`
}

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

type EmployeeSummaryDTO struct {
	EmployeeID   string `json:"employeeId"`
	Name         string `json:"name,omitempty"`
	EmployeeCode string `json:"employeeCode,omitempty"`
	AssignedOn   string `json:"assignedOn,omitempty"`
}

type AssetDetailDTO struct {
	ID              string              `json:"id"`
	AssetCode       string              `json:"assetCode"`
	AssetName       string              `json:"assetName"`
	AssetType       string              `json:"assetType"`
	AssetCategory   string              `json:"assetCategory"`
	Status          string              `json:"status"`
	SerialNo        string              `json:"serialNo"`
	Brand           string              `json:"brand"`
	Model           string              `json:"model"`
	PurchaseDate    string              `json:"purchaseDate"`
	PurchaseCostINR float64             `json:"purchaseCostINR"`
	Vendor          string              `json:"vendor"`
	WarrantyExpiry  string              `json:"warrantyExpiry"`
	Location        string              `json:"location"`
	Notes           string              `json:"notes"`
	CurrentAssignee *EmployeeSummaryDTO `json:"currentAssignee,omitempty"`
	AssignmentID    *uuid.UUID          `json:"assignmentId,omitempty"`
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

type ResponseMeta struct {
	Page       int `json:"page,omitempty"`
	Limit      int `json:"limit,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"totalPages,omitempty"`
}

type StandardResponse struct {
	Success bool          `json:"success"`
	Data    any           `json:"data,omitempty"`
	Message string        `json:"message,omitempty"`
	Meta    *ResponseMeta `json:"meta,omitempty"`
}

type AssetFilter struct {
	Status   string
	Type     string
	Category string
	Search   string
	Page     int
	Limit    int
	Offset   int
}

type StatusRequest struct {
	Status string `json:"status"`
}
