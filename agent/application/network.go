package application

import (
	"errors"
	"net/rpc"

	structs "github.com/seashell/drago/drago/application/structs"
)

// NetworkService
type NetworkService interface {
	Get(*structs.GetNetworkInput) (*structs.GetNetworkOutput, error)
	Create(*structs.CreateNetworkInput) (*structs.CreateNetworkOutput, error)
	Update(*structs.UpdateNetworkInput) (*structs.UpdateNetworkOutput, error)
	Delete(*structs.DeleteNetworkInput) (*structs.DeleteNetworkOutput, error)
	List(*structs.ListNetworksInput) (*structs.ListNetworksOutput, error)
}

type networkService struct {
	rpcClient *rpc.Client
}

// NewNetworkService
func NewNetworkService(client *rpc.Client) NetworkService {
	return &networkService{client}
}

// GetByID returns a Network entity by ID
func (s *networkService) Get(in *structs.GetNetworkInput) (*structs.GetNetworkOutput, error) {

	var reply structs.GetNetworkOutput

	err := s.rpcClient.Call("Networks.Get", in, &reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

// Create creates a new Network entity
func (s *networkService) Create(in *structs.CreateNetworkInput) (*structs.CreateNetworkOutput, error) {
	var reply structs.CreateNetworkOutput

	err := s.rpcClient.Call("Networks.Create", in, &reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

// Update updates an already existing Network entity
func (s *networkService) Update(in *structs.UpdateNetworkInput) (*structs.UpdateNetworkOutput, error) {
	return nil, errors.New("not implemented")
}

// Delete :
func (s *networkService) Delete(in *structs.DeleteNetworkInput) (*structs.DeleteNetworkOutput, error) {
	return nil, errors.New("not implemented")
}

// List :
func (s *networkService) List(in *structs.ListNetworksInput) (*structs.ListNetworksOutput, error) {
	return nil, errors.New("not implemented")
}
