package api
// NetworkInput :
type NetworkInput struct {
	Name           string `json:"name"`
	IPAddressRange string `json:"ipAddressRange"`
}

type Network struct {
	ID				string 	`json:"id,omitempty"`
	Name			string 	`json:"name,omitempty"`
	IPAddressRange	string 	`json:"ipAddressRange,omitempty"`
	//...
}

// NetworksList :
type NetworksList struct {
	Items []*Network `json:"items"`
}

// Networks is used to query the network-related endpoints.
type Networks struct {
	client *Client
}

// Networks returns a handle on the networks endpoints.
func (c *Client) Networks() *Networks {
	return &Networks{client: c}
}

// ListNetwork :
func (n *Networks) ListNetworks() (*NetworksList,error) {
	var r NetworksList
	err := n.client.Get("/networks", &r, nil)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

//CreateNetwork :
func (n *Networks) CreateNetwork(ni *NetworkInput) (*Network,error) {
	var r Network
	err := n.client.Post("/networks", ni, &r, nil)
	if err != nil {
		return nil,err
	}
	return &r,nil
}
