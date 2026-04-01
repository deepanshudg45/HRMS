package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"WITS/internal/assets/model"
	"WITS/package/utils"
)

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

func (s *AssetService) GetAssets(ctx context.Context, filters *model.AssetFilter) ([]model.AssetListDTO, int, error) {
	return s.repo.GetAssets(ctx, filters)
}

func (s *AssetService) GetAssetByID(ctx context.Context, assetID string) (*model.AssetDetailDTO, error) {
	assetID = strings.TrimSpace(assetID)
	if assetID == "" {
		return nil, errors.New("asset id is required")
	}

	return s.repo.GetAssetByID(ctx, assetID)
}

func (s *AssetService) UpdateAsset(ctx context.Context, assetID string, req model.UpdateAssetRequest) (*model.AssetDTO, error) {
	assetID = strings.TrimSpace(assetID)
	if assetID == "" {
		return nil, errors.New("asset id is required")
	}

	status, err := s.repo.GetAssetStatusByID(ctx, assetID)
	if err != nil {
		return nil, err
	}

	if status == "RETIRED" {
		return nil, errors.New("retired asset cannot be updated")
	}

	req.Brand = strings.TrimSpace(req.Brand)
	req.Model = strings.TrimSpace(req.Model)
	req.WarrantyExpiry = strings.TrimSpace(req.WarrantyExpiry)
	req.Location = strings.TrimSpace(req.Location)
	req.Notes = strings.TrimSpace(req.Notes)

	if err := s.repo.UpdateAsset(ctx, assetID, req); err != nil {
		return nil, err
	}

	detail, err := s.repo.GetAssetByID(ctx, assetID)
	if err != nil {
		return nil, err
	}

	return &model.AssetDTO{
		AssetCode: detail.AssetCode,
		AssetName: detail.AssetName,
		AssetType: detail.AssetType,
		Status:    detail.Status,
	}, nil
}

func (s *AssetService) DeleteAsset(ctx context.Context, assetID string) error {
	status, err := s.repo.GetAssetStatusByID(ctx, assetID)
	if err != nil {
		return err
	}

	if status != "AVAILABLE" && status != "RETIRED" {
		return errors.New("asset is in use")
	}

	return s.repo.SoftDeleteAsset(ctx, assetID)
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
