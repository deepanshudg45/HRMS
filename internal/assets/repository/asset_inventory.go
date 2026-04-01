package repository

import (
	"context"
	"database/sql"
	"fmt"

	"WITS/internal/assets/model"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateAsset(ctx context.Context, asset model.Asset) error {
	query := `
        INSERT INTO asset_inventory (
            asset_code,
            asset_type,
            name,
            brand,
            model,
            category,
            serial_no,
            purchase_date,
            purchase_cost_inr,
            vendor,
            warranty_expiry,
            location,
            notes
        )
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
    `

	_, err := r.DB.Exec(
		ctx,
		query,
		asset.AssetCode,
		asset.AssetType,
		asset.AssetName,
		asset.Brand,
		asset.Model,
		asset.AssetCategory,
		asset.SerialNo,
		asset.PurchaseDate,
		asset.PurchaseCostINR,
		asset.Vendor,
		asset.WarrantyExpiry,
		asset.Location,
		asset.Notes,
	)

	return err
}

func (r *Repository) NextAssetSeq(ctx context.Context) (int64, error) {
	var seq int64

	query := `SELECT nextval('asset_code_seq')`

	err := r.DB.QueryRow(ctx, query).Scan(&seq)
	if err != nil {
		return 0, err
	}

	return seq, nil
}

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

	countSQL := `SELECT COUNT(*) FROM (` + query + `) AS filtered`
	var total int
	err := r.DB.QueryRow(ctx, countSQL, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

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

func (r *Repository) GetAssetStatusByID(ctx context.Context, assetID string) (string, error) {
	query := `
        SELECT status
        FROM asset_inventory
        WHERE id = $1 AND is_deleted = FALSE
    `

	var status string
	err := r.DB.QueryRow(ctx, query, assetID).Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (r *Repository) UpdateAsset(ctx context.Context, assetID string, req model.UpdateAssetRequest) error {
	query := `
        UPDATE asset_inventory
        SET
            brand = $2,
            model = $3,
            warranty_expiry = NULLIF($4, '')::date,
            location = $5,
            notes = $6
        WHERE id = $1 AND is_deleted = FALSE
    `

	tag, err := r.DB.Exec(ctx, query,
		assetID,
		req.Brand,
		req.Model,
		req.WarrantyExpiry,
		req.Location,
		req.Notes,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *Repository) SoftDeleteAsset(ctx context.Context, assetID string) error {
	query := `
        UPDATE asset_inventory
        SET is_deleted = TRUE
        WHERE id = $1 AND is_deleted = FALSE
    `

	tag, err := r.DB.Exec(ctx, query, assetID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *Repository) UpdateAssetStatus(ctx context.Context, assetID string, status string) (*model.AssetDTO, error) {
	query := `
        UPDATE asset_inventory
        SET status = $2
        WHERE id = $1 AND is_deleted = FALSE
        RETURNING asset_code, name, asset_type, status
    `

	var asset model.AssetDTO
	err := r.DB.QueryRow(ctx, query, assetID, status).Scan(
		&asset.AssetCode,
		&asset.AssetName,
		&asset.AssetType,
		&asset.Status,
	)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (r *Repository) GenerateReport(ctx context.Context, filters *model.AssetFilter) ([]model.AssetReportDTO, error) {
	query := `
        SELECT
            id,
            asset_code,
            name,
            asset_type,
            COALESCE(category::text, ''),
            status,
            COALESCE(serial_no, ''),
            COALESCE(location, ''),
            COALESCE(vendor, ''),
            created_at::text
        FROM asset_inventory
        WHERE is_deleted = FALSE
    `

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

	query += " ORDER BY created_at DESC"

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reports := []model.AssetReportDTO{}
	for rows.Next() {
		var report model.AssetReportDTO
		if err := rows.Scan(
			&report.ID,
			&report.AssetCode,
			&report.AssetName,
			&report.AssetType,
			&report.Category,
			&report.Status,
			&report.SerialNo,
			&report.Location,
			&report.Vendor,
			&report.CreatedAt,
		); err != nil {
			return nil, err
		}

		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *Repository) GetExpiringWarrantyAssetIDs(ctx context.Context) ([]string, error) {
	query := `
        SELECT id::text
        FROM asset_inventory
        WHERE warranty_expiry = CURRENT_DATE + INTERVAL '30 days'
          AND status NOT IN ('RETIRED', 'LOST')
          AND is_deleted = FALSE
    `

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assetIDs := []string{}
	for rows.Next() {
		var assetID string
		if err := rows.Scan(&assetID); err != nil {
			return nil, err
		}

		assetIDs = append(assetIDs, assetID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assetIDs, nil
}
