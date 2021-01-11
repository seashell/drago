package structs

type Interface struct {
	Name       string
	Address    string
	ListenPort int
	Peers      []*Peer
}

type Peer struct {
	PublicKey           string
	Address             string
	Port                int
	AllowedIPs          []string
	PersistentKeepalive int
}
