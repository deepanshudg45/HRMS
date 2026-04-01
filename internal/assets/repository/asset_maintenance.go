package repository

import (
	"context"
	"errors"

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
