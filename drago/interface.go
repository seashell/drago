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
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "interface", args.InterfaceID, InterfaceRead); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	n, err := s.state.InterfaceByID(ctx, args.InterfaceID)
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

// UpsertInterface upserts a new Interface entity, which results in a node being added to a network
func (s *InterfaceService) UpsertInterface(args *structs.InterfaceUpsertRequest, out *structs.GenericResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "interface", "", InterfaceWrite); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	i := args.Interface

	// Make sure the input is valid
	err := i.Validate()
	if err != nil {
		return structs.NewInvalidInputError(err.Error())
	}

	// If the interface already exists, we simply merge the new values into the existing struct.
	// Otherwise, we generate a new ID and set the protected attributes in preparation for inserting
	// the new struct into the repository.
	if i.ID != "" {
		old, err := s.state.InterfaceByID(ctx, i.ID)
		if err != nil {
			return structs.ErrNotFound // interface does not exist
		}
		i = old.Merge(i)
	} else {
		i.ID = uuid.Generate()
		i.Name = nil                // Setting name is responsibility of the client node
		i.Address = nil             // TODO: set with leasing plugin if it is loaded and enabled
		i.Peers = []*structs.Peer{} // TODO: set with meshing plugin if it is loaded and enabled
		i.CreatedAt = time.Now()
	}

	// Retrieve the network to which the interface is meant to be added, throwing an error if it does not exist
	network, err := s.state.NetworkByID(ctx, i.NetworkID)
	if err != nil {
		return structs.ErrInternal // network does not exist
	}

	// Retrieve the node to which the interface is meant to be added, throwing an error if it does not exist
	node, err := s.state.NodeByID(ctx, i.NodeID)
	if err != nil {
		return structs.ErrInternal // node does not exist
	}

	// Retrieve the already existing interfaces of the targeted node
	nodeInterfaces, err := s.state.InterfacesByNodeID(ctx, node.ID)
	if err != nil {
		return structs.ErrInternal // error getting node interfaces
	}

	// Make sure that the node will not end up with two interfaces in the same network
	for _, iface := range nodeInterfaces {
		if iface.NetworkID == network.ID && i.ID != iface.ID {
			return structs.NewInternalError("Network already joined")
		}
	}

	i.UpdatedAt = time.Now()

	// TODO: wrap in a transaction

	node.UpsertInterface(i.ID)
	err = s.state.UpsertNode(ctx, node)
	if err != nil {
		return structs.ErrInternal // could not update network with the new interface
	}

	network.UpsertInterface(i.ID)
	err = s.state.UpsertNetwork(ctx, network)
	if err != nil {
		return structs.ErrInternal // could not update network with the new interface
	}

	err = s.state.UpsertInterface(ctx, i)
	if err != nil {
		return structs.ErrInternal // could not create interface
	}

	return nil
}

// DeleteInterface deletes an interface entity from the repository
func (s *InterfaceService) DeleteInterface(args *structs.InterfaceDeleteRequest, out *structs.GenericResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "interface", "", InterfaceWrite); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	for _, id := range args.InterfaceIDs {
		if iface, err := s.state.InterfaceByID(ctx, id); err == nil {

			network, err := s.state.NetworkByID(ctx, iface.NetworkID)
			if err != nil {
				return structs.NewInternalError(err.Error())
			}

			node, err := s.state.NodeByID(ctx, iface.NodeID)
			if err != nil {
				return structs.NewInternalError(err.Error())
			}

			connections, err := s.state.ConnectionsByInterfaceID(ctx, iface.ID)
			if err != nil {
				return structs.NewInternalError(err.Error())
			}

			connectionIDs := []string{}
			for _, c := range connections {
				connectionIDs = append(connectionIDs, c.ID)
				network.RemoveConnection(c.ID)
			}

			network.RemoveInterface(iface.ID)
			node.RemoveInterface(iface.ID)

			if err := s.state.DeleteConnections(ctx, connectionIDs); err != nil {
				return structs.NewInternalError(err.Error())
			}

			if err := s.state.UpsertNetwork(ctx, network); err != nil {
				return structs.NewInternalError(err.Error())
			}

			if err := s.state.UpsertNode(ctx, node); err != nil {
				return structs.NewInternalError(err.Error())
			}
		}
	}

	if err := s.state.DeleteInterfaces(ctx, args.InterfaceIDs); err != nil {
		return structs.ErrInternal
	}

	return nil
}
