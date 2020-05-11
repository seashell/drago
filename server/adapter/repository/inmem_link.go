package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/seashell/drago/server/domain"
)

type inmemLinkRepositoryAdapter struct {
	db *sqlx.DB
}

// NewInMemLinkRepositoryAdapter :
func NewInMemLinkRepositoryAdapter(backend Backend) (domain.LinkRepository, error) {
	if db, ok := backend.DB().(*sqlx.DB); ok {
		return &inmemLinkRepositoryAdapter{db}, nil
	}

	return nil, errors.New("Error creating in-memory backend adapter for link repository")
}

func (a *inmemLinkRepositoryAdapter) GetByID(id string) (*domain.Link, error) {
	return nil, nil
}

func (a *inmemLinkRepositoryAdapter) Create(l *domain.Link) (*string, error) {
	return nil, nil
}

func (a *inmemLinkRepositoryAdapter) Update(l *domain.Link) (*string, error) {
	return nil, nil
}

func (a *inmemLinkRepositoryAdapter) DeleteByID(id string) (*string, error) {
	return nil, nil
}

func (a *inmemLinkRepositoryAdapter) FindAllByNetworkID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error) {
	return nil, nil, nil
}

func (a *inmemLinkRepositoryAdapter) FindAllByHostID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error) {
	return nil, nil, nil
}
