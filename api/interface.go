package api

import "context"

const (
	interfacesPath = "/api/interfaces"
)

type Interface struct {
	ID         *string `json:"id,omitempty"`
	Name       *string `json:"name"`
	HostID     *string `json:"hostId"`
	NetworkID  *string `json:"networkId,omitempty"`
	IPAddress  *string `json:"ipAddress,omitempty"`
	ListenPort *string `json:"listenPort,omitempty"`
	PublicKey  *string `json:"publicKey"`
	Table      *string `json:"table"`
	DNS        *string `json:"dns"`
	MTU        *string `json:"mtu"`
	PreUp      *string `json:"preUp"`
	PostUp     *string `json:"postUp"`
	PreDown    *string `json:"preDown"`
	PostDown   *string `json:"postDown"`
}

// Interfaces is a handle to the interfaces API
type Interfaces struct {
	client *Client
}

// Interfaces returns a handle on the interfaces endpoints.
func (c *Client) Interfaces() *Interfaces {
	return &Interfaces{client: c}
}

// Get :
func (i *Interfaces) Get(ctx context.Context, id string) (*Interface, error) {
	receiver := struct {
		*Interface
	}{}

	err := i.client.getResource(id, interfacesPath, &receiver)
	if err != nil {
		return nil, err
	}

	return receiver.Interface, nil
}

// Create :
func (i *Interfaces) Create(ctx context.Context, iface *Interface) (*string, error) {
	receiver := struct {
		ID *string `json:"id"`
	}{}

	err := i.client.createResource(interfacesPath, i, &receiver)
	if err != nil {
		return nil, err
	}

	return receiver.ID, nil
}

// Update :
func (i *Interfaces) Update(ctx context.Context, iface *Interface) (*string, error) {

	receiver := struct {
		*Interface
	}{}

	err := i.client.updateResource(*iface.ID, interfacesPath, iface, &receiver)
	if err != nil {
		return nil, err
	}

	return receiver.ID, nil
}

// Delete :
func (i *Interfaces) Delete(ctx context.Context, id string) (*string, error) {

	receiver := struct {
		*Interface
	}{}

	err := i.client.deleteResource(id, interfacesPath, &receiver)
	if err != nil {
		return nil, err
	}

	return receiver.ID, nil
}

// ListByHostID :
func (i *Interfaces) ListByHostID(ctx context.Context, id string, page *Page) ([]*Interface, *Page, error) {
	receiver := struct {
		Interfaces []*Interface `json:"items"`
		*Page
	}{}

	filters := map[string]string{
		"hostId": id,
	}

	err := i.client.listResources(interfacesPath, page, filters, &receiver)
	if err != nil {
		return nil, nil, err
	}

	return receiver.Interfaces, receiver.Page, nil
}

// ListByNetworkID :
func (i *Interfaces) ListByNetworkID(ctx context.Context, id string, page *Page) ([]*Interface, *Page, error) {
	receiver := struct {
		Interfaces []*Interface `json:"items"`
		*Page
	}{}

	filters := map[string]string{
		"networkId": id,
	}

	err := i.client.listResources(interfacesPath, page, filters, &receiver)
	if err != nil {
		return nil, nil, err
	}

	return receiver.Interfaces, receiver.Page, nil
}
