package service

import (
	"context"
	"errors"
	"strings"

	"WITS/internal/assets/model"

	"github.com/google/uuid"
)

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
	_, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, errors.New("invalid employee id format")
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

	if err := s.dispatcher.Dispatch(ctx, EventAssetAssigned, req.EmployeeID, assetID); err != nil {
		return nil, err
	}

	return assignment, nil
}

func (s *AssetService) GetActiveAssets(ctx context.Context, employeeID string) ([]model.MyAssetDTO, error) {
	employeeID = strings.TrimSpace(employeeID)
	if employeeID == "" {
		return nil, errors.New("employee id is required")
	}

	assets, err := s.repo.GetActiveAssetsByEmployeeID(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	// Debug: return empty array if no assets found
	if len(assets) == 0 {
		return []model.MyAssetDTO{}, nil
	}

	return assets, nil
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
