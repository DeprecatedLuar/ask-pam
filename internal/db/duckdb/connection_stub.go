//go:build !cgo

package duckdb

import (
	"fmt"

	"github.com/eduardofuncao/squix/internal/db"
)

type Connection struct {
	*db.BaseConnection
}

func New(name, connStr string) (db.DatabaseConnection, error) {
	return nil, fmt.Errorf("duckdb driver not available: build with CGO_ENABLED=1 to enable")
}

func (d *Connection) GetUniqueConstraints(tableName string) ([]string, error) {
	return nil, fmt.Errorf("duckdb driver not available")
}

func init() {
	db.Register("duckdb", New)
}
