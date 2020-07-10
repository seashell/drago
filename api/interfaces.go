package api

// InterfaceInput :
type InterfaceInput struct {
	Name       	string	`json:"name"`
	HostID     	string 	`json:"hostId"`
	NetworkID  	string 	`json:"networkId"`
	IPAddress  	string 	`json:"ipAddress"`
	ListenPort 	string 	`json:"listenPort,omitempty"`
	//...
}

// Interface :
type Interface struct {	
	ID			string 	`json:"id,omitempty"`
	Name       	string	`json:"name,omitempty"`
	HostID     	string	`json:"hostId,omitempty"`
	NetworkID  	string	`json:"networkId,omitempty"`
	IPAddress 	string 	`json:"ipAddress,omitempty"`
	ListenPort	string 	`json:"listenPort,omitempty"`
	//...
}


// ListInterfacesInput :
type ListInterfacesInput struct {
	HostIDFilter    string `url:"hostId"`
	NetworkIDFilter string `url:"networkId"`
}

// InterfacesList :
type InterfacesList struct {
	Items []*Interface `json:"items"`
}

// Interfaces is used to query the interface-related endpoints.
type Interfaces struct {
	client *Client
}

// Interfaces returns a handle on the networks endpoints.
func (c *Client) Interfaces() *Interfaces {
	return &Interfaces{client: c}
}

// ListInterfaces :
func (i *Interfaces) ListInterfaces(q ListInterfacesInput) (*InterfacesList,error) {
	var r InterfacesList
	err := i.client.Get("/interfaces", &r, &q)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// CreateInterface :
func (i *Interfaces) CreateInterface(ni *InterfaceInput) (*Interface,error) {
	var r Interface
	err := i.client.Post("/interfaces", ni, &r, nil)
	if err != nil {
		return nil, err
	}
	return &r, nil
}