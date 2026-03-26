package database

import (
	"context"
	"errors"

	"WITS/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
    pool  *pgxpool.Pool
}

func Connect(cfg config.Config) (*Client, error) {
    if cfg.DBURL == "" {
        return nil, errors.New("DATABASE_URL is required")
    }

    pool, err := pgxpool.New(context.Background(), cfg.DBURL)
    if err != nil {
        return nil, err
    }

	if err := pool.Ping(context.Background()); err != nil {
        pool.Close()
        return nil, err
    }

    return &Client{pool: pool}, nil
}

func (c *Client) Exec(ctx context.Context, sql string, arguments ...any) (any, error) {
    return c.pool.Exec(ctx, sql, arguments...)
}

func (c *Client) Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error) {
    return c.pool.Query(ctx, sql, arguments...)
}

func (c *Client) QueryRow(ctx context.Context, sql string, arguments ...any) interface{ Scan(dest ...any) error } {
    return c.pool.QueryRow(ctx, sql, arguments...)
}

func (c *Client) Close() error {
	c.pool.Close()
    return nil
}