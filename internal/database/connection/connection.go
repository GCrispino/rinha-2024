package connection

import (
	"context"
	"fmt"

	"github.com/GCrispino/rinha-2024/internal/utils"
	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConn struct {
	Conn *pgxpool.Pool
}

func NewDBConn(ctx context.Context, driverName, connString string, maxOpenConns int) (*DBConn, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("cannot parse config from conn string: %w", err)
	}

	cfg.MaxConns = int32(maxOpenConns)

	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to db: %w", err)
	}

	err = backoff.Retry(func() error {
		return db.Ping(ctx)
	}, utils.DefaultBackoff())
	if err != nil {
		return nil, err
	}

	return &DBConn{Conn: db}, err
}
