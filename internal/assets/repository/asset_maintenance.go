package repository

import (
	"context"
	"errors"
	"fmt"

	"WITS/internal/assets/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrAssetUnavailable = errors.New("asset is not available for assignment")

type Row interface {
	Scan(dest ...any) error
}

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...any) interface{ Scan(dest ...any) error }
	Begin(ctx context.Context) (pgx.Tx, error)
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
            COALESCE(description, ''),
            COALESCE(sent_for_repair_at::text, ''),
            COALESCE(vendor, ''),
            COALESCE(notes, ''),
            maint_status,
            COALESCE(returned_from_repair_at::text, ''),
            COALESCE(repair_cost_inr, 0),
            created_at::text,
            updated_at::text
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
			&record.Notes,
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

func (r *Repository) CreateMaintenanceRecord(ctx context.Context, assetID string, req model.MaintenanceRequest) (*model.MaintenanceDTO, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	statusQuery := `
        SELECT status
        FROM asset_inventory
        WHERE id = $1 AND is_deleted = FALSE
        FOR UPDATE
    `

	var currentStatus string
	if err := tx.QueryRow(ctx, statusQuery, assetID).Scan(&currentStatus); err != nil {
		return nil, err
	}

	if currentStatus == "RETIRED" || currentStatus == "LOST" {
		return nil, ErrAssetUnavailable
	}

	insertQuery := `
        INSERT INTO asset_maintenance_logs (
            asset_id,
            maintenance_type,
            description,
            sent_for_repair_at,
            vendor,
            notes,
            maint_status
        )
        VALUES ($1, $2, $3, NULLIF($4, '')::date, $5, '', 'IN_PROGRESS')
        RETURNING
            id,
            asset_id,
            maintenance_type,
            COALESCE(description, ''),
            COALESCE(sent_for_repair_at::text, ''),
            COALESCE(vendor, ''),
            COALESCE(notes, ''),
            maint_status,
            COALESCE(returned_from_repair_at::text, ''),
            COALESCE(repair_cost_inr, 0),
            created_at::text,
            updated_at::text
    `

	var record model.MaintenanceDTO
	if err := tx.QueryRow(
		ctx,
		insertQuery,
		assetID,
		req.MaintenanceType,
		req.Description,
		req.SentForRepairAt,
		req.Vendor,
	).Scan(
		&record.ID,
		&record.AssetID,
		&record.MaintenanceType,
		&record.Description,
		&record.SentForRepairAt,
		&record.Vendor,
		&record.Notes,
		&record.MaintStatus,
		&record.ReturnedFromRepairAt,
		&record.RepairCostINR,
		&record.CreatedAt,
		&record.UpdatedAt,
	); err != nil {
		return nil, err
	}

	updateAssetQuery := `
        UPDATE asset_inventory
        SET status = 'UNDER_REPAIR'
        WHERE id = $1 AND is_deleted = FALSE
    `

	tag, err := tx.Exec(ctx, updateAssetQuery, assetID)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	record.MaintStatus = "IN_PROGRESS"
	return &record, nil
}

func (r *Repository) UpdateMaintenanceRecord(ctx context.Context, assetID string, maintenanceID string, req model.UpdateMaintenanceRequest) (*model.MaintenanceDTO, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	nextAssetStatus := "UNDER_REPAIR"
	switch req.Status {
	case "COMPLETED":
		nextAssetStatus = "AVAILABLE"
	case "SCRAPPED":
		nextAssetStatus = "RETIRED"
	case "IN_PROGRESS":
		nextAssetStatus = "UNDER_REPAIR"
	}

	updateMaintenanceQuery := `
        UPDATE asset_maintenance_logs
        SET
            maint_status = $3,
            returned_from_repair_at = NULLIF($4, '')::date,
            repair_cost_inr = $5,
            vendor = $6,
            notes = $7,
            updated_at = NOW()
        WHERE id = $2 AND asset_id = $1
        RETURNING
            id,
            asset_id,
            maintenance_type,
            COALESCE(description, ''),
            COALESCE(sent_for_repair_at::text, ''),
            COALESCE(vendor, ''),
            COALESCE(notes, ''),
            maint_status,
            COALESCE(returned_from_repair_at::text, ''),
            COALESCE(repair_cost_inr, 0),
            created_at::text,
            updated_at::text
    `

	var record model.MaintenanceDTO
	if err := tx.QueryRow(
		ctx,
		updateMaintenanceQuery,
		assetID,
		maintenanceID,
		req.Status,
		req.ReturnedFromRepairAt,
		req.RepairCostINR,
		req.Vendor,
		req.Notes,
	).Scan(
		&record.ID,
		&record.AssetID,
		&record.MaintenanceType,
		&record.Description,
		&record.SentForRepairAt,
		&record.Vendor,
		&record.Notes,
		&record.MaintStatus,
		&record.ReturnedFromRepairAt,
		&record.RepairCostINR,
		&record.CreatedAt,
		&record.UpdatedAt,
	); err != nil {
		return nil, err
	}

	updateAssetQuery := `
        UPDATE asset_inventory
        SET status = $2
        WHERE id = $1 AND is_deleted = FALSE
    `

	tag, err := tx.Exec(ctx, updateAssetQuery, assetID, nextAssetStatus)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &record, nil
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

func (r *Repository) UpdateAsset(ctx context.Context, assetID string, asset model.Asset) error {
	query := `
        UPDATE asset_inventory
        SET
            asset_type = $2,
            name = $3,
            brand = $4,
            model = $5,
            category = $6,
            serial_no = $7,
            purchase_date = NULLIF($8, '')::date,
            purchase_cost_inr = $9,
            vendor = $10,
            warranty_expiry = NULLIF($11, '')::date,
            location = $12,
            notes = $13
        WHERE id = $1 AND is_deleted = FALSE
    `

	tag, err := r.DB.Exec(ctx, query,
		assetID,
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

func (r *Repository) AssignAsset(ctx context.Context, assetID string, req model.AssignRequest) (*model.AssignAssetDTO, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	statusQuery := `
        SELECT status
        FROM asset_inventory
        WHERE id = $1 AND is_deleted = FALSE
        FOR UPDATE
    `

	var currentStatus string
	if err := tx.QueryRow(ctx, statusQuery, assetID).Scan(&currentStatus); err != nil {
		return nil, err
	}

	if currentStatus != "AVAILABLE" {
		return nil, ErrAssetUnavailable
	}

	insertQuery := `
        INSERT INTO asset_assignments (
            asset_id,
            assigned_to,
            assigned_on,
            condition_at_assignment,
            notes,
            is_active,
            acknowledgement_status
        )
        VALUES ($1, $2, NULLIF($3, '')::timestamp, $4, $5, TRUE, 'PENDING')
        RETURNING id
    `

	var assignmentID string
	if err := tx.QueryRow(
		ctx,
		insertQuery,
		assetID,
		req.EmployeeID,
		req.AssignedOn,
		req.ConditionAtAssignment,
		req.Notes,
	).Scan(&assignmentID); err != nil {
		return nil, err
	}

	updateAssetQuery := `
        UPDATE asset_inventory
        SET status = 'ASSIGNED'
        WHERE id = $1 AND is_deleted = FALSE
    `

	tag, err := tx.Exec(ctx, updateAssetQuery, assetID)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &model.AssignAssetDTO{
		ID:                    assignmentID,
		AssetID:               assetID,
		EmployeeID:            req.EmployeeID,
		AssignedOn:            req.AssignedOn,
		ConditionAtAssignment: req.ConditionAtAssignment,
		Notes:                 req.Notes,
		AcknowledgementStatus: "PENDING",
		AssetStatus:           "ASSIGNED",
	}, nil
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
            COALESCE(assigned_on::text, created_at::text),
            COALESCE(condition_at_assignment, ''),
            COALESCE(notes, ''),
            COALESCE(acknowledgement_status, '')
        FROM asset_assignments
        WHERE asset_id = $1
        ORDER BY COALESCE(assigned_on, created_at) DESC
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
			&assignment.ConditionAtAssignment,
			&assignment.Notes,
			&assignment.AcknowledgementStatus,
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
