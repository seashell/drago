package inmem

import (
	"github.com/seashell/drago/server/adapter/repository"
)

// InmemBackend:
type InMemBackend struct {
	db map[string]interface{}
}

// NewInmemBackend:
func NewInmemBackend() (*InMemBackend, error) {
	return &InMemBackend{
		db: make(map[string]interface{}),
	}, nil
}

// DB:
func (b *InMemBackend) DB() interface{} {
	return b.db
}

// Type:
func (b *InMemBackend) Type() repository.BackendType {
	return repository.BackendInMemory
}

// ApplyMigrations:
func (b *InMemBackend) ApplyMigrations(migrations ...string) error {
	return nil
}
