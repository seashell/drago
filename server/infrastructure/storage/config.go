package storage

import "github.com/seashell/drago/server/adapter/repository"

type Config struct {
	Type               repository.BackendType
	Path               string
	PostgreSQLAddress  string
	PostgreSQLPort     uint16
	PostgreSQLDatabase string
	PostgreSQLUsername string
	PostgreSQLPassword string
	PostgreSQLSSLMode  string
	SQLiteFilename     string
}
