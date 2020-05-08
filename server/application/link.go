package application

import (
	"github.com/seashell/drago/server/domain"
)

// LinkService :
type LinkService interface {
	GetByID(id string) (*domain.Link, error)
	Create(l *domain.Link) (*domain.Link, error)
	Update(l *domain.Link) (*domain.Link, error)
	DeleteByID(id string) error
	FindAllByNetworkID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error)
	FindAllByHostID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error)
}

type linkService struct {
	lr domain.LinkRepository
}

// NewLinkService :
func NewLinkService(lr domain.LinkRepository) (LinkService, error) {
	return &linkService{lr}, nil
}

// GetByID :
func (ls *linkService) GetByID(id string) (*domain.Link, error) {
	return ls.lr.GetByID(id)
}

// Create :
func (ls *linkService) Create(l *domain.Link) (*domain.Link, error) {
	id, err := ls.lr.Create(l)
	if err != nil {
		return nil, err
	}
	return &domain.Link{ID: id}, nil
}

// Update :
func (ls *linkService) Update(l *domain.Link) (*domain.Link, error) {
	link, err := ls.lr.GetByID(*l.ID)
	if err != nil {
		return nil, err
	}

	mergeLinkUpdate(link, l)

	id, err := ls.lr.Update(link)
	if err != nil {
		return nil, err
	}
	return &domain.Link{ID: id}, nil
}

// Delete :
func (ls *linkService) DeleteByID(id string) error {
	return ls.lr.DeleteByID(id)
}

// FindAllByNetworkID :
func (ls *linkService) FindAllByNetworkID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error) {
	return ls.lr.FindAllByNetworkID(id, pageInfo)
}

// FindAllByHostID :
func (ls *linkService) FindAllByHostID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error) {
	return ls.lr.FindAllByHostID(id, pageInfo)
}

func mergeLinkUpdate(current, update *domain.Link) {
	if update.AllowedIPs != nil {
		current.AllowedIPs = update.AllowedIPs
	}

	if update.PersistentKeepalive != nil {
		current.PersistentKeepalive = update.PersistentKeepalive
	}
}
