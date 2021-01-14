package drago

import (
	"context"
	"time"

	auth "github.com/seashell/drago/drago/auth"
	state "github.com/seashell/drago/drago/state"
	structs "github.com/seashell/drago/drago/structs"
	log "github.com/seashell/drago/pkg/log"
	uuid "github.com/seashell/drago/pkg/uuid"
)

const (
	NetworkList  = "list"
	NetworkRead  = "read"
	NetworkWrite = "write"
)

type NetworkService struct {
	config      *Config
	logger      log.Logger
	state       state.Repository
	authHandler auth.AuthorizationHandler
}

// NewNetworkService ...
func NewNetworkService(config *Config, logger log.Logger, state state.Repository, authHandler auth.AuthorizationHandler) *NetworkService {
	return &NetworkService{
		config:      config,
		logger:      logger,
		state:       state,
		authHandler: authHandler,
	}
}

// GetNetwork returns a Network entity by ID
func (s *NetworkService) GetNetwork(args *structs.NetworkSpecificRequest, out *structs.SingleNetworkResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "network", args.ID, NetworkRead); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	n, err := s.state.NetworkByID(ctx, args.ID)
	if err != nil {
		return structs.ErrNotFound
	}

	out.Network = n

	return nil
}

// ListNetworks retrieves all network entities in the repository
func (s *NetworkService) ListNetworks(args *structs.NetworkListRequest, out *structs.NetworkListResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "network", "", NetworkList); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	networks, err := s.state.Networks(ctx)
	if err != nil {
		return structs.ErrInternal
	}

	out.Items = nil

	for _, n := range networks {
		out.Items = append(out.Items, n.Stub())
	}

	return nil
}

// UpsertNetwork upserts a new Network entity
func (s *NetworkService) UpsertNetwork(args *structs.NetworkUpsertRequest, out *structs.GenericResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "network", "", NetworkWrite); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	n := args.Network

	err := n.Validate()
	if err != nil {
		return structs.ErrInvalidInput
	}

	if n.ID == "" {
		n.ID = uuid.Generate()
		n.CreatedAt = time.Now()
		n.ModifyIndex = 0
	} else {
		old, err := s.state.NetworkByID(ctx, n.ID)
		if err != nil {
			return structs.ErrNotFound
		}
		n = old.Merge(n)
	}

	n.UpdatedAt = time.Now()
	n.ModifyIndex++

	err = s.state.UpsertNetwork(ctx, n)
	if err != nil {
		return structs.ErrInternal
	}

	return nil
}

// DeleteNetwork deletes a network entity from the repository
func (s *NetworkService) DeleteNetwork(args *structs.NetworkDeleteRequest, out *structs.GenericResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "network", "", NetworkWrite); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	err := s.state.DeleteNetworks(ctx, args.IDs)
	if err != nil {
		return structs.ErrInternal
	}

	return nil
}
