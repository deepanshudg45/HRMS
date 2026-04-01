package database

import (
	"context"
	"errors"

	"WITS/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	pool *pgxpool.Pool
}

func Connect(cfg config.Config) (*Client, error) {
	if cfg.DBURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &Client{pool: pool}, nil
}

func (c *Client) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return c.pool.Exec(ctx, sql, arguments...)
}

func (c *Client) Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error) {
	return c.pool.Query(ctx, sql, arguments...)
}

func (c *Client) QueryRow(ctx context.Context, sql string, arguments ...any) interface{ Scan(dest ...any) error } {
	return c.pool.QueryRow(ctx, sql, arguments...)
}

func (c *Client) Begin(ctx context.Context) (pgx.Tx, error) {
	return c.pool.Begin(ctx)
}

func (c *Client) Close() {
	c.pool.Close()
}
