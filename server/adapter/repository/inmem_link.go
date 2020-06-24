package repository

import (
	"github.com/pkg/errors"
	"github.com/seashell/drago/server/domain"
)

const bucketTypeLink = "link"

type inmemLinkRepositoryAdapter struct {
}

// NewInmemLinkRepositoryAdapter :
func NewInmemLinkRepositoryAdapter(backend Backend) (domain.LinkRepository, error) {
	return nil, errors.New("Not implemented")
}
