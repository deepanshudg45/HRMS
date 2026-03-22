package db

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
	once sync.Once
)

func InitDB(connStr string) {
	once.Do(func() {
		config, err := pgxpool.ParseConfig(connStr)
		if err != nil {
			log.Fatal(err)
		}

		config.MaxConns = 25

		pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			log.Fatal(err)
		}
	})
}

func GetPool() *pgxpool.Pool {
	return pool
}