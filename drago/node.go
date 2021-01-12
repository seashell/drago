package drago

import (
	"context"
	"time"

	auth "github.com/seashell/drago/drago/auth"
	state "github.com/seashell/drago/drago/state"
	structs "github.com/seashell/drago/drago/structs"
	uuid "github.com/seashell/drago/pkg/uuid"
)

const (
	NodeList  = "list"
	NodeRead  = "read"
	NodeWrite = "write"
)

type NodeService struct {
	config      *Config
	state       state.Repository
	authHandler auth.AuthorizationHandler
}

// NewNodeService ...
func NewNodeService(config *Config, state state.Repository, authHandler auth.AuthorizationHandler) *NodeService {
	return &NodeService{
		config:      config,
		state:       state,
		authHandler: authHandler,
	}
}

func (s *NodeService) Register(args *structs.NodeRegisterRequest, reply *structs.NodeUpdateResponse) error {

	ctx := context.TODO()

	err := args.Validate()
	if err != nil {
		return structs.ErrInvalidInput
	}

	n := args.Node

	if n.Status == "" {
		n.Status = structs.NodeStatusInit
	}

	err = n.Validate()
	if err != nil {
		return structs.ErrInvalidInput
	}

	old, err := s.state.NodeByID(ctx, n.ID)
	if err != nil {
		n.CreatedAt = time.Now()
	} else {

		// Check if the SecretID has been tampered with
		if old != nil {
			if args.Node.SecretID != old.SecretID {
				return structs.ErrInvalidInput // node secret ID does not match
			}
		}

		n = old.Merge(n)
	}

	n.UpdatedAt = time.Now()

	err = s.state.UpsertNode(ctx, n)
	if err != nil {
		return structs.ErrInternal
	}

	return nil
}

func (s *NodeService) GetInterfaces(args *structs.NodeSpecificRequest, reply *structs.NodeInterfacesResponse) error {
	return nil
}

func (s *NodeService) GetPeers(args *structs.NodeSpecificRequest, reply *structs.NodePeersResponse) error {
	return nil
}

func (s *NodeService) UpdateInterfaces(args *structs.InterfaceUpdateRequest, reply *structs.GenericResponse) error {
	return nil
}

func (s *NodeService) UpdatePeers(args *structs.PeerUpdateRequest, reply *structs.GenericResponse) error {
	return nil
}

// GetNode returns a Node entity by ID
func (s *NodeService) GetNode(args *structs.NodeSpecificRequest, out *structs.SingleNodeResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "network", args.ID, NodeRead); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	n, err := s.state.NodeByID(ctx, args.ID)
	if err != nil {
		return structs.ErrNotFound
	}

	out.Node = n

	return nil
}

// ListNodes retrieves all network entities in the repository
func (s *NodeService) ListNodes(args *structs.NodeListRequest, out *structs.NodeListResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "network", "", NodeList); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	networks, err := s.state.Nodes(ctx)
	if err != nil {
		return structs.ErrInternal
	}

	out.Items = nil

	for _, n := range networks {
		out.Items = append(out.Items, n.Stub())
	}

	return nil
}

// UpsertNode upserts a new Node entity
func (s *NodeService) UpsertNode(args *structs.NodeUpsertRequest, out *structs.GenericResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "network", "", NodeWrite); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	n := args.Node

	err := n.Validate()
	if err != nil {
		return structs.ErrInvalidInput
	}

	if n.ID == "" {
		n.ID = uuid.Generate()
		n.CreatedAt = time.Now()
	} else {
		old, err := s.state.NodeByID(ctx, n.ID)
		if err != nil {
			return structs.ErrNotFound
		}
		n = old.Merge(n)
	}

	n.UpdatedAt = time.Now()

	err = s.state.UpsertNode(ctx, n)
	if err != nil {
		return structs.ErrInternal
	}

	return nil
}

// DeleteNode deletes a network entity from the repository
func (s *NodeService) DeleteNode(args *structs.NodeDeleteRequest, out *structs.GenericResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "network", "", NodeWrite); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	err := s.state.DeleteNodes(ctx, args.IDs)
	if err != nil {
		return structs.ErrInternal
	}

	return nil
}
