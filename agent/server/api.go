package server

import "time"

type HostList struct {
	Count int            `json:"count"`
	Items []*HostSummary `json:"items"`
}

type HostSummary struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Address          string `json:"address"`
	PublicKey        string `json:"publicKey,omitempty"`
	AdvertiseAddress string `json:"advertiseAddress,omitempty"`
	ListenPort       string `json:"listenPort,omitempty"`
}

type HostDetails struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Address          string    `json:"address"`
	AdvertiseAddress string    `json:"advertiseAddress,omitempty"`
	ListenPort       string    `json:"listenPort,omitempty"`
	Table            string    `json:"table,omitempty"`
	DNS              string    `json:"dns,omitempty"`
	Mtu              string    `json:"mtu,omitempty"`
	PreUp            string    `json:"preUp,omitempty"`
	PostUp           string    `json:"postUp,omitempty"`
	PreDown          string    `json:"preDown,omitempty"`
	PostDown         string    `json:"postDown,omitempty"`
	PublicKey        string    `json:"publicKey,omitempty"`
	Links            LinkList  `json:"links,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type LinkList struct {
	Count int            `json:"count"`
	Items []*LinkDetails `json:"items"`
}

type LinkDetails struct {
	ID                  int          `json:"id"`
	From                *HostSummary `json:"from,omitempty"`
	To                  *HostSummary `json:"to,omitempty"`
	AllowedIPs          string       `json:"allowedIPs,omitempty"`
	PersistentKeepalive int          `json:"persistentKeepalive,omitempty"`
}

type LinkSummary struct {
	ID                  int    `json:"id"`
	FromID              int    `json:"fromId"`
	ToID                int    `json:"toId"`
	AllowedIPs          string `json:"allowedIPs,omitempty"`
	PersistentKeepalive int    `json:"persistentKeepalive,omitempty"`
}

type HostSettings struct {
	ID int `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Interface WireguardInterface `json:"interface"`
	Peers     []*WireguardPeer   `json:"peers,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type WireguardInterface struct {
	Address          string `json:"address"`
	AdvertiseAddress string `json:"advertiseAddress,omitempty"`
	ListenPort       string `json:"listenPort,omitempty"`
	Table            string `json:"table,omitempty"`
	DNS              string `json:"dns,omitempty"`
	Mtu              string `json:"mtu,omitempty"`
	PreUp            string `json:"preUp,omitempty"`
	PostUp           string `json:"postUp,omitempty"`
	PreDown          string `json:"preDown,omitempty"`
	PostDown         string `json:"postDown,omitempty"`
}

type WireguardPeer struct {
	Name                string `json:"name"`
	PublicKey           string `json:"publicKey,omitempty"`
	AllowedIPs          string `json:"allowedIPs,omitempty"`
	Endpoint            string `json:"endpoint,omitempty"`
	PersistentKeepalive int    `json:"persistentKeepalive,omitempty"`
}

type ApiError struct {
	Message string `json:"message"`
}

type CreateHostInput struct {
	Name             string `json:"name"`
	Address          string `json:"address"`
	AdvertiseAddress string `json:"advertiseAddress,omitempty"`
	ListenPort       string `json:"listenPort,omitempty"`
}

type UpdateHostInput struct {
	ID               int    `json:"id,omitempty"`
	Name             string `json:"name"`
	Address          string `json:"address"`
	AdvertiseAddress string `json:"advertiseAddress,omitempty"`
	ListenPort       string `json:"listenPort,omitempty"`
	Table            string `json:"table,omitempty"`
	DNS              string `json:"dns,omitempty"`
	Mtu              string `json:"mtu,omitempty"`
	PreUp            string `json:"preUp,omitempty"`
	PostUp           string `json:"postUp,omitempty"`
	PreDown          string `json:"preDown,omitempty"`
	PostDown         string `json:"postDown,omitempty"`
	PublicKey        string `json:"publicKey,omitempty"`
}

type DeleteHostInput struct {
	ID int `json:"id,omitempty"`
}

type GetAllHostsInput struct {
}

type GetHostInput struct {
	ID int `json:"id,omitempty"`
}

type CreateLinkInput struct {
	FromID              int    `json:"from,omitempty"`
	ToID                int    `json:"to,omitempty"`
	AllowedIPs          string `json:"allowedIPs,omitempty"`
	PersistentKeepalive int    `json:"persistentKeepalive,omitempty"`
}

type UpdateLinkInput struct {
	ID                  int    `json:"id"`
	AllowedIPs          string `json:"allowedIPs,omitempty"`
	PersistentKeepalive int    `json:"persistentKeepalive,omitempty"`
}

type DeleteLinkInput struct {
	ID int `json:"id"`
}

type GetAllLinksInput struct {
}

type SyncHostInput struct {
	ID        int
	HostType  string
	PublicKey string `json:"publicKey"`
}
