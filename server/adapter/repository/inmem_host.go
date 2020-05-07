package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/seashell/drago/server/domain"
)

type inmemHostRepositoryAdapter struct {
	db *sqlx.DB
}

// NewInMemHostRepositoryAdapter :
func NewInMemHostRepositoryAdapter(backend Backend) (domain.HostRepository, error) {
	if db, ok := backend.DB().(*sqlx.DB); ok {
		return &inmemHostRepositoryAdapter{db}, nil
	}

	return nil, errors.New("Error creating in-memory backend adapter for host repository")
}

func (a *inmemHostRepositoryAdapter) GetByID(id string) (*domain.Host, error) {
	return nil, nil
}

func (a *inmemHostRepositoryAdapter) Create(h *domain.Host) (id *string, err error) {
	return nil, nil
}

func (a *inmemHostRepositoryAdapter) Update(h *domain.Host) (id *string, err error) {
	return nil, nil
}

func (a *inmemHostRepositoryAdapter) DeleteByID(id string) error {
	return nil
}

func (a *inmemHostRepositoryAdapter) FindAllByNetworkID(id string, pageInfo domain.PageInfo) ([]*domain.Host, *domain.Page, error) {
	return nil, nil, nil
}
