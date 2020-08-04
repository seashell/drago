package boltdb

import (
	"fmt"
)

// Backend:
type Backend struct{}

// NewBoltBackend:
func NewBackend(filename string) (*Backend, error) {
	return nil, nil
}

// DB:
func (b *Backend) DB() interface{} {
	return nil
}

// ApplyMigrations:
func (b *Backend) ApplyMigrations(migrations ...string) error {
	for _, m := range migrations {
		fmt.Println(m)
	}
	return nil
}
