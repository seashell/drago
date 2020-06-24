package domain

import "time"

// Host :
type Host struct {
	ID               *string    `json:"id"`
	Name             *string    `json:"name,omitempty"`
	AdvertiseAddress *string    `json:"advertiseAddress,omitempty"`
	Labels           []string   `json:"labels,omitempty"`
	CreatedAt        *time.Time `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time `json:"updatedAt,omitempty"`
}

// HostSettings :
type HostSettings struct {
	Interfaces []*WgInterfaceSettings `json:"interfaces,omitempty"`
	Peers      []*WgPeerSettings      `json:"peers,omitempty"`
}

// HostState :
type HostState struct {
	Interfaces []*WgInterfaceState `json:"interfaces,omitempty"`
	Peers      []*WgPeerState      `json:"peers,omitempty"`
}

// WgInterfaceState :
type WgInterfaceState struct {
	Name      *string `json:"name,omitempty"`
	PublicKey *string `json:"publicKey,omitempty"`
}

// WgPeerState :
type WgPeerState struct {
	LatencyMs     uint64    `json:"latencyMs,omitempty"`
	LastHandshake time.Time `json:"lastHandshake,omitempty"`
}

// WgInterfaceSettings :
type WgInterfaceSettings struct {
	Name       *string `json:"name,omitempty"`
	Address    *string `json:"address,omitempty"`
	ListenPort *string `json:"listenPort,omitempty"`
	Table      *string `json:"table,omitempty"`
	DNS        *string `json:"dns,omitempty"`
	MTU        *string `json:"mtu,omitempty"`
	PreUp      *string `json:"preUp,omitempty"`
	PostUp     *string `json:"postUp,omitempty"`
	PreDown    *string `json:"preDown,omitempty"`
	PostDown   *string `json:"postDown,omitempty"`
}

// WgPeerSettings :
type WgPeerSettings struct {
	Interface           string   `json:"interface,omitempty"`
	Address             *string  `json:"address,omitempty"`
	Port                *string  `json:"port,omitempty"`
	PublicKey           *string  `json:"publicKey,omitempty"`
	AllowedIPs          []string `json:"allowedIps"`
	PersistentKeepalive *int     `json:"persistentKeepalive,omitempty"`
}

// HostRepository :
type HostRepository interface {
	GetByID(id string) (*Host, error)
	Create(h *Host) (*string, error)
	Update(h *Host) (*string, error)
	DeleteByID(id string) (*string, error)
	FindAll(pageInfo PageInfo) ([]*Host, *Page, error)
	FindAllByNetworkID(id string, pageInfo PageInfo) ([]*Host, *Page, error)
}
