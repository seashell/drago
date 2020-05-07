package repository

import (
	"github.com/pkg/errors"
)

type BackendType string

const (
	maxQueryRows int = 50

	BackendInMemory   BackendType = "inmem"
	BackendPostgreSQL BackendType = "postgresql"
	BackendSQLite     BackendType = "sqlite"
	BackendBoltDB     BackendType = "boltdb"
)

type Backend interface {
	DB() interface{}
	Type() BackendType
	ApplyMigrations(...string) error
}

func (bt BackendType) IsValid() error {
	switch bt {
	case BackendInMemory, BackendPostgreSQL, BackendSQLite, BackendBoltDB:
		return nil
	}

	return errors.New("Invalid backend type")
}
