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
	uuid "github.com/seashell/drago/pkg/uuid"
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
		n.ModifyIndex = 0
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
	n.ModifyIndex++

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
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "node", args.NodeID, NodeWrite); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	if args.NodeID == "" {
		return structs.ErrInvalidInput
	}
	if !structs.IsValidNodeStatus(args.Status) {
		return structs.ErrInvalidInput
	}

	n, err := s.state.NodeByID(ctx, args.NodeID)
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
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "node", args.NodeID, NodeRead); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	if args.NodeID == "" {
		return structs.ErrInvalidInput
	}

	interfaces, err := s.state.InterfacesByNodeID(ctx, args.NodeID)
	if err != nil {
		return structs.ErrNotFound
	}

	out.Items = nil

	for _, i := range interfaces {
		out.Items = append(out.Items, i)
	}

	return nil
}

func (s *NodeService) UpdateInterfaces(args *structs.NodeInterfaceUpdateRequest, out *structs.GenericResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "node", "", NodeWrite); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	node, err := s.state.NodeByID(ctx, args.NodeID)
	if err != nil {
		return structs.ErrNotFound
	}

	nodeInterfaces, err := s.state.InterfacesByNodeID(ctx, node.ID)
	if err != nil {
		return structs.ErrInternal
	}

	// Create a map for more efficient lookup
	nodeInterfacesMap := map[string]*structs.Interface{}
	for _, i := range nodeInterfaces {
		nodeInterfacesMap[i.ID] = i
	}

	for _, i := range args.Interfaces {
		old, found := nodeInterfacesMap[i.ID]
		if !found {
			return structs.ErrNotFound // Node does not contain interface being updated
		}

		i = old.Merge(i)
		i.UpdatedAt = time.Now()

		err := s.state.UpsertInterface(ctx, i)
		if err != nil {
			return structs.ErrInternal // Error updating interface
		}
	}

	return structs.ErrInvalidInput
}

// GetNode returns a Node entity by ID
func (s *NodeService) GetNode(args *structs.NodeSpecificRequest, out *structs.SingleNodeResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "network", args.NodeID, NodeRead); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	n, err := s.state.NodeByID(ctx, args.NodeID)
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

// JoinNetwork : connects a node to a network
func (s *NetworkService) JoinNetwork(args *structs.NodeJoinNetworkRequest, out *structs.GenericResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "node", args.NodeID, NodeList); err != nil {
			return structs.ErrPermissionDenied
		}
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "network", args.NetworkID, NodeList); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	network, err := s.state.NetworkByID(ctx, args.NodeID)
	if err != nil {
		return structs.ErrNotFound // network not found
	}

	node, err := s.state.NodeByID(ctx, args.NodeID)
	if err != nil {
		return structs.ErrNotFound // node not found
	}

	interfaces, err := s.state.InterfacesByNodeID(ctx, node.ID)
	if err != nil {
		return structs.ErrInternal // error getting node interfaces
	}

	// Check whether node has already joined the network
	for _, iface := range interfaces {
		if iface.NetworkID == network.ID {
			return fmt.Errorf("Network already joined")
		}
	}

	iface := &structs.Interface{
		ID:          uuid.Generate(),
		NodeID:      node.ID,
		NetworkID:   network.ID,
		Address:     "",                // to be set if leasing plugin is loaded and enabled
		Peers:       []*structs.Peer{}, // to be set if meshing plugin is loaded and enabled
		ModifyIndex: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.state.UpsertInterface(ctx, iface)
	if err != nil {
		return structs.ErrInternal // could not create interface
	}

	return nil
}

// LeaveNetwork : disconnects a node from a network
func (s *NetworkService) LeaveNetwork(args *structs.NodeLeaveNetworkRequest, out *structs.GenericResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "node", args.NodeID, NodeList); err != nil {
			return structs.ErrPermissionDenied
		}
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "network", args.NetworkID, NodeList); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	network, err := s.state.NetworkByID(ctx, args.NodeID)
	if err != nil {
		return structs.ErrNotFound // network not found
	}

	node, err := s.state.NodeByID(ctx, args.NodeID)
	if err != nil {
		return structs.ErrNotFound // node not found
	}

	interfaces, err := s.state.InterfacesByNodeID(ctx, node.ID)
	if err != nil {
		return structs.ErrInternal // error getting node interfaces
	}

	// Check whether node is in the network
	for _, iface := range interfaces {
		if iface.NetworkID == network.ID {
			if err := s.state.DeleteInterfaces(ctx, []string{iface.ID}); err != nil {
				return structs.ErrInternal
			}
		}
	}

	return nil
}
