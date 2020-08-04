package receiver

import (
	application "github.com/seashell/drago/drago/application"
	structs "github.com/seashell/drago/drago/application/structs"
)

// NetworkReceiverAdapter
type NetworkReceiverAdapter struct {
	service application.NetworkService
}

// NewNetworkReceiverAdapter
func NewNetworkReceiverAdapter(service application.NetworkService) *NetworkReceiverAdapter {
	return &NetworkReceiverAdapter{
		service: service,
	}
}

// Get
func (r *NetworkReceiverAdapter) Get(in *structs.GetNetworkInput, out *structs.GetNetworkOutput) error {
	out, err := r.service.Get(in)
	if err != nil {
		return err
	}
	return nil
}

// Create
func (r *NetworkReceiverAdapter) Create(in *structs.CreateNetworkInput, out *structs.CreateNetworkOutput) error {
	out, err := r.service.Create(in)
	if err != nil {
		return err
	}
	return nil
}

// Update
func (r *NetworkReceiverAdapter) Update(in *structs.UpdateNetworkInput, out *structs.UpdateNetworkOutput) error {
	out, err := r.service.Update(in)
	if err != nil {
		return err
	}
	return nil
}

// Delete
func (r *NetworkReceiverAdapter) Delete(in *structs.DeleteNetworkInput, out *structs.DeleteNetworkOutput) error {
	out, err := r.service.Delete(in)
	if err != nil {
		return err
	}
	return nil
}

// List
func (r *NetworkReceiverAdapter) List(in *structs.ListNetworksInput, out *structs.ListNetworksOutput) error {
	out, err := r.service.List(in)
	if err != nil {
		return err
	}
	return nil
}
