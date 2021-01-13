package drago

import (
	"context"

	auth "github.com/seashell/drago/drago/auth"
	state "github.com/seashell/drago/drago/state"
	structs "github.com/seashell/drago/drago/structs"
	log "github.com/seashell/drago/pkg/log"
)

const (
	InterfaceList  = "list"
	InterfaceRead  = "read"
	InterfaceWrite = "write"
)

type InterfaceService struct {
	config      *Config
	logger      log.Logger
	state       state.Repository
	authHandler auth.AuthorizationHandler
}

// NewInterfaceService ...
func NewInterfaceService(config *Config, logger log.Logger, state state.Repository, authHandler auth.AuthorizationHandler) *InterfaceService {
	return &InterfaceService{
		config:      config,
		logger:      logger,
		state:       state,
		authHandler: authHandler,
	}
}

// GetInterface returns an Interface entity by ID
func (s *InterfaceService) GetInterface(args *structs.InterfaceSpecificRequest, out *structs.SingleInterfaceResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "interface", args.ID, InterfaceRead); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	n, err := s.state.InterfaceByID(ctx, args.ID)
	if err != nil {
		return structs.ErrNotFound
	}

	out.Interface = n

	return nil
}

// ListInterfaces retrieves all interface entities in the repository
func (s *InterfaceService) ListInterfaces(args *structs.InterfaceListRequest, out *structs.InterfaceListResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "interface", "", InterfaceList); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	out.Items = nil

	var err error
	var interfaces []*structs.Interface

	if args.NodeID != "" {
		if interfaces, err = s.state.InterfacesByNodeID(ctx, args.NodeID); err != nil {
			return structs.ErrInternal
		}
	} else if args.NetworkID != "" {
		if interfaces, err = s.state.InterfacesByNetworkID(ctx, args.NetworkID); err != nil {
			return structs.ErrInternal
		}
	} else {
		if interfaces, err = s.state.Interfaces(ctx); err != nil {
			return structs.ErrInternal
		}
	}

	for _, i := range interfaces {
		if args.NetworkID != "" {
			if i.NetworkID == args.NetworkID {
				out.Items = append(out.Items, i.Stub())
			}
		} else {
			out.Items = append(out.Items, i.Stub())
		}
	}

	return nil
}
