package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"WITS/internal/assets/model"
	"WITS/package/utils"

	"github.com/jackc/pgx/v5"
)

var (
	ErrMissingRequiredFields     = errors.New("assetType, assetName and assetCategory are required")
	ErrInvalidAssetType          = errors.New("invalid asset type")
	ErrAssetInUse                = errors.New("asset is in use")
	ErrAssetNotFound             = errors.New("asset not found")
	ErrInvalidAssetStatus        = errors.New("invalid asset status")
	ErrInvalidStatusChange       = errors.New("only RETIRED and LOST transitions are allowed")
	ErrEmployeeIDRequired        = errors.New("employee id is required")
	ErrActiveAssignmentNotFound  = errors.New("active assignment not found")
	ErrConditionAtReturnRequired = errors.New("conditionAtReturn is required")
	ErrForbiddenAssignmentAccess = errors.New("you do not own this assignment")
)

type AssetRepository interface {
	NextAssetSeq(ctx context.Context) (int64, error)
	CreateAsset(ctx context.Context, asset model.Asset) error
	GetAssets(ctx context.Context, filters *model.AssetFilter) ([]model.AssetListDTO, int, error)
	GetMaintenanceRecordsByAssetID(ctx context.Context, assetID string) ([]model.MaintenanceDTO, error)
	GetAssetStatusByID(ctx context.Context, assetID string) (string, error)
	SoftDeleteAsset(ctx context.Context, assetID string) error
	UpdateAssetStatus(ctx context.Context, assetID string, status string) (*model.AssetDTO, error)
	GetActiveAssetsByEmployeeID(ctx context.Context, employeeID string) ([]model.MyAssetDTO, error)
	GetActiveAssignmentIDByAssetID(ctx context.Context, assetID string) (string, error)
	DeactivateAssignment(ctx context.Context, assignmentID string) error
	GetAssignmentOwnerByID(ctx context.Context, assignmentID string) (string, error)
	GetAssignmentsByAssetID(ctx context.Context, assetID string) ([]model.AssignmentHistoryDTO, error)
	GetAssetsByEmployeeID(ctx context.Context, employeeID string) ([]model.MyAssetDTO, error)
	GenerateReport(ctx context.Context, filters *model.AssetFilter) ([]model.AssetReportDTO, error)
}

type AssetService struct {
	repo AssetRepository
}

func NewAssetService(repo AssetRepository) *AssetService {
	return &AssetService{repo: repo}
}

func (s *AssetService) CreateAsset(ctx context.Context, req model.CreateAssetRequest) (*model.AssetDTO, error) {
	req.AssetName = strings.TrimSpace(req.AssetName)
	req.AssetCategory = strings.TrimSpace(req.AssetCategory)
	req.Brand = strings.TrimSpace(req.Brand)
	req.Model = strings.TrimSpace(req.Model)
	req.SerialNo = strings.TrimSpace(req.SerialNo)
	req.PurchaseDate = strings.TrimSpace(req.PurchaseDate)
	req.Vendor = strings.TrimSpace(req.Vendor)
	req.WarrantyExpiry = strings.TrimSpace(req.WarrantyExpiry)
	req.Location = strings.TrimSpace(req.Location)
	req.Notes = strings.TrimSpace(req.Notes)

	if req.AssetType == "" || req.AssetName == "" || req.AssetCategory == "" {
		return nil, ErrMissingRequiredFields
	}

	validTypes := map[string]bool{
		"LAPTOP":      true,
		"MOBILE":      true,
		"ACCESS_CARD": true,
	}

	assetType := strings.ToUpper(strings.TrimSpace(req.AssetType))
	if !validTypes[assetType] {
		return nil, ErrInvalidAssetType
	}
	req.AssetType = assetType

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

	err = s.repo.CreateAsset(ctx, asset)
	if err != nil {
		return nil, err
	}

	return &model.AssetDTO{
		AssetCode: assetCode,
		AssetName: req.AssetName,
		AssetType: req.AssetType,
		Status:    "CREATED",
	}, nil
}

// GetAssets retrieves assets with filters and pagination
func (s *AssetService) GetAssets(ctx context.Context, filters *model.AssetFilter) ([]model.AssetListDTO, int, error) {
	return s.repo.GetAssets(ctx, filters)
}

func (s *AssetService) GetMaintenanceRecordsByAssetID(ctx context.Context, assetID string) ([]model.MaintenanceDTO, error) {
	return s.repo.GetMaintenanceRecordsByAssetID(ctx, assetID)
}

func (s *AssetService) DeleteAsset(ctx context.Context, assetID string) error {
	status, err := s.repo.GetAssetStatusByID(ctx, assetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrAssetNotFound
		}
		return err
	}

	if status != "AVAILABLE" && status != "RETIRED" {
		return ErrAssetInUse
	}

	if err := s.repo.SoftDeleteAsset(ctx, assetID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrAssetNotFound
		}
		return err
	}

	return nil
}

func (s *AssetService) UpdateAssetStatus(ctx context.Context, assetID string, req model.StatusRequest) (*model.AssetDTO, error) {
	nextStatus := strings.ToUpper(strings.TrimSpace(req.Status))
	if nextStatus == "" {
		return nil, ErrInvalidAssetStatus
	}

	if nextStatus != "RETIRED" && nextStatus != "LOST" {
		return nil, ErrInvalidStatusChange
	}

	currentStatus, err := s.repo.GetAssetStatusByID(ctx, assetID)
	if err != nil {
		return nil, err
	}

	if currentStatus == "ASSIGNED" {
		return nil, ErrAssetInUse
	}

	return s.repo.UpdateAssetStatus(ctx, assetID, nextStatus)
}

func (s *AssetService) GetActiveAssets(ctx context.Context, employeeID string) ([]model.MyAssetDTO, error) {
	employeeID = strings.TrimSpace(employeeID)
	if employeeID == "" {
		return nil, ErrEmployeeIDRequired
	}

	return s.repo.GetActiveAssetsByEmployeeID(ctx, employeeID)
}

func (s *AssetService) ReturnAsset(ctx context.Context, assetID string, req model.ReturnRequest) (*model.ReturnDTO, error) {
	req.ReturnedOn = strings.TrimSpace(req.ReturnedOn)
	req.ConditionAtReturn = strings.ToUpper(strings.TrimSpace(req.ConditionAtReturn))
	req.ReturnReason = strings.TrimSpace(req.ReturnReason)

	if req.ConditionAtReturn == "" {
		return nil, ErrConditionAtReturnRequired
	}

	assignmentID, err := s.repo.GetActiveAssignmentIDByAssetID(ctx, assetID)
	if err != nil {
		return nil, ErrActiveAssignmentNotFound
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
		return nil, ErrConditionAtReturnRequired
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
		return nil, ErrEmployeeIDRequired
	}

	ownerID, err := s.repo.GetAssignmentOwnerByID(ctx, assignmentID)
	if err != nil {
		return nil, err
	}

	if ownerID != employeeID {
		return nil, ErrForbiddenAssignmentAccess
	}

	acknowledgedAt := time.Now().Format(time.RFC3339)
	return &model.AssignmentDTO{
		ID:                    assignmentID,
		EmployeeID:            employeeID,
		AcknowledgementStatus: "ACKNOWLEDGED",
		AcknowledgedAt:        acknowledgedAt,
	}, nil
}

func (s *AssetService) GetMyAssignments(ctx context.Context, assetID string) ([]model.AssignmentHistoryDTO, error) {
	assetID = strings.TrimSpace(assetID)
	if assetID == "" {
		return nil, ErrInvalidAssetStatus
	}

	return s.repo.GetAssignmentsByAssetID(ctx, assetID)
}

func (s *AssetService) GetAssetsByEmployeeID(ctx context.Context, employeeID string) ([]model.MyAssetDTO, error) {
	employeeID = strings.TrimSpace(employeeID)
	if employeeID == "" {
		return nil, ErrEmployeeIDRequired
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
