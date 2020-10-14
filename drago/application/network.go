package application

import (
	"context"
	"errors"

	structs "github.com/seashell/drago/drago/application/structs"
	domain "github.com/seashell/drago/drago/domain"
)

const (
	NetworkList  = "list"
	NetworkRead  = "read"
	NetworkWrite = "write"
)

const (
	errNetworkNotFound = "network not found"
)

var (
	// ErrNetworkNotFound ...
	ErrNetworkNotFound = errors.New(errNetworkNotFound)
)

type networkService struct {
	config *Config
}

// NewNetworkService ...
func NewNetworkService(config *Config) NetworkService {
	return &networkService{config}
}

// GetByID returns a Network entity by ID
func (s *networkService) GetByID(ctx context.Context, in *structs.NetworkGetInput) (*structs.NetworkGetOutput, error) {

	// Check if authorized
	if err := s.authorize(ctx, in.Subject, in.ID, NetworkRead); err != nil {
		return nil, err
	}

	n, err := s.config.NetworkRepository.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	out := &structs.NetworkGetOutput{
		Network: structs.Network{
			ID:           n.ID,
			Name:         n.Name,
			AddressRange: n.AddressRange,
			CreatedAt:    n.CreatedAt,
			UpdatedAt:    n.UpdatedAt,
		},
	}

	return out, nil
}

// Create creates a new Network entity
func (s *networkService) Create(ctx context.Context, in *structs.NetworkCreateInput) (*structs.NetworkCreateOutput, error) {

	// Check if authorized
	if err := s.authorize(ctx, in.Subject, "", NetworkWrite); err != nil {
		return nil, ErrUnauthorized
	}

	id, err := s.config.NetworkRepository.Create(ctx, &domain.Network{
		Name:         in.Name,
		AddressRange: in.AddressRange,
	})
	if err != nil {
		return nil, err
	}

	n, err := s.config.NetworkRepository.GetByID(ctx, *id)
	if err != nil {
		return nil, err
	}

	out := &structs.NetworkCreateOutput{
		ID: n.ID,
	}

	return out, nil
}

// Delete deletes a network entity from the repository
func (s *networkService) Delete(ctx context.Context, in *structs.NetworkDeleteInput) (*structs.NetworkDeleteOutput, error) {

	// Check if authorized
	if err := s.authorize(ctx, in.Subject, "", NetworkWrite); err != nil {
		return nil, ErrUnauthorized
	}

	_, err := s.config.NetworkRepository.DeleteByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	out := &structs.NetworkDeleteOutput{}
	return out, nil
}

// List retrieves all network entities in the repository
func (s *networkService) List(ctx context.Context, in *structs.NetworkListInput) (*structs.NetworkListOutput, error) {

	// Check if authorized
	if err := s.authorize(ctx, in.Subject, "", NetworkList); err != nil {
		return nil, ErrUnauthorized
	}

	networks, err := s.config.NetworkRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	out := &structs.NetworkListOutput{
		Items: []*structs.NetworkListItem{},
	}

	for _, n := range networks {
		out.Items = append(out.Items, &structs.NetworkListItem{
			ID:           n.ID,
			Name:         n.Name,
			AddressRange: n.AddressRange,
			CreatedAt:    n.CreatedAt,
			UpdatedAt:    n.UpdatedAt,
		})
	}

	return out, nil
}

func (s *networkService) authorize(ctx context.Context, sub, id, op string) error {
	if s.config.ACLEnabled {
		if err := s.config.AuthHandler.Authorize(ctx, sub, ResourceNetwork, id, op); err != nil {
			return err
		}
	}
	return nil
}
