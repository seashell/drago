package api

import "context"

const (
	hostsPath = "/api/hosts"
)

type Host struct {
	ID               *string  `json:"id,omitempty"`
	Name             *string  `json:"name"`
	AdvertiseAddress *string  `json:"advertiseAddress,omitempty"`
	Labels           []string `json:"labels"`
}

// Hosts is a handle to the hosts API
type Hosts struct {
	client *Client
}

// Hosts returns a handle on the hosts endpoints.
func (c *Client) Hosts() *Hosts {
	return &Hosts{client: c}
}

// Get :
func (h *Hosts) Get(ctx context.Context, id string) (*Host, error) {
	receiver := struct {
		*Host
	}{}

	err := h.client.getResource(hostsPath, id, &receiver)
	if err != nil {
		return nil, err
	}

	return receiver.Host, nil
}

// Create :
func (h *Hosts) Create(ctx context.Context, host *Host) (*string, error) {
	receiver := struct {
		ID *string `json:"id"`
	}{}

	err := h.client.createResource(hostsPath, host, &receiver)
	if err != nil {
		return nil, err
	}

	return receiver.ID, nil
}

// Update :
func (h *Hosts) Update(ctx context.Context, host *Host) (*string, error) {
	receiver := struct {
		*Host
	}{}

	err := h.client.updateResource(hostsPath, *host.ID, host, &receiver)
	if err != nil {
		return nil, err
	}

	return receiver.ID, nil
}

// Delete :
func (h *Hosts) Delete(ctx context.Context, id string) (*string, error) {
	receiver := struct {
		*Host
	}{}

	err := h.client.deleteResource(id, hostsPath, &receiver)
	if err != nil {
		return nil, err
	}

	return receiver.ID, nil
}

// List :
func (h *Hosts) List(ctx context.Context, label string, page *Page) ([]*Host, *Page, error) {
	receiver := struct {
		Hosts []*Host `json:"items"`
		*Page
	}{}

	filters := map[string]string{
		"label": label,
	}

	err := h.client.listResources(hostsPath, page, filters, &receiver)
	if err != nil {
		return nil, nil, err
	}

	return receiver.Hosts, receiver.Page, nil
}
