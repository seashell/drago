package drago

import (
	"context"
	"fmt"
	"sync"
	"time"

	auth "github.com/seashell/drago/drago/auth"
	state "github.com/seashell/drago/drago/state"
	structs "github.com/seashell/drago/drago/structs"
	log "github.com/seashell/drago/pkg/log"
)

const (
	NodeList  = "list"
	NodeRead  = "read"
	NodeWrite = "write"

	defaultExpectedHeartbeatInterval = 3 * time.Second
)

type NodeService struct {
	config              *Config
	logger              log.Logger
	state               state.Repository
	authHandler         auth.AuthorizationHandler
	heartbeatTimers     map[string]*time.Timer
	heartbeatTimersLock sync.Mutex
}

// NewNodeService ...
func NewNodeService(config *Config, logger log.Logger, state state.Repository, authHandler auth.AuthorizationHandler) (*NodeService, error) {

	s := &NodeService{
		config:          config,
		state:           state,
		authHandler:     authHandler,
		logger:          logger,
		heartbeatTimers: map[string]*time.Timer{},
	}

	err := s.setupHeartbeatTimers()
	if err != nil {
		return nil, fmt.Errorf("error setting up heartbeat timers: %v", err)
	}

	return s, nil
}

func (s *NodeService) setupHeartbeatTimers() error {

	ctx := context.TODO()

	nodes, err := s.state.Nodes(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve nodes from state: %v", err)
	}

	for _, n := range nodes {
		s.resetHeartbeatTimer(n.ID)
	}

	return nil
}

func (s *NodeService) resetHeartbeatTimer(id string) {

	s.heartbeatTimersLock.Lock()
	defer s.heartbeatTimersLock.Unlock()

	if timer, ok := s.heartbeatTimers[id]; ok {
		timer.Reset(defaultExpectedHeartbeatInterval)
		return
	}

	timer := time.AfterFunc(defaultExpectedHeartbeatInterval, func() {

		ctx := context.TODO()

		s.logger.Debugf("heartbeat missed by node %s", id)

		old, err := s.state.NodeByID(ctx, id)
		if err != nil {
			s.logger.Debugf("failed to set node status after heartbeat miss: %v", err)
		}

		n := old.Merge(&structs.Node{
			ID:     id,
			Status: structs.NodeStatusDown,
		})
		n.UpdatedAt = time.Now()

		if err := s.state.UpsertNode(ctx, n); err != nil {
			s.logger.Debugf("failed to set node status after hearbeat miss: %v", err)
		}
	})

	s.heartbeatTimers[id] = timer
}

func (s *NodeService) Register(args *structs.NodeRegisterRequest, out *structs.NodeUpdateResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "node", "", NodeWrite); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	err := args.Validate()
	if err != nil {
		return structs.ErrInvalidInput
	}

	n := args.Node

	if n.Status == "" {
		n.Status = structs.NodeStatusInit
	}

	if !structs.IsValidNodeStatus(n.Status) {
		return structs.ErrInvalidInput
	}

	old, err := s.state.NodeByID(ctx, n.ID)
	if err != nil {
		s.logger.Debugf("registering a new node with id %s!", n.ID)
		n.CreatedAt = time.Now()
	} else {
		s.logger.Debugf("node %s already registered.", n.ID)
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

	s.resetHeartbeatTimer(n.ID)

	return nil
}

func (s *NodeService) UpdateStatus(args *structs.NodeUpdateStatusRequest, out *structs.NodeUpdateResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "node", args.ID, NodeWrite); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	if args.ID == "" {
		return structs.ErrInvalidInput
	}
	if !structs.IsValidNodeStatus(args.Status) {
		return structs.ErrInvalidInput
	}

	n, err := s.state.NodeByID(ctx, args.ID)
	if err != nil {
		return structs.ErrNotFound
	}

	n.Status = args.Status
	n.UpdatedAt = time.Now()

	err = s.state.UpsertNode(ctx, n)
	if err != nil {
		return structs.ErrInternal
	}

	out.Servers = []string{s.config.RPCAdvertiseAddr}

	s.logger.Debugf("heartbeat from node %s", n.ID)
	s.resetHeartbeatTimer(n.ID)

	return nil
}

func (s *NodeService) GetInterfaces(args *structs.NodeSpecificRequest, out *structs.NodeInterfacesResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "node", args.ID, NodeRead); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	if args.ID == "" {
		return structs.ErrInvalidInput
	}

	interfaces, err := s.state.InterfacesByNodeID(ctx, args.ID)
	if err != nil {
		return structs.ErrNotFound
	}

	out.Items = nil

	for _, i := range interfaces {
		out.Items = append(out.Items, i)
	}

	return nil
}

func (s *NodeService) GetPeers(args *structs.NodeSpecificRequest, out *structs.NodePeersResponse) error {
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
