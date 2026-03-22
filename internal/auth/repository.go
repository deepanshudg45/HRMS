// internal/auth/repository.go
package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func (r *Repository) CreateEmployee(ctx context.Context, email, password, code string) error {
	_, err := r.DB.Exec(ctx,
		`INSERT INTO employees (work_email, password_hash, employee_code)
		 VALUES ($1,$2,$3)`,
		email, password, code,
	)
	return err
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (Employee, error) {
	var e Employee
	err := r.DB.QueryRow(ctx,
		`SELECT id, password_hash FROM employees WHERE work_email=$1`,
		email,
	).Scan(&e.ID, &e.PasswordHash)
	return e, err
}

func (r *Repository) UpdateRefreshToken(ctx context.Context, id string, tokenHash string, expiry string) error {
	_, err := r.DB.Exec(ctx,
		`UPDATE employees SET refresh_token=$1, refresh_expiry=$2 WHERE id=$3`,
		tokenHash, expiry, id,
	)
	return err
}

func (r *Repository) ClearRefreshToken(ctx context.Context, id string) error {
	_, err := r.DB.Exec(ctx,
		`UPDATE employees SET refresh_token=NULL, refresh_expiry=NULL WHERE id=$1`,
		id,
	)
	return err
}