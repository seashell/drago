package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/seashell/drago/server/domain"
)

type inmemNetworkRepositoryAdapter struct {
	db *sqlx.DB
}

// NewInMemNetworkRepositoryAdapter :
func NewInMemNetworkRepositoryAdapter(backend Backend) (domain.NetworkRepository, error) {
	if db, ok := backend.DB().(*sqlx.DB); ok {
		return &inmemNetworkRepositoryAdapter{db}, nil
	}

	return nil, errors.New("Error creating in-memory backend adapter for network repository")
}

func (a *inmemNetworkRepositoryAdapter) GetByID(id string) (*domain.Network, error) {
	return nil, nil
}

func (a *inmemNetworkRepositoryAdapter) Create(h *domain.Network) (id *string, err error) {
	return nil, nil
}

func (a *inmemNetworkRepositoryAdapter) Update(h *domain.Network) (id *string, err error) {
	return nil, nil
}

func (a *inmemNetworkRepositoryAdapter) DeleteByID(id string) error {
	return nil
}

func (a *inmemNetworkRepositoryAdapter) FindAll(pageInfo domain.PageInfo) ([]*domain.Network, *domain.Page, error) {
	return nil, nil, nil
}
