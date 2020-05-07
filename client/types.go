package client

import "time"

type HostID = int

type KeyPair struct {
	PublicKey  string `json:"publicKey,omitempty"`
	PrivateKey string `json:"-"`
}

type WireguardInterface struct {
	Address    string `json:"address,omitempty"`
	ListenPort string `json:"listenPort,omitempty"`
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
	ID        HostID    `json:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`

	Name string `json:"name,omitempty"`

	Interface WireguardInterface `json:"interface,omitempty"`
	Peers     []WireguardPeer    `json:"peers,omitempty"`

	// Control fields
	Keys          KeyPair `json:"keys,omitempty"`
	AdvertiseAddr string  `json:"advertiseAddr,omitempty"`
	Jwt           string  `json:"jwt,omitempty"`
}
