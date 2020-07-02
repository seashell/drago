package repository

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/seashell/drago/server/domain"
)

// NewLinkRepositoryAdapter :
func NewLinkRepositoryAdapter(backend Backend) (domain.LinkRepository, error) {
	switch backend.Type() {
	case BackendPostgreSQL:
		return NewPostgreSQLLinkRepositoryAdapter(backend)

	case BackendInMemory:
		return NewInmemLinkRepositoryAdapter(backend)

	default:
		return nil, errors.New(fmt.Sprintf("Error creating adapter for backend of type %s", backend.Type()))
	}
}
