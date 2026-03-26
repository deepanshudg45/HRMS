package repository

import (
	"context"
	"fmt"

	"WITS/internal/assets/model"
)

// GetAssets retrieves assets with dynamic filtering and pagination
func (r *Repository) GetAssets(ctx context.Context, filters *model.AssetFilter) ([]model.AssetListDTO, int, error) {
	query := `SELECT id, asset_code, name, serial_no, asset_type, category, status 
	          FROM asset_inventory WHERE is_deleted = FALSE`

	args := []interface{}{}
	paramCount := 1

	if filters.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", paramCount)
		args = append(args, filters.Status)
		paramCount++
	}

	if filters.Type != "" {
		query += fmt.Sprintf(" AND asset_type = $%d", paramCount)
		args = append(args, filters.Type)
		paramCount++
	}

	if filters.Category != "" {
		query += fmt.Sprintf(" AND category = $%d", paramCount)
		args = append(args, filters.Category)
		paramCount++
	}

	if filters.Search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR serial_no ILIKE $%d OR asset_code ILIKE $%d)", paramCount, paramCount, paramCount)
		args = append(args, "%"+filters.Search+"%")
		paramCount++
	}

	// Get total count using a subquery
	countSQL := `SELECT COUNT(*) FROM (` + query + `) AS filtered`
	var total int
	err := r.DB.QueryRow(ctx, countSQL, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
	args = append(args, filters.Limit, filters.Offset)

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	assets := []model.AssetListDTO{}
	for rows.Next() {
		var a model.AssetListDTO
		if err := rows.Scan(&a.ID, &a.Code, &a.Name, &a.SerialNo, &a.Type, &a.Category, &a.Status); err != nil {
			return nil, 0, err
		}
		assets = append(assets, a)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return assets, total, nil
}