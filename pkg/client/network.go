package client

import "context"

// Network resource methods
func (d *DragoAPIClient) GetNetwork(ctx context.Context, id string) (*Network, error) {
	receiver := struct {
		*Network
	}{}

	err := d.getResource(
		id,
		networksPath,
		&receiver,
	)
	if err != nil {
		return nil, err
	}

	return receiver.Network, nil
}

func (d *DragoAPIClient) AddNetwork(ctx context.Context, n *Network) (*string, error) {
	receiver := struct {
		ID *string `json:"id"`
	}{}

	err := d.createResource(
		networksPath,
		n,
		&receiver,
	)
	if err != nil {
		return nil, err
	}

	return receiver.ID, nil
}

func (d *DragoAPIClient) UpdateNetwork(ctx context.Context, n *Network) (*string, error) {
	receiver := struct {
		ID *string `json:"id"`
	}{}

	err := d.updateResource(
		*n.ID,
		networksPath,
		n,
		&receiver,
	)
	if err != nil {
		return nil, err
	}

	return receiver.ID, nil
}

func (d *DragoAPIClient) DeleteNetwork(ctx context.Context, id string) error {
	err := d.deleteResource(
		id,
		networksPath,
	)
	if err != nil {
		return err
	}

	return nil
}

func (d *DragoAPIClient) ListNetworks(
	ctx context.Context,
	page *Page,
) ([]*Network, *Page, error) {
	s := struct {
		N []*Network `json:"items"`
		*Page
	}{}

	err := d.listResources(
		networksPath,
		page,
		nil,
		&s,
	)
	if err != nil {
		return nil, nil, err
	}

	return s.N, s.Page, nil
}
