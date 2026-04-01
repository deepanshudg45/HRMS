package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"WITS/internal/assets/model"
	"WITS/package/utils"
)

type AssetRepository interface {
	NextAssetSeq(ctx context.Context) (int64, error)
	CreateAsset(ctx context.Context, asset model.Asset) error
	GetAssets(ctx context.Context, filters *model.AssetFilter) ([]model.AssetListDTO, int, error)
	GetAssetByID(ctx context.Context, assetID string) (*model.AssetDetailDTO, error)
	UpdateAsset(ctx context.Context, assetID string, asset model.Asset) error
	GetMaintenanceRecordsByAssetID(ctx context.Context, assetID string) ([]model.MaintenanceDTO, error)
	CreateMaintenanceRecord(ctx context.Context, assetID string, req model.MaintenanceRequest) (*model.MaintenanceDTO, error)
	UpdateMaintenanceRecord(ctx context.Context, assetID string, maintenanceID string, req model.UpdateMaintenanceRequest) (*model.MaintenanceDTO, error)
	GetAssetStatusByID(ctx context.Context, assetID string) (string, error)
	AssignAsset(ctx context.Context, assetID string, req model.AssignRequest) (*model.AssignAssetDTO, error)
	SoftDeleteAsset(ctx context.Context, assetID string) error
	UpdateAssetStatus(ctx context.Context, assetID string, status string) (*model.AssetDTO, error)
	GetActiveAssetsByEmployeeID(ctx context.Context, employeeID string) ([]model.MyAssetDTO, error)
	GetActiveAssignmentIDByAssetID(ctx context.Context, assetID string) (string, error)
	DeactivateAssignment(ctx context.Context, assignmentID string) error
	GetAssignmentOwnerByID(ctx context.Context, assignmentID string) (string, error)
	UpdateAssignmentAcknowledgement(ctx context.Context, assignmentID string) (*model.AssignmentDTO, error)
	GetAssignmentsByAssetID(ctx context.Context, assetID string) ([]model.AssignmentHistoryDTO, error)
	GetAssetsByEmployeeID(ctx context.Context, employeeID string) ([]model.MyAssetDTO, error)
	GenerateReport(ctx context.Context, filters *model.AssetFilter) ([]model.AssetReportDTO, error)
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

func (s *AssetService) CreateAsset(ctx context.Context, req model.CreateAssetRequest) (*model.AssetDTO, error) {
	req.AssetType = strings.ToUpper(strings.TrimSpace(req.AssetType))
	req.AssetName = strings.TrimSpace(req.AssetName)
	req.Brand = strings.TrimSpace(req.Brand)
	req.Model = strings.TrimSpace(req.Model)
	req.AssetCategory = strings.TrimSpace(req.AssetCategory)
	req.SerialNo = strings.TrimSpace(req.SerialNo)
	req.PurchaseDate = strings.TrimSpace(req.PurchaseDate)
	req.Vendor = strings.TrimSpace(req.Vendor)
	req.WarrantyExpiry = strings.TrimSpace(req.WarrantyExpiry)
	req.Location = strings.TrimSpace(req.Location)
	req.Notes = strings.TrimSpace(req.Notes)

	if req.AssetType == "" || req.AssetName == "" || req.AssetCategory == "" {
		return nil, errors.New("assetType, assetName and assetCategory are required")
	}

	seq, err := s.repo.NextAssetSeq(ctx)
	if err != nil {
		return nil, err
	}

	year := time.Now().Year()
	assetCode := utils.AssetCodeGen(year, int(seq))

	asset := model.Asset{
		AssetCode:       assetCode,
		AssetType:       req.AssetType,
		AssetName:       req.AssetName,
		Brand:           req.Brand,
		Model:           req.Model,
		AssetCategory:   req.AssetCategory,
		SerialNo:        req.SerialNo,
		PurchaseDate:    req.PurchaseDate,
		PurchaseCostINR: req.PurchaseCostINR,
		Vendor:          req.Vendor,
		WarrantyExpiry:  req.WarrantyExpiry,
		Location:        req.Location,
		Notes:           req.Notes,
	}

	if err := s.repo.CreateAsset(ctx, asset); err != nil {
		return nil, err
	}

	return &model.AssetDTO{
		AssetCode: assetCode,
		AssetName: req.AssetName,
		AssetType: req.AssetType,
		Status:    "CREATED",
	}, nil
}

func (s *AssetService) UpdateAsset(ctx context.Context, assetID string, req model.UpdateAssetRequest) (*model.AssetDetailDTO, error) {
	assetID = strings.TrimSpace(assetID)
	if assetID == "" {
		return nil, errors.New("asset id is required")
	}

	normalized, err := normalizeAssetUpsertRequest(
		req.AssetType,
		req.AssetName,
		req.Brand,
		req.Model,
		req.AssetCategory,
		req.SerialNo,
		req.PurchaseDate,
		req.PurchaseCostINR,
		req.Vendor,
		req.WarrantyExpiry,
		req.Location,
		req.Notes,
	)
	if err != nil {
		return nil, err
	}

	asset := model.Asset{
		AssetType:       normalized.AssetType,
		AssetName:       normalized.AssetName,
		Brand:           normalized.Brand,
		Model:           normalized.Model,
		AssetCategory:   normalized.AssetCategory,
		SerialNo:        normalized.SerialNo,
		PurchaseDate:    normalized.PurchaseDate,
		PurchaseCostINR: normalized.PurchaseCostINR,
		Vendor:          normalized.Vendor,
		WarrantyExpiry:  normalized.WarrantyExpiry,
		Location:        normalized.Location,
		Notes:           normalized.Notes,
	}

	if err := s.repo.UpdateAsset(ctx, assetID, asset); err != nil {
		return nil, err
	}

	return s.repo.GetAssetByID(ctx, assetID)
}

func (s *AssetService) AssignAsset(ctx context.Context, assetID string, req model.AssignRequest) (*model.AssignAssetDTO, error) {
	assetID = strings.TrimSpace(assetID)
	if assetID == "" {
		return nil, errors.New("asset id is required")
	}

	req.EmployeeID = strings.TrimSpace(req.EmployeeID)
	req.AssignedOn = strings.TrimSpace(req.AssignedOn)
	req.ConditionAtAssignment = strings.ToUpper(strings.TrimSpace(req.ConditionAtAssignment))
	req.Notes = strings.TrimSpace(req.Notes)

	if req.EmployeeID == "" {
		return nil, errors.New("employee id is required")
	}
	if req.AssignedOn == "" {
		return nil, errors.New("assignedOn is required")
	}
	if req.ConditionAtAssignment == "" {
		return nil, errors.New("conditionAtAssignment is required")
	}

	assignment, err := s.repo.AssignAsset(ctx, assetID, req)
	if err != nil {
		return nil, err
	}

	if err := s.dispatcher.Dispatch(ctx, EventAssetAssigned, req.EmployeeID); err != nil {
		return nil, err
	}

	return assignment, nil
}

func (s *AssetService) GetAssets(ctx context.Context, filters *model.AssetFilter) ([]model.AssetListDTO, int, error) {
	return s.repo.GetAssets(ctx, filters)
}

func (s *AssetService) GetAssetByID(ctx context.Context, assetID string) (*model.AssetDetailDTO, error) {
	assetID = strings.TrimSpace(assetID)
	if assetID == "" {
		return nil, errors.New("asset id is required")
	}

	asset, err := s.repo.GetAssetByID(ctx, assetID)
	if err != nil {
		return nil, err
	}

	return asset, nil
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

func (s *AssetService) DeleteAsset(ctx context.Context, assetID string) error {
	status, err := s.repo.GetAssetStatusByID(ctx, assetID)
	if err != nil {
		return err
	}

	if status != "AVAILABLE" && status != "RETIRED" {
		return errors.New("asset is in use")
	}

	if err := s.repo.SoftDeleteAsset(ctx, assetID); err != nil {
		return err
	}

	return nil
}

func (s *AssetService) UpdateAssetStatus(ctx context.Context, assetID string, req model.StatusRequest) (*model.AssetDTO, error) {
	nextStatus := strings.ToUpper(strings.TrimSpace(req.Status))
	if nextStatus == "" {
		return nil, errors.New("invalid asset status")
	}

	if nextStatus != "RETIRED" && nextStatus != "LOST" {
		return nil, errors.New("only RETIRED and LOST transitions are allowed")
	}

	currentStatus, err := s.repo.GetAssetStatusByID(ctx, assetID)
	if err != nil {
		return nil, err
	}

	if currentStatus == "ASSIGNED" {
		return nil, errors.New("asset is in use")
	}

	return s.repo.UpdateAssetStatus(ctx, assetID, nextStatus)
}

func (s *AssetService) GetActiveAssets(ctx context.Context, employeeID string) ([]model.MyAssetDTO, error) {
	employeeID = strings.TrimSpace(employeeID)
	if employeeID == "" {
		return nil, errors.New("employee id is required")
	}

	return s.repo.GetActiveAssetsByEmployeeID(ctx, employeeID)
}

func (s *AssetService) ReturnAsset(ctx context.Context, assetID string, req model.ReturnRequest) (*model.ReturnDTO, error) {
	req.ReturnedOn = strings.TrimSpace(req.ReturnedOn)
	req.ConditionAtReturn = strings.ToUpper(strings.TrimSpace(req.ConditionAtReturn))
	req.ReturnReason = strings.TrimSpace(req.ReturnReason)

	if req.ConditionAtReturn == "" {
		return nil, errors.New("conditionAtReturn is required")
	}

	assignmentID, err := s.repo.GetActiveAssignmentIDByAssetID(ctx, assetID)
	if err != nil {
		return nil, errors.New("active assignment not found")
	}

	if err := s.repo.DeactivateAssignment(ctx, assignmentID); err != nil {
		return nil, err
	}

	nextStatus := "AVAILABLE"
	switch req.ConditionAtReturn {
	case "GOOD":
		nextStatus = "AVAILABLE"
	case "DAMAGED", "UNDER_REPAIR":
		nextStatus = "UNDER_REPAIR"
	case "LOST":
		nextStatus = "LOST"
	default:
		return nil, errors.New("conditionAtReturn is required")
	}

	if _, err := s.repo.UpdateAssetStatus(ctx, assetID, nextStatus); err != nil {
		return nil, err
	}

	return &model.ReturnDTO{
		AssetID:           assetID,
		AssignmentID:      assignmentID,
		ConditionAtReturn: req.ConditionAtReturn,
		ReturnReason:      req.ReturnReason,
		ReturnedOn:        req.ReturnedOn,
		AssetStatus:       nextStatus,
	}, nil
}

func (s *AssetService) UpdateAssignment(ctx context.Context, assignmentID string, employeeID string) (*model.AssignmentDTO, error) {
	assignmentID = strings.TrimSpace(assignmentID)
	employeeID = strings.TrimSpace(employeeID)

	if employeeID == "" {
		return nil, errors.New("employee id is required")
	}

	ownerID, err := s.repo.GetAssignmentOwnerByID(ctx, assignmentID)
	if err != nil {
		return nil, err
	}

	if ownerID != employeeID {
		return nil, errors.New("you do not own this assignment")
	}

	return s.repo.UpdateAssignmentAcknowledgement(ctx, assignmentID)
}

func (s *AssetService) GetMyAssignments(ctx context.Context, assetID string) ([]model.AssignmentHistoryDTO, error) {
	assetID = strings.TrimSpace(assetID)
	if assetID == "" {
		return nil, errors.New("asset id is required")
	}

	return s.repo.GetAssignmentsByAssetID(ctx, assetID)
}

func (s *AssetService) GetAssetsByEmployeeID(ctx context.Context, employeeID string) ([]model.MyAssetDTO, error) {
	employeeID = strings.TrimSpace(employeeID)
	if employeeID == "" {
		return nil, errors.New("employee id is required")
	}

	return s.repo.GetAssetsByEmployeeID(ctx, employeeID)
}

func (s *AssetService) GenerateReport(ctx context.Context, filters *model.AssetFilter) ([]model.AssetReportDTO, error) {
	return s.repo.GenerateReport(ctx, filters)
}

func (s *AssetService) ImportAssets(ctx context.Context, requests []model.CreateAssetRequest) (*model.ImportResultDTO, error) {
	result := &model.ImportResultDTO{
		Errors: []model.RowError{},
	}

	for idx, req := range requests {
		if _, err := s.CreateAsset(ctx, req); err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, model.RowError{
				Row:    idx + 2,
				Reason: err.Error(),
			})
			continue
		}

		result.Imported++
	}

	return result, nil
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
	Dispatch(ctx context.Context, eventName string, targetEmployee string) error
}

const EventAssetAssigned = "EventAssetAssigned"

type noopAssetEventDispatcher struct{}

func (noopAssetEventDispatcher) Dispatch(ctx context.Context, eventName string, targetEmployee string) error {
	return nil
}
