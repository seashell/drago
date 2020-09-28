package structs

import "time"

// Host :
type Host struct {
	Name       string       `json:"name"`
	Interfaces []*Interface `json:"interfaces,omitempty"`
	Peers      []*Peer      `json:"peers,omitempty"`
}

// Interface :
type Interface struct {
	Name       string `json:"name,omitempty"`
	PublicKey  string `json:"publicKey,omitempty"`
	Address    string `json:"address,omitempty"`
	ListenPort string `json:"listenPort,omitempty"`
	Table      string `json:"table,omitempty"`
	DNS        string `json:"dns,omitempty"`
	MTU        string `json:"mtu,omitempty"`
	PreUp      string `json:"preUp,omitempty"`
	PostUp     string `json:"postUp,omitempty"`
	PreDown    string `json:"preDown,omitempty"`
	PostDown   string `json:"postDown,omitempty"`
}

// Peer :
type Peer struct {
	Interface           string    `json:"interface,omitempty"`
	Address             string    `json:"address,omitempty"`
	Port                string    `json:"port,omitempty"`
	PublicKey           string    `json:"publicKey,omitempty"`
	AllowedIPs          []string  `json:"allowedIps"`
	PersistentKeepalive int       `json:"persistentKeepalive,omitempty"`
	LatencyMs           uint64    `json:"latencyMs,omitempty"`
	LastHandshake       time.Time `json:"lastHandshake,omitempty"`
}

// HostSynchronizeInput :
type HostSynchronizeInput struct {
	Host
}

// HostSynchronizeOutput :
type HostSynchronizeOutput struct {
	Host
}
