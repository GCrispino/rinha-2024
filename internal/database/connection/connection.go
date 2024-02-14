package connection

import (
	"database/sql"
	"fmt"

	"github.com/GCrispino/rinha-2024/internal/utils"
	"github.com/cenkalti/backoff/v4"
	_ "github.com/lib/pq"
)

type DBConn struct {
	Conn *sql.DB
}

func NewDBConn(driverName, connString string, maxOpenConns int) (*DBConn, error) {
	db, err := sql.Open(driverName, connString)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to db: %w", err)
	}

	db.SetMaxOpenConns(maxOpenConns)

	err = backoff.Retry(func() error {
		return db.Ping()
	}, utils.DefaultBackoff())
	if err != nil {
		return nil, err
	}

	return &DBConn{Conn: db}, err
}
