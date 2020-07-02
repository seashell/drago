package application

import (
	"github.com/google/uuid"
	"github.com/seashell/drago/server/domain"
)

// AdmissionService :
type AdmissionService interface {
	GetHostOrCreate(in *domain.Host) (*domain.Host, error)
}

type admissionService struct {
	hr domain.HostRepository
}

// NewAdmissionService :
func NewAdmissionService(hs domain.HostRepository) (AdmissionService, error) {
	return &admissionService{hs}, nil
}

// GetHostOrCreate :
func (s *admissionService) GetHostOrCreate(h *domain.Host) (*domain.Host, error) {
	exists, err := s.hr.ExistsByID(*h.ID)
	if err != nil {
		return nil, err
	}

	if exists {
		h, err := s.hr.GetByID(*h.ID)
		if err != nil {
			return nil, err
		}
		if h != nil {
			return h, nil
		}
	}

	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	sguid := uid.String()

	host := &domain.Host{
		ID:     h.ID,
		Name:   &sguid,
		Labels: h.Labels,
	}

	id, err := s.hr.CreateWithID(host)
	if err != nil {
		return nil, err
	}

	return &domain.Host{ID: id}, nil
}
