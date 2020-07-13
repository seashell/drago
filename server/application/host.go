package application

import (
	"github.com/seashell/drago/server/domain"
)

// HostService :
type HostService interface {
	GetByID(id string) (*domain.Host, error)
	Create(h *domain.Host) (*domain.Host, error)
	Update(h *domain.Host) (*domain.Host, error)
	DeleteByID(id string) (*domain.Host, error)
	FindAll(pageInfo domain.PageInfo) ([]*domain.Host, *domain.Page, error)
	FindAllByLabels(labels []string, pageInfo domain.PageInfo) ([]*domain.Host, *domain.Page, error)
	FindAllByNetworkID(id string, pageInfo domain.PageInfo) ([]*domain.Host, *domain.Page, error)
}

type hostService struct {
	repo domain.HostRepository
}

// NewHostService :
func NewHostService(repo domain.HostRepository) (HostService, error) {
	return &hostService{repo}, nil
}

// GetByID :
func (s *hostService) GetByID(id string) (*domain.Host, error) {
	return s.repo.GetByID(id)
}

// Create :
func (s *hostService) Create(h *domain.Host) (*domain.Host, error) {

	id, err := s.repo.Create(h)
	if err != nil {
		return nil, err
	}

	return &domain.Host{ID: id}, nil
}

// Update :
func (s *hostService) Update(h *domain.Host) (*domain.Host, error) {
	host, err := s.repo.GetByID(*h.ID)
	if err != nil {
		return nil, err
	}
	mergeHostUpdate(host, h)

	id, err := s.repo.Update(host)
	if err != nil {
		return nil, err
	}

	return &domain.Host{ID: id}, nil
}

// Delete :
func (s *hostService) DeleteByID(id string) (*domain.Host, error) {
	_id, err := s.repo.DeleteByID(id)
	if err != nil {
		return nil, err
	}
	return &domain.Host{ID: _id}, nil
}

// FindAll :
func (s *hostService) FindAll(pageInfo domain.PageInfo) ([]*domain.Host, *domain.Page, error) {
	return s.repo.FindAll(pageInfo)
}

// FindAllByNetworkID :
func (s *hostService) FindAllByNetworkID(id string, pageInfo domain.PageInfo) ([]*domain.Host, *domain.Page, error) {
	return s.repo.FindAllByNetworkID(id, pageInfo)
}

// FindAllByLabels :
func (s *hostService) FindAllByLabels(labels []string, pageInfo domain.PageInfo) ([]*domain.Host, *domain.Page, error) {
	return s.repo.FindAllByLabels(labels, pageInfo)
}

func mergeHostUpdate(current, update *domain.Host) {
	if update.Name != nil {
		current.Name = update.Name
	}
	if update.Labels != nil {
		current.Labels = update.Labels
	}
	if update.AdvertiseAddress != nil {
		current.AdvertiseAddress = update.AdvertiseAddress
	}
}
