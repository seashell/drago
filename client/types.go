package client

import "strings"

type HostID = int

type Host struct {
	Name          string               `json:"name,omitempty"`
	AdvertiseAddr *string              `json:"advertiseAddr,omitempty"`
	Interfaces    []WireguardInterface `json:"interface,omitempty"`
	Peers         []WireguardPeer      `json:"peers,omitempty"`
}

type WireguardInterface struct {
	Name       string  `json:"name,omitempty"`
	Address    string  `json:"address,omitempty"`
	ListenPort *string `json:"listenPort,omitempty"`
	Table      *string `json:"table,omitempty"`
	DNS        *string `json:"dns,omitempty"`
	MTU        *string `json:"mtu,omitempty"`
	PreUp      *string `json:"preUp,omitempty"`
	PostUp     *string `json:"postUp,omitempty"`
	PreDown    *string `json:"preDown,omitempty"`
	PostDown   *string `json:"postDown,omitempty"`
	KeyPair    *string `json:"-"`
}

type KeyPair struct {
	PublicKey  string `json:"publicKey,omitempty"`
	PrivateKey string `json:"-"`
}

type WireguardPeer struct {
	Endpoint            *string  `json:"endpoint,omitempty"`
	AllowedIPs          []string `json:"allowedIPs,omitempty"`
	PublicKey           *string  `json:"publicKey,omitempty"`
	PersistentKeepalive *int     `json:"persistentKeepalive,omitempty"`
}

func (WireguardPeer) AllowedIPsToString(allowedIPs []string) string {
	return strings.Join(allowedIPs, ",")
}
