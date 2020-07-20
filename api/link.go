package api

import "context"

const (
	linksPath = "/api/networks"
)

type Link struct {
	ID                  *string  `json:"id,omitempty"`
	FromInterfaceID     *string  `json:"fromInterfaceId"`
	ToInterfaceID       *string  `json:"toInterfaceId"`
	AllowedIPs          []string `json:"allowedIPs"`
	PersistentKeepalive *int     `json:"persistentKeepalive,omitempty"`
}

// Links is a handle to the links API
type Links struct {
	client *Client
}

// Links returns a handle on the links endpoints.
func (c *Client) Links() *Links {
	return &Links{client: c}
}

// Get :
func (l *Links) Get(ctx context.Context, id string) (*Link, error) {
	receiver := struct {
		*Link
	}{}

	err := l.client.getResource(id, linksPath, &receiver)
	if err != nil {
		return nil, err
	}

	return receiver.Link, nil
}

// Create :
func (l *Links) Create(ctx context.Context, link *Link) (*string, error) {
	receiver := struct {
		ID *string `json:"id"`
	}{}

	err := l.client.createResource(linksPath, l, &receiver)
	if err != nil {
		return nil, err
	}

	return receiver.ID, nil
}

// Update :
func (l *Links) Update(ctx context.Context, link *Link) error {
	err := l.client.updateResource(*link.ID, linksPath, link)
	if err != nil {
		return err
	}
	return nil
}

// Delete :
func (l *Links) Delete(ctx context.Context, id string) error {
	err := l.client.deleteResource(id, linksPath)
	if err != nil {
		return err
	}
	return nil
}
