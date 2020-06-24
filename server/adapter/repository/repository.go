package repository

import (
	"github.com/pkg/errors"
)

type BackendType string

const (
	MaxQueryRows      int         = 50
	BackendInMemory   BackendType = "inmem"
	BackendPostgreSQL BackendType = "postgresql"
	BackendSQLite     BackendType = "sqlite"
	BackendBoltDB     BackendType = "boltdb"
)

type Backend interface {
	DB() interface{}
	Type() BackendType
}

func (bt BackendType) IsValid() error {
	switch bt {
	case BackendInMemory, BackendPostgreSQL:
		return nil
	}
	return errors.New("Invalid backend type")
}
