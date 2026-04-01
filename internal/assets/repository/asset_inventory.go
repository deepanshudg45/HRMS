package repository

import (
	"context"
	"database/sql"
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

func (r *Repository) GetAssetByID(ctx context.Context, assetID string) (*model.AssetDetailDTO, error) {
	query := `
        SELECT
            ai.id,
            ai.asset_code,
            ai.name,
            ai.asset_type,
            COALESCE(ai.category::text, ''),
            ai.status,
            COALESCE(ai.serial_no, ''),
            COALESCE(ai.brand, ''),
            COALESCE(ai.model, ''),
            COALESCE(ai.purchase_date::text, ''),
            COALESCE(ai.purchase_cost_inr, 0),
            COALESCE(ai.vendor, ''),
            COALESCE(ai.warranty_expiry::text, ''),
            COALESCE(ai.location, ''),
            COALESCE(ai.notes, ''),
            a.id::text,
            a.assigned_to::text
        FROM asset_inventory ai
        LEFT JOIN asset_assignments a
            ON a.asset_id = ai.id AND a.is_active = TRUE
        WHERE ai.id = $1 AND ai.is_deleted = FALSE
    `

	var detail model.AssetDetailDTO
	var assignmentID sql.NullString
	var employeeID sql.NullString

	err := r.DB.QueryRow(ctx, query, assetID).Scan(
		&detail.ID,
		&detail.AssetCode,
		&detail.AssetName,
		&detail.AssetType,
		&detail.AssetCategory,
		&detail.Status,
		&detail.SerialNo,
		&detail.Brand,
		&detail.Model,
		&detail.PurchaseDate,
		&detail.PurchaseCostINR,
		&detail.Vendor,
		&detail.WarrantyExpiry,
		&detail.Location,
		&detail.Notes,
		&assignmentID,
		&employeeID,
	)
	if err != nil {
		return nil, err
	}

	if assignmentID.Valid {
		detail.AssignmentID = &assignmentID.String
	}

	if employeeID.Valid {
		detail.CurrentAssignee = &model.EmployeeSummaryDTO{
			EmployeeID: employeeID.String,
		}
	}

	return &detail, nil
}
