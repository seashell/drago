package domain

import "time"

// Host :
type Host struct {
	ID               *string
	Name             *string
	AdvertiseAddress *string
	Labels           []string
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
}

// HostRepository :
type HostRepository interface {
	GetByID(id string) (*Host, error)
	Create(h *Host) (*string, error)
	CreateWithID(h *Host) (*string, error)
	Update(h *Host) (*string, error)
	DeleteByID(id string) (*string, error)
	FindAll(pageInfo PageInfo) ([]*Host, *Page, error)
	FindAllByNetworkID(id string, pageInfo PageInfo) ([]*Host, *Page, error)
	FindAllByLabels(labels []string, pageInfo PageInfo) ([]*Host, *Page, error)
	ExistsByID(id string) (bool, error)
}

func (h *Host) Merge(b *Host) *Host {

	if b == nil {
		return h
	}

	result := *h

	return &result
}
