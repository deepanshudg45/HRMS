package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"WITS/internal/assets/model"
	"WITS/package/utils"
)

var (
    ErrMissingRequiredFields = errors.New("assetType, assetName and assetCategory are required")
    ErrInvalidAssetType      = errors.New("invalid asset type")
)

type AssetRepository interface {
    NextAssetSeq(ctx context.Context) (int64, error)
    CreateAsset(ctx context.Context, asset model.Asset) error
    GetAssets(ctx context.Context, filters *model.AssetFilter) ([]model.AssetListDTO, int, error)
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
