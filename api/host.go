package api

//TODO refactor names

type Settings struct {
	Interfaces []Interface 	`json:"interfaces"`
	Peers []Peer 				`json:"peers"`
}

type Interface struct {
	Name 		string `json:"name"`
	ListenPort 	string `json:"listenPort"`
	Address 	string `json:"address"`
	Table      	string `json:"table"`
	DNS        	string `json:"dns"`
	MTU        	string `json:"mtu"`
	PreUp      	string `json:"preUp"`
	PostUp     	string `json:"postUp"`
	PreDown    	string `json:"preDown"`
	PostDown   	string `json:"postDown"`
}

type Peer struct {
	Interface string 			`json:"interface"`
	Address string				`json:"address"`
	Port string 				`json:"port"`
	PublicKey string 			`json:"publicKey"`
	AllowedIps []string 		`json:"allowedIps"`
	PersistentKeepalive int		`json:"persistentKeepalive"`
}

type InterfaceState struct {
	Name string
	PublicKey string
}

type State struct {
	Interfaces []InterfaceState
}

// Hosts is used to query the host-related endpoints.
type Hosts struct {
	client *Client
}

// Settings returns a handle on the settings endpoints.
func (c *Client) Hosts() *Hosts {
	return &Hosts{client: c}
}

// Fetch self settings
func (h Hosts) GetSelfSettings() (*Settings, error) {
	var resp Settings
	err := h.client.Get("/hosts/self/settings", &resp)
	if err != nil {
		return nil,err
	}
	return &resp,nil
}

// Fetch self settings
func (h Hosts) PostSelfState(state *State) (error) {
		err := h.client.Post("/hosts/self/state", state, nil)
	if err != nil {
		return err
	}
	return nil
}

// sync self
func (h Hosts) PostSelfSync(state *State) (*Settings,error) {
	var resp Settings
		err := h.client.Post("/hosts/self/sync", state, &resp)
	if err != nil {
		return nil,err
	}
	return &resp,nil
}

//cli.Add("GET", "hosts/self/settings", h.GetSelfSettings)
//cli.Add("POST", "hosts/self/state", h.UpdateSelfState)
//cli.Add("POST", "hosts/self/sync", h.SynchronizeSelf)