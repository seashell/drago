package application

import (
	"github.com/seashell/drago/server/domain"
)

// LinkService :
type LinkService interface {
	GetByID(id string) (*domain.Link, error)
	Create(l *domain.Link) (*domain.Link, error)
	Update(l *domain.Link) (*domain.Link, error)
	DeleteByID(id string) (*domain.Link, error)
	FindAll(pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error)
	FindAllByNetworkID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error)
	FindAllBySourceHostID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error)
	FindAllByTargetHostID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error)
	FindAllByTargetInterfaceID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error)
	FindAllBySourceInterfaceID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error)
}

type linkService struct {
	repo domain.LinkRepository
}

// NewLinkService :
func NewLinkService(lr domain.LinkRepository) (LinkService, error) {
	return &linkService{lr}, nil
}

// GetByID :
func (s *linkService) GetByID(id string) (*domain.Link, error) {
	return s.repo.GetByID(id)
}

// Create :
func (s *linkService) Create(l *domain.Link) (*domain.Link, error) {
	id, err := s.repo.Create(l)
	if err != nil {
		return nil, err
	}
	return &domain.Link{ID: id}, nil
}

// Update :
func (s *linkService) Update(l *domain.Link) (*domain.Link, error) {
	link, err := s.repo.GetByID(*l.ID)
	if err != nil {
		return nil, err
	}

	mergeLinkUpdate(link, l)

	id, err := s.repo.Update(link)
	if err != nil {
		return nil, err
	}
	return &domain.Link{ID: id}, nil
}

// Delete :
func (s *linkService) DeleteByID(id string) (*domain.Link, error) {
	_id, err := s.repo.DeleteByID(id)
	if err != nil {
		return nil, err
	}
	return &domain.Link{ID: _id}, nil
}

// FindAll :
func (s *linkService) FindAll(pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error) {
	return s.repo.FindAll(pageInfo)
}

// FindAllByNetworkID :
func (s *linkService) FindAllByNetworkID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error) {
	return s.repo.FindAllByNetworkID(id, pageInfo)
}

// FindAllBySourceHostID :
func (s *linkService) FindAllBySourceHostID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error) {
	return s.repo.FindAllBySourceHostID(id, pageInfo)
}

// FindAllByTargetHostID :
func (s *linkService) FindAllByTargetHostID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error) {
	return s.repo.FindAllByTargetHostID(id, pageInfo)
}

// FindAllBySourceInterfaceID :
func (s *linkService) FindAllBySourceInterfaceID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error) {
	return s.repo.FindAllBySourceInterfaceID(id, pageInfo)
}

// FindAllByTargetInterfaceID :
func (s *linkService) FindAllByTargetInterfaceID(id string, pageInfo domain.PageInfo) ([]*domain.Link, *domain.Page, error) {
	return s.repo.FindAllByTargetInterfaceID(id, pageInfo)
}

func mergeLinkUpdate(current, update *domain.Link) {
	if update.AllowedIPs != nil {
		current.AllowedIPs = update.AllowedIPs
	}
	if update.PersistentKeepalive != nil {
		current.PersistentKeepalive = update.PersistentKeepalive
	}
}
