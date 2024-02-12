package connection

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DBConn struct {
	Conn *sql.DB
}

func NewDBConn(driverName, connString string) (*DBConn, error) {
	db, err := sql.Open(driverName, connString)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DBConn{Conn: db}, err
}
