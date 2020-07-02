package sqlite

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/seashell/drago/server/adapter/repository"
)

// SQLiteBackend:
type SQLiteBackend struct {
	db *sqlx.DB
}

// NewSQLiteBackend:
func NewSQLiteBackend(filename string) (*SQLiteBackend, error) {
	db, err := sqlx.Connect("sqlite3", filename)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating new SQLite3 backend")
	}

	return &SQLiteBackend{db}, nil
}

// DB:
func (b *SQLiteBackend) DB() interface{} {
	return b.db
}

// Type:
func (b *SQLiteBackend) Type() repository.BackendType {
	return repository.BackendSQLite
}

// ApplyMigrations:
func (b *SQLiteBackend) ApplyMigrations(migrations ...string) error {
	for _, m := range migrations {
		b.db.MustExec(m)
	}
	return nil
}
