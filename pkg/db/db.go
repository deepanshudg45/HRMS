package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(connStr string) *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	return db
}
