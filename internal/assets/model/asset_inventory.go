package model

import "github.com/google/uuid"

type AssetInventory struct {
	ID        uuid.UUID `json:"id"`
	AssetCode string    `json:"asset_code"`
	Name      string    `json:"name"`
	SerialNo  *string   `json:"serial_no"`
	AssetType string    `json:"asset_type"`
	Category  *string   `json:"category"`
	Status    string    `json:"status" default:"AVAILABLE"`
	IsDeleted bool      `json:"is_deleted" default:"false"`
	CreatedAt string    `json:"created_at"`
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
	Reason string `json:"reason"`
}
