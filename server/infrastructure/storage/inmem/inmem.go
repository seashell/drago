package inmem

import (
	"github.com/seashell/drago/server/adapter/repository"
)

type InMemBackend struct {
	db map[string]interface{}
}

func NewInmemBackend() (*InMemBackend, error) {
	return &InMemBackend{
		db: make(map[string]interface{}),
	}, nil
}

func (b *InMemBackend) DB() interface{} {
	return b.db
}

func (b *InMemBackend) Type() repository.BackendType {
	return repository.BackendInMemory
}

func (b *InMemBackend) ApplyMigrations(migrations ...string) error {
	return nil
}
