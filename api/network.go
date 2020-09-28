package api

import "context"

const (
	networksPath = "/api/networks"
)

type Network struct {
	ID             *string `json:"id,omitempty"`
	Name           *string `json:"name,omitempty"`
	IPAddressRange *string `json:"ipAddressRange,omitempty"`
}

// Networks is a handle to the networks API
type Networks struct {
	client *Client
}

// Networks returns a handle on the networks endpoints.
func (c *Client) Networks() *Networks {
	return &Networks{client: c}
}

// Get :
func (n *Networks) Get(ctx context.Context, id string) (*Network, error) {
	receiver := struct {
		*Network
	}{}

	err := n.client.getResource(id, networksPath, &receiver)
	if err != nil {
		return nil, err
	}

	return receiver.Network, nil
}

// Create :
func (n *Networks) Create(ctx context.Context, network *Network) (*string, error) {
	receiver := struct {
		ID *string `json:"id"`
	}{}

	err := n.client.createResource(networksPath, network, &receiver)
	if err != nil {
		return nil, err
	}

	return receiver.ID, nil
}

// Update :
func (n *Networks) Update(ctx context.Context, network *Network) (*string, error) {

	receiver := struct {
		*Network
	}{}

	err := n.client.updateResource(*network.ID, networksPath, network, &receiver)
	if err != nil {
		return nil, err
	}

	return receiver.ID, nil
}

// Delete :
func (n *Networks) Delete(ctx context.Context, id string) (*string, error) {
	receiver := struct {
		*Network
	}{}

	err := n.client.deleteResource(id, networksPath, &receiver)
	if err != nil {
		return nil, err
	}

	return receiver.ID, nil
}

// List :
func (n *Networks) List(ctx context.Context, page *Page) ([]*Network, *Page, error) {
	receiver := struct {
		Networks []*Network `json:"items"`
		*Page
	}{}

	filters := map[string]string{}

	err := n.client.listResources(networksPath, page, filters, &receiver)
	if err != nil {
		return nil, nil, err
	}

	return receiver.Networks, receiver.Page, nil
}
