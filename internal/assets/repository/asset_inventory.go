package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

// ✅ WRITE IT HERE
func (r *Repository) NextAssetSeq(ctx context.Context) (int, error) {
	var seq int
	err := r.DB.QueryRow(ctx, "SELECT nextval('asset_code_seq')").Scan(&seq)
	return seq, err
}