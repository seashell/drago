package repository

import (
	"github.com/pkg/errors"
	"github.com/seashell/drago/server/domain"
)

const bucketTypeInterface = "interface"

type inmemInterfaceRepositoryAdapter struct {
}

// NewInmemInterfaceRepositoryAdapter :
func NewInmemInterfaceRepositoryAdapter(backend Backend) (domain.InterfaceRepository, error) {
	return nil, errors.New("Not implemented")
}
