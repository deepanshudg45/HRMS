package repository

import (
	"context"

	"WITS/internal/assets/model"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetActiveAssetsByEmployeeID(ctx context.Context, employeeID string) ([]model.MyAssetDTO, error) {
	query := `
        SELECT
            ai.asset_code,
            ai.asset_type,
            ai.name,
            aa.assigned_on,
            aa.acknowledgement_status,
            aa.condition_at_assignment
        FROM asset_assignments aa
        JOIN asset_inventory ai ON ai.id = aa.asset_id
        WHERE aa.assigned_to = $1
          AND aa.is_active = TRUE
          AND ai.is_deleted = FALSE
        ORDER BY aa.assigned_on DESC
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
			&asset.AcknowledgementStatus,
			&asset.ConditionAtAssignment,
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

func (r *Repository) UpdateAssignmentAcknowledgement(ctx context.Context, assignmentID string) (*model.AssignmentDTO, error) {
	query := `
        UPDATE asset_assignments
        SET
            acknowledgement_status = 'ACKNOWLEDGED',
            acknowledged_at = NOW()
        WHERE id = $1
        RETURNING
            id,
            assigned_to,
            acknowledgement_status,
            COALESCE(acknowledged_at::text, '')
    `

	var assignment model.AssignmentDTO
	err := r.DB.QueryRow(ctx, query, assignmentID).Scan(
		&assignment.ID,
		&assignment.EmployeeID,
		&assignment.AcknowledgementStatus,
		&assignment.AcknowledgedAt,
	)
	if err != nil {
		return nil, err
	}

	return &assignment, nil
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
