package service

import (
	"context"

	"WITS/internal/assets/model"
	"WITS/internal/assets/repository"
)

type Service struct {
	Repo *repository.Repository
}

// GetAssets retrieves assets with filters and pagination
func (s *Service) GetAssets(ctx context.Context, filters *model.AssetFilter) ([]model.AssetListDTO, int, error) {
	return s.Repo.GetAssets(ctx, filters)
}
