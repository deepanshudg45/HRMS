package db

import (
	// "database/sql"
	"fmt"
	"os"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(pool *pgxpool.Pool) error {

	file, err := os.ReadFile("migrations/005_asset_inventory.up.sql")
	if err != nil {
		return err
	}

	query := string(file)


	_, err = pool.Exec(context.Background(), query)
	if err != nil {
		return err
	}

	fmt.Println("Database migrations applied ✅")

	return nil
}