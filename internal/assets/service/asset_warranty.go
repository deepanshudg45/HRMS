package service

import "context"

func (s *AssetService) CheckWarrantyExpiry(ctx context.Context) error {
	assetIDs, err := s.repo.GetExpiringWarrantyAssetIDs(ctx)
	if err != nil {
		return err
	}

	for _, assetID := range assetIDs {
		if err := s.dispatcher.Dispatch(ctx, EventWarrantyExpiring, HRAdminID, assetID); err != nil {
			return err
		}
	}

	return nil
}
