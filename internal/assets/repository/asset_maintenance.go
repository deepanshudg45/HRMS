package repository

import (
	"context"
	"fmt"

	"WITS/internal/assets/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Row interface {
	Scan(dest ...any) error
}

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...any) interface{ Scan(dest ...any) error }
}

type Repository struct {
	DB DB
}

func NewRepository(db DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateAsset(ctx context.Context, asset model.Asset) error {
	query := `
        INSERT INTO asset_inventory 
        (
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

	_, err := r.DB.Exec(ctx, query,
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

func (r *Repository) GetMaintenanceRecordsByAssetID(ctx context.Context, assetID string) ([]model.MaintenanceDTO, error) {
	query := `
        SELECT
            id,
            asset_id,
            maintenance_type,
            description,
            sent_for_repair_at,
            vendor,
            maint_status,
            returned_from_repair_at,
            repair_cost_inr,
            created_at,
            updated_at
        FROM asset_maintenance_logs
        WHERE asset_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.DB.Query(ctx, query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := []model.MaintenanceDTO{}
	for rows.Next() {
		var record model.MaintenanceDTO

		err := rows.Scan(
			&record.ID,
			&record.AssetID,
			&record.MaintenanceType,
			&record.Description,
			&record.SentForRepairAt,
			&record.Vendor,
			&record.MaintStatus,
			&record.ReturnedFromRepairAt,
			&record.RepairCostINR,
			&record.CreatedAt,
			&record.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func (r *Repository) GetAssetStatusByID(ctx context.Context, assetID string) (string, error) {
	query := `
        SELECT status
        FROM asset_inventory
        WHERE id = $1
    `

	var status string
	err := r.DB.QueryRow(ctx, query, assetID).Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (r *Repository) SoftDeleteAsset(ctx context.Context, assetID string) error {
	query := `
        WITH deleted_maintenance AS (
            DELETE FROM asset_maintenance_logs
            WHERE asset_id = $1
        ),
        deleted_assignments AS (
            DELETE FROM asset_assignments
            WHERE asset_id = $1
        )
        DELETE FROM asset_inventory
        WHERE id = $1
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

func (r *Repository) GetActiveAssetsByEmployeeID(ctx context.Context, employeeID string) ([]model.MyAssetDTO, error) {
	query := `
        SELECT
            ai.asset_code,
            ai.asset_type,
            ai.name,
            aa.created_at
        FROM asset_assignments aa
        JOIN asset_inventory ai ON ai.id = aa.asset_id
        WHERE aa.assigned_to = $1
          AND aa.is_active = TRUE
          AND ai.is_deleted = FALSE
        ORDER BY aa.created_at DESC
    `

	rows, err := r.DB.Query(ctx, query, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assets := []model.MyAssetDTO{}
	for rows.Next() {
		var asset model.MyAssetDTO
		if err := rows.Scan(
			&asset.AssetCode,
			&asset.AssetType,
			&asset.AssetName,
			&asset.AssignedOn,
		); err != nil {
			return nil, err
		}

		assets = append(assets, asset)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assets, nil
}

func (r *Repository) GetActiveAssignmentIDByAssetID(ctx context.Context, assetID string) (string, error) {
	query := `
        SELECT id
        FROM asset_assignments
        WHERE asset_id = $1 AND is_active = TRUE
    `

	var assignmentID string
	err := r.DB.QueryRow(ctx, query, assetID).Scan(&assignmentID)
	if err != nil {
		return "", err
	}

	return assignmentID, nil
}

func (r *Repository) DeactivateAssignment(ctx context.Context, assignmentID string) error {
	query := `
        UPDATE asset_assignments
        SET is_active = FALSE
        WHERE id = $1 AND is_active = TRUE
    `

	_, err := r.DB.Exec(ctx, query, assignmentID)
	return err
}

func (r *Repository) GetAssignmentOwnerByID(ctx context.Context, assignmentID string) (string, error) {
	query := `
        SELECT assigned_to
        FROM asset_assignments
        WHERE id = $1
    `

	var employeeID string
	err := r.DB.QueryRow(ctx, query, assignmentID).Scan(&employeeID)
	if err != nil {
		return "", err
	}

	return employeeID, nil
}

func (r *Repository) GetAssignmentsByAssetID(ctx context.Context, assetID string) ([]model.AssignmentHistoryDTO, error) {
	query := `
        SELECT
            id,
            asset_id,
            assigned_to,
            is_active,
            created_at
        FROM asset_assignments
        WHERE asset_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.DB.Query(ctx, query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignments := []model.AssignmentHistoryDTO{}
	for rows.Next() {
		var assignment model.AssignmentHistoryDTO
		if err := rows.Scan(
			&assignment.ID,
			&assignment.AssetID,
			&assignment.EmployeeID,
			&assignment.IsActive,
			&assignment.AssignedOn,
		); err != nil {
			return nil, err
		}

		assignments = append(assignments, assignment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assignments, nil
}

func (r *Repository) GetAssetsByEmployeeID(ctx context.Context, employeeID string) ([]model.MyAssetDTO, error) {
	return r.GetActiveAssetsByEmployeeID(ctx, employeeID)
}

func (r *Repository) GenerateReport(ctx context.Context, filters *model.AssetFilter) ([]model.AssetReportDTO, error) {
	query := `
        SELECT
            id,
            asset_code,
            name,
            asset_type,
            category,
            status,
            serial_no,
            location,
            vendor,
            created_at
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
