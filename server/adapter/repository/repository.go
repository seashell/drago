package repository

import (
	"github.com/pkg/errors"
)

type BackendType string

const (
	// Maximum rows that can be returned in a query
	MaxQueryRows int = 50
	// In-memory storage backend type
	BackendInMemory BackendType = "inmem"
	// PostgreSQL storage backend type
	BackendPostgreSQL BackendType = "postgresql"
	// SQLite storage backend type
	BackendSQLite BackendType = "sqlite"
	// BoltDB storage backend type
	BackendBoltDB BackendType = "boltdb"
)

// Backend : Storage backend interface
type Backend interface {
	DB() interface{}
	Type() BackendType
}

// IsValid : Return and error in case the backend type is invalid
func (bt BackendType) IsValid() error {
	switch bt {
	case BackendInMemory, BackendPostgreSQL:
		return nil
	}
	return errors.New("Invalid backend type")
}
