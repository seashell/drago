package domain

import "time"

// Link :
type Link struct {
	ID                  *string    `json:"id"`
	FromInterfaceID     *string    `json:"fromInterface,omitempty"`
	ToInterfaceID       *string    `json:"toInterface,omitempty"`
	AllowedIPs          []string   `json:"allowedIps,omitempty"`
	PersistentKeepalive *int       `json:"persistentKeepalive,omitempty"`
	CreatedAt           *time.Time `json:"createdAt,omitempty"`
	UpdatedAt           *time.Time `json:"updatedAt,omitempty"`
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
