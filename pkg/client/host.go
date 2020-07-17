package client

import "context"

// Host resource methods
func (d *DragoAPIClient) GetHost(ctx context.Context, id string) (*Host, error) {
	receiver := struct {
		*Host
	}{}

	err := d.getResource(
		id,
		hostsPath,
		&receiver,
	)
	if err != nil {
		return nil, err
	}

	return receiver.Host, nil
}

func (d *DragoAPIClient) AddHost(ctx context.Context, h *Host) (*string, error) {
	receiver := struct {
		ID *string `json:"id"`
	}{}

	err := d.createResource(
		hostsPath,
		h,
		&receiver,
	)
	if err != nil {
		return nil, err
	}

	return receiver.ID, nil
}

func (d *DragoAPIClient) UpdateHost(ctx context.Context, h *Host) (*string, error) {
	receiver := struct {
		ID *string `json:"id"`
	}{}

	err := d.updateResource(
		*h.ID,
		hostsPath,
		h,
		&receiver,
	)
	if err != nil {
		return nil, err
	}

	return receiver.ID, nil
}

func (d *DragoAPIClient) DeleteHost(ctx context.Context, id string) error {
	err := d.deleteResource(
		id,
		hostsPath,
	)
	if err != nil {
		return err
	}

	return nil
}

func (d *DragoAPIClient) ListHostsByLabel(
	ctx context.Context,
	label string,
	page *Page,
) ([]*Host, *Page, error) {
	receiver := struct {
		Hosts []*Host `json:"items"`
		*Page
	}{}

	err := d.listResources(
		hostsPath,
		page,
		map[string]string{labelQueryParamKey: label},
		&receiver,
	)
	if err != nil {
		return nil, nil, err
	}

	return receiver.Hosts, receiver.Page, nil
}
