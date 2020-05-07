package storage

import (
	"errors"

	"github.com/seashell/drago/server/adapter/repository"
	"github.com/seashell/drago/server/infrastructure/storage/inmem"
	"github.com/seashell/drago/server/infrastructure/storage/postgres"
	"github.com/seashell/drago/server/infrastructure/storage/sqlite"
)

func NewBackend(conf *Config) (repository.Backend, error) {

	switch conf.Type {
	case repository.BackendPostgreSQL:
		return postgres.NewPostgreSQLBackend(
			conf.PostgreSQLAddress,
			conf.PostgreSQLPort,
			conf.PostgreSQLDatabase,
			conf.PostgreSQLUsername,
			conf.PostgreSQLPassword,
			conf.PostgreSQLSSLMode,
		)

	case repository.BackendInMemory:
		return inmem.NewInmemBackend()

	case repository.BackendSQLite:
		return sqlite.NewSQLiteBackend(conf.SQLiteFilename)

	default:
		return nil, errors.New("Invalid backend")
	}
}
