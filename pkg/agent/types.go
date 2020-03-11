package agent

import "time"

type HostID = int

type Entity struct {
	ID        HostID    `json:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type WireguardInterface struct {
	Address    string `json:"address,omitempty"`
	ListenPort string `json:"listenPort,omitempty"`
	PrivateKey string `json:"privateKey,omitempty"`
	Table      string `json:"table,omitempty"`
	DNS        string `json:"DNS,omitempty"`
	MTU        string `json:"MTU,omitempty"`
	PreUp      string `json:"PreUp,omitempty"`
	PostUp     string `json:"PostUp,omitempty"`
	PreDown    string `json:"PreDown,omitempty"`
	PostDown   string `json:"PostDown,omitempty"`
}

type WireguardPeer struct {
	Endpoint            string `json:"endpoint,omitempty"`
	AllowedIPs          string `json:"allowedIPs,omitempty"`
	PublicKey           string `json:"publicKey,omitempty"`
	PersistentKeepalive int    `json:"persistentKeepalive,omitempty"`
}

type Host struct {
	Entity
	Name string `json:"name,omitempty"`

	Interface WireguardInterface `json:"interface,omitempty"`
	Peers     []WireguardPeer    `json:"peers,omitempty"`

	// Control fields
	PublicKey     string `json:"publicKey,omitempty"`
	AdvertiseAddr string `json:"advertiseAddr,omitempty"`
	Jwt           string `json:"jwt,omitempty"`
}

type Version struct {
	Version string `json:"version"`
}

type Error struct {
	Message string `json:"message"`
}
