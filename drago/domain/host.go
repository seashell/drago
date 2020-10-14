package domain

import (
	"context"
	"time"
)

// Host ...
type Host struct {
	ID               string
	Name             string
	AdvertiseAddress string
	Labels           []string
	Interfaces       map[string]*Interface
	Links            []*Link
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Interface ...
type Interface struct {
	Name       string
	HostID     string
	NetworkID  string
	IPAddress  string
	ListenPort string
	Table      string
	DNS        string
	MTU        string
	PreUp      string
	PostUp     string
	PreDown    string
	PostDown   string
	PublicKey  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Link ...
type Link struct {
	ToHostID            string
	ToInterfaceName     string
	AllowedIPs          []string
	PersistentKeepalive int
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// HostRepository : Host repository interface
type HostRepository interface {
	GetByID(ctx context.Context, id string) (*Host, error)
	Create(ctx context.Context, h *Host) (*string, error)
	Update(ctx context.Context, h *Host) (*string, error)
	DeleteByID(ctx context.Context, id string) (*string, error)
	FindAll(ctx context.Context) ([]*Host, error)
	FindByLabels(ctx context.Context, labels []string) ([]*Host, error)
}

func (h *Host) UpsertInterface(i *Interface) error {
	return nil
}

func (h *Host) DeleteInterface(i string) error {
	return nil
}

func (h *Host) UpsertLink(l *Link) error {
	return nil
}

func (h *Host) DeleteLink(l string) error {
	return nil
}
