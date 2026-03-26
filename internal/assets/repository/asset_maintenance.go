package repository

import (
	"context"

	"WITS/internal/assets/model"

	"github.com/jackc/pgx/v5"
)

type Row interface {
    Scan(dest ...any) error
}

type DB interface {
    Exec(ctx context.Context, sql string, arguments ...any) (any, error)
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
