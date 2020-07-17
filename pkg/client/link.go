package client

import "context"

// Link resource methods
func (d *DragoAPIClient) GetLink(ctx context.Context, id string) (*Link, error) {
	receiver := struct {
		*Link
	}{}

	err := d.getResource(
		id,
		linksPath,
		&receiver,
	)
	if err != nil {
		return nil, err
	}

	return receiver.Link, nil
}

func (d *DragoAPIClient) AddLink(ctx context.Context, l *Link) (*string, error) {
	receiver := struct {
		ID *string `json:"id"`
	}{}

	err := d.createResource(
		linksPath,
		l,
		&receiver,
	)
	if err != nil {
		return nil, err
	}

	return receiver.ID, nil
}

func (d *DragoAPIClient) UpdateLink(ctx context.Context, l *Link) (*string, error) {
	receiver := struct {
		ID *string `json:"id"`
	}{}

	err := d.updateResource(
		*l.ID,
		linksPath,
		l,
		&receiver,
	)
	if err != nil {
		return nil, err
	}

	return receiver.ID, nil
}

func (d *DragoAPIClient) DeleteLink(ctx context.Context, id string) error {
	err := d.deleteResource(
		id,
		linksPath,
	)
	if err != nil {
		return err
	}

	return nil
}
