package domain

import "time"

// Link :
type Link struct {
	ID                  *string
	FromInterfaceID     *string
	ToInterfaceID       *string
	AllowedIPs          []string
	PersistentKeepalive *int
	CreatedAt           *time.Time
	UpdatedAt           *time.Time
}

// LinkRepository :
type LinkRepository interface {
	GetByID(id string) (*Link, error)
	Create(l *Link) (*string, error)
	Update(l *Link) (*string, error)
	DeleteByID(id string) (*string, error)
	FindAll(pageInfo PageInfo) ([]*Link, *Page, error)
	FindAllByNetworkID(id string, pageInfo PageInfo) ([]*Link, *Page, error)
	FindAllBySourceHostID(id string, pageInfo PageInfo) ([]*Link, *Page, error)
	FindAllByTargetHostID(id string, pageInfo PageInfo) ([]*Link, *Page, error)
	FindAllBySourceInterfaceID(id string, pageInfo PageInfo) ([]*Link, *Page, error)
	FindAllByTargetInterfaceID(id string, pageInfo PageInfo) ([]*Link, *Page, error)
}

func (l *Link) Merge(b *Link) *Link {

	if b == nil {
		return l
	}

	result := *l

	return &result
}
