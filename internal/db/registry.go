package db

import (
	"fmt"
)

// Constructor creates a DatabaseConnection for a registered engine.
type Constructor func(name, connString string) (DatabaseConnection, error)

var registry = map[string]Constructor{}

// Register associates a dbType (or alias) with a Constructor. Engine packages
// call this from their init() so registration happens before any
// CreateConnection call, as long as the engine package is imported
// (see internal/db/engines).
func Register(dbType string, c Constructor) {
	registry[dbType] = c
}

// CreateConnection looks up the constructor registered for dbType and
// invokes it with the given name and connection string.
func CreateConnection(name, dbType, connString string) (DatabaseConnection, error) {
	c, ok := registry[dbType]
	if !ok {
		return nil, fmt.Errorf("driver not implemented for %s", dbType)
	}
	return c(name, connString)
}
