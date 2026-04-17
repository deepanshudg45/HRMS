package service

import (
	"context"
	"errors"
	"strings"

	"WITS/internal/assets/model"
)

type AssetRepository interface {
	NextAssetSeq(ctx context.Context) (int64, error)
	CreateAsset(ctx context.Context, asset model.Asset) error
	GetAssets(ctx context.Context, filters *model.AssetFilter) ([]model.AssetListDTO, int, error)
	GetAssetByID(ctx context.Context, assetID string) (*model.AssetDetailDTO, error)
	IsHREmployee(ctx context.Context, employeeID string) (bool, error)
	UpdateAsset(ctx context.Context, assetID string, req model.UpdateAssetRequest) error
	GetMaintenanceRecordsByAssetID(ctx context.Context, assetID string) ([]model.MaintenanceDTO, error)
	CreateMaintenanceRecord(ctx context.Context, assetID string, req model.MaintenanceRequest) (*model.MaintenanceDTO, error)
	UpdateMaintenanceRecord(ctx context.Context, assetID string, maintenanceID string, req model.UpdateMaintenanceRequest) (*model.MaintenanceDTO, error)
	GetAssetStatusByID(ctx context.Context, assetID string) (string, error)
	AssignAsset(ctx context.Context, assetID string, req model.AssignRequest) (*model.AssignAssetDTO, error)
	DeleteAsset(ctx context.Context, assetID string) error
	UpdateAssetStatus(ctx context.Context, assetID string, status string) (*model.AssetDTO, error)
	GetActiveAssetsByEmployeeID(ctx context.Context, employeeID string) ([]model.MyAssetDTO, error)
	GetActiveAssignmentIDByAssetID(ctx context.Context, assetID string) (string, error)
	DeactivateAssignment(ctx context.Context, assignmentID string) error
	GetAssignmentOwnerByID(ctx context.Context, assignmentID string) (string, error)
	UpdateAssignmentAcknowledgement(ctx context.Context, assignmentID string) (*model.AssignmentDTO, error)
	GetAssignmentsByAssetID(ctx context.Context, assetID string) ([]model.AssignmentHistoryDTO, error)
	GetAssetsByEmployeeID(ctx context.Context, employeeID string) ([]model.MyAssetDTO, error)
	GenerateReport(ctx context.Context, filters *model.AssetFilter) ([]model.AssetReportDTO, error)
	GetExpiringWarrantyAssetIDs(ctx context.Context) ([]string, error)
}

type AssetService struct {
	repo       AssetRepository
	dispatcher AssetEventDispatcher
}

func NewAssetService(repo AssetRepository) *AssetService {
	return &AssetService{
		repo:       repo,
		dispatcher: noopAssetEventDispatcher{},
	}
}

func (s *AssetService) GetMaintenanceRecordsByAssetID(ctx context.Context, assetID string) ([]model.MaintenanceDTO, error) {
	return s.repo.GetMaintenanceRecordsByAssetID(ctx, assetID)
}

func (s *AssetService) CreateMaintenanceRecord(ctx context.Context, assetID string, req model.MaintenanceRequest) (*model.MaintenanceDTO, error) {
	assetID = strings.TrimSpace(assetID)
	if assetID == "" {
		return nil, errors.New("asset id is required")
	}

	req.MaintenanceType = strings.ToUpper(strings.TrimSpace(req.MaintenanceType))
	req.Description = strings.TrimSpace(req.Description)
	req.SentForRepairAt = strings.TrimSpace(req.SentForRepairAt)
	req.Vendor = strings.TrimSpace(req.Vendor)

	if req.Description == "" {
		return nil, errors.New("description is required")
	}

	validTypes := map[string]bool{
		"REPAIR":     true,
		"SERVICE":    true,
		"INSPECTION": true,
		"DISPOSAL":   true,
	}
	if !validTypes[req.MaintenanceType] {
		return nil, errors.New("invalid maintenance type")
	}

	record, err := s.repo.CreateMaintenanceRecord(ctx, assetID, req)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (s *AssetService) UpdateMaintenanceRecord(ctx context.Context, assetID string, maintenanceID string, req model.UpdateMaintenanceRequest) (*model.MaintenanceDTO, error) {
	assetID = strings.TrimSpace(assetID)
	maintenanceID = strings.TrimSpace(maintenanceID)
	if assetID == "" || maintenanceID == "" {
		return nil, errors.New("asset id and maintenance id are required")
	}

	req.Status = strings.ToUpper(strings.TrimSpace(req.Status))
	req.ReturnedFromRepairAt = strings.TrimSpace(req.ReturnedFromRepairAt)
	req.Vendor = strings.TrimSpace(req.Vendor)
	req.Notes = strings.TrimSpace(req.Notes)

	if req.Status != "COMPLETED" && req.Status != "SCRAPPED" && req.Status != "IN_PROGRESS" {
		return nil, errors.New("invalid maintenance status")
	}

	record, err := s.repo.UpdateMaintenanceRecord(ctx, assetID, maintenanceID, req)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func normalizeAssetUpsertRequest(assetType string, assetName string, brand string, modelName string, assetCategory string, serialNo string, purchaseDate string, purchaseCostINR float64, vendor string, warrantyExpiry string, location string, notes string) (*model.CreateAssetRequest, error) {
	req := &model.CreateAssetRequest{
		AssetType:       strings.ToUpper(strings.TrimSpace(assetType)),
		AssetName:       strings.TrimSpace(assetName),
		Brand:           strings.TrimSpace(brand),
		Model:           strings.TrimSpace(modelName),
		AssetCategory:   strings.TrimSpace(assetCategory),
		SerialNo:        strings.TrimSpace(serialNo),
		PurchaseDate:    strings.TrimSpace(purchaseDate),
		PurchaseCostINR: purchaseCostINR,
		Vendor:          strings.TrimSpace(vendor),
		WarrantyExpiry:  strings.TrimSpace(warrantyExpiry),
		Location:        strings.TrimSpace(location),
		Notes:           strings.TrimSpace(notes),
	}

	if req.AssetType == "" || req.AssetName == "" || req.AssetCategory == "" {
		return nil, errors.New("assetType, assetName and assetCategory are required")
	}

	validTypes := map[string]bool{
		"LAPTOP":    true,
		"MOBILE":    true,
		"DESKTOP":   true,
		"FURNITURE": true,
		"OTHER":     true,
	}

	if !validTypes[req.AssetType] {
		return nil, errors.New("invalid asset type")
	}

	return req, nil
}

type AssetEventDispatcher interface {
	Dispatch(ctx context.Context, eventName string, targetEmployee string, assetID string) error
}

const EventAssetAssigned = "EventAssetAssigned"
const EventWarrantyExpiring = "EventWarrantyExpiring"
const HRAdminID = "hrAdmin"

type noopAssetEventDispatcher struct{}

func (noopAssetEventDispatcher) Dispatch(ctx context.Context, eventName string, targetEmployee string, assetID string) error {
	return nil
}
