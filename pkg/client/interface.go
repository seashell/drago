package client

import "context"

// Interface resource methods
func (d *DragoAPIClient) GetInterface(ctx context.Context, id string) (*Interface, error) {
	receiver := struct {
		*Interface
	}{}

	err := d.getResource(
		id,
		interfacesPath,
		&receiver,
	)
	if err != nil {
		return nil, err
	}

	return receiver.Interface, nil
}

func (d *DragoAPIClient) AddInterface(ctx context.Context, i *Interface) (*string, error) {
	receiver := struct {
		ID *string `json:"id"`
	}{}

	err := d.createResource(
		interfacesPath,
		i,
		&receiver,
	)
	if err != nil {
		return nil, err
	}

	return receiver.ID, nil
}

func (d *DragoAPIClient) UpdateInterface(ctx context.Context, i *Interface) (*string, error) {
	receiver := struct {
		ID *string `json:"id"`
	}{}

	err := d.updateResource(
		*i.ID,
		interfacesPath,
		i,
		&receiver,
	)
	if err != nil {
		return nil, err
	}

	return receiver.ID, nil
}

func (d *DragoAPIClient) DeleteInteface(ctx context.Context, id string) error {
	err := d.deleteResource(
		id,
		interfacesPath,
	)
	if err != nil {
		return err
	}

	return nil
}

func (d *DragoAPIClient) ListInterfacesByHostID(
	ctx context.Context,
	id string,
	page *Page,
) ([]*Interface, *Page, error) {
	receiver := struct {
		Interfaces []*Interface `json:"items"`
		*Page
	}{}

	err := d.listResources(
		interfacesPath,
		page,
		map[string]string{hostIDQueryParamKey: id},
		&receiver,
	)
	if err != nil {
		return nil, nil, err
	}

	return receiver.Interfaces, receiver.Page, nil
}
