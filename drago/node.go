package drago

import (
	"context"
	"fmt"
	"path"
	"strings"
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

	defaultExpectedHeartbeatInterval = 10 * time.Second
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
		return structs.NewInvalidInputError(err.Error())
	}

	n := args.Node

	if n.Status == "" {
		n.Status = structs.NodeStatusInit
	}

	if !structs.IsValidNodeStatus(n.Status) {
		return structs.NewInvalidInputError(err.Error())
	}

	old, err := s.state.NodeByID(ctx, n.ID)
	if err != nil {
		s.logger.Debugf("registering a new node with id %s!", n.ID)
		n.CreatedAt = time.Now()
	} else {
		s.logger.Debugf("node %s already registered.", n.ID)
		if old != nil {
			if args.Node.SecretID != old.SecretID {
				return structs.NewInvalidInputError("Node secret does not match")
			}
		}
		n = old.Merge(n)
	}

	n.UpdatedAt = time.Now()

	err = s.state.UpsertNode(ctx, n)
	if err != nil {
		return structs.NewInternalError(err.Error())
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
		return structs.NewInvalidInputError("Invalid node status")
	}

	n, err := s.state.NodeByID(ctx, args.NodeID)
	if err != nil {
		return structs.NewInternalError(err.Error())
	}

	n.Status = args.Status
	n.AdvertiseAddress = args.AdvertiseAddress

	if args.Meta != nil {
		n.Meta = args.Meta
	}

	n.UpdatedAt = time.Now()

	err = s.state.UpsertNode(ctx, n)
	if err != nil {
		return structs.NewInternalError(err.Error())
	}

	out.Servers = []string{s.config.RPCAdvertiseAddr}

	s.logger.Debugf("heartbeat from node %s", n.ID)
	s.resetHeartbeatTimer(n.ID)

	return nil
}

// GetInterfaces :
func (s *NodeService) GetInterfaces(args *structs.NodeSpecificRequest, out *structs.NodeInterfacesResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "node", args.NodeID, NodeRead); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	if args.NodeID == "" {
		return structs.NewInvalidInputError("Missing NodeID")
	}

	interfaces, err := s.state.InterfacesByNodeID(ctx, args.NodeID)
	if err != nil {
		return structs.ErrNotFound
	}

	for _, iface := range interfaces {

		iface.Peers = []*structs.Peer{}

		connections, err := s.state.ConnectionsByInterfaceID(ctx, iface.ID)
		if err != nil {
			s.logger.Warnf("couldn't get connections for interface %s", iface.ID)
		}

		for _, conn := range connections {

			ifaceSettings := conn.PeerSettingsByInterfaceID(iface.ID)
			peerSettings := conn.OtherPeerSettingsByInterfaceID(iface.ID)

			if ifaceSettings == nil || peerSettings == nil {
				s.logger.Warnf("couldn't get settings for connection")
			}

			peerIface, err := s.state.InterfaceByID(ctx, peerSettings.InterfaceID)
			if err != nil {
				s.logger.Warnf("couldn't get peer interface %s", peerSettings.InterfaceID)
			}

			peerNode, err := s.state.NodeByID(ctx, peerIface.NodeID)
			if err != nil {
				s.logger.Warnf("couldn't get peer node %s", peerIface.NodeID)
			}

			peer := &structs.Peer{
				PublicKey:           peerIface.PublicKey,
				Address:             &peerNode.AdvertiseAddress,
				Port:                peerIface.ListenPort,
				AllowedIPs:          []string{},
				PersistentKeepalive: conn.PersistentKeepalive,
			}

			if ifaceSettings.RoutingRules != nil {
				peer.AllowedIPs = ifaceSettings.RoutingRules.AllowedIPs
			}

			iface.Peers = append(iface.Peers, peer)

		}
	}

	out.Items = append(out.Items, interfaces...)

	return nil
}

// UpdateInterfaces :
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
		return structs.NewInternalError(err.Error())
	}

	// Create a map for more efficient lookup
	nodeInterfacesMap := map[string]*structs.Interface{}
	for _, i := range nodeInterfaces {
		nodeInterfacesMap[i.ID] = i
	}

	for _, i := range args.Interfaces {
		old, found := nodeInterfacesMap[i.ID]
		if !found {
			return structs.NewInternalError("Interface does not belong to node")
		}

		i = old.Merge(i)
		i.UpdatedAt = time.Now()

		err := s.state.UpsertInterface(ctx, i)
		if err != nil {
			return structs.NewInternalError("Can't update interface")
		}
	}

	return nil
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

// ListNodes retrieves all node entities in the repository
func (s *NodeService) ListNodes(args *structs.NodeListRequest, out *structs.NodeListResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "node", "", NodeList); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	nodes, err := s.state.Nodes(ctx)
	if err != nil {
		return structs.NewInternalError(err.Error())
	}

	out.Items = nil

	for _, n := range nodes {
		out.Items = append(out.Items, n.Stub())
	}

	out.Items = filterNodes(out.Items, args.Filters)

	return nil
}

// JoinNetwork : connects a node to a network
func (s *NetworkService) JoinNetwork(args *structs.NodeJoinNetworkRequest, out *structs.GenericResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "node", args.NodeID, NodeWrite); err != nil {
			return structs.ErrPermissionDenied
		}
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "network", args.NetworkID, NetworkWrite); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	network, err := s.state.NetworkByID(ctx, args.NetworkID)
	if err != nil {
		return structs.ErrNotFound // network not found
	}

	node, err := s.state.NodeByID(ctx, args.NodeID)
	if err != nil {
		return structs.ErrNotFound // node not found
	}

	interfaces, err := s.state.InterfacesByNodeID(ctx, node.ID)
	if err != nil {
		return structs.NewInternalError(err.Error())
	}

	// Check whether node has already joined the network
	for _, iface := range interfaces {
		if iface.NetworkID == network.ID {
			return structs.NewInternalError("Network already joined")
		}
	}

	iface := &structs.Interface{
		ID:        uuid.Generate(),
		NodeID:    node.ID,
		NetworkID: network.ID,
		Name:      nil,               // Setting name is responsibility of the client node
		Address:   nil,               // TODO: set with leasing plugin if it is loaded and enabled
		Peers:     []*structs.Peer{}, // TODO: set with meshing plugin if it is loaded and enabled
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.state.UpsertInterface(ctx, iface)
	if err != nil {
		return structs.NewInternalError("Can't create interface")
	}

	node.UpsertInterface(iface.ID)
	err = s.state.UpsertNode(ctx, node)
	if err != nil {
		return structs.NewInternalError("Can't add interface to node")
	}

	network.UpsertInterface(iface.ID)
	err = s.state.UpsertNetwork(ctx, network)
	if err != nil {
		return structs.NewInternalError("Can't add interface to network")
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
		return structs.NewInternalError("Network does not exist")
	}

	node, err := s.state.NodeByID(ctx, args.NodeID)
	if err != nil {
		return structs.NewInternalError("Node does not exist")
	}

	interfaces, err := s.state.InterfacesByNodeID(ctx, node.ID)
	if err != nil {
		return structs.NewInternalError("Can't retrieve node interfaces")
	}

	// Check whether node is in the network
	for _, iface := range interfaces {
		if iface.NetworkID == network.ID {

			network.RemoveInterface(iface.ID)
			if err := s.state.UpsertNetwork(ctx, network); err != nil {
				return structs.NewInternalError("Can't update network")
			}

			node.RemoveInterface(iface.ID)
			if err := s.state.UpsertNode(ctx, node); err != nil {
				return structs.NewInternalError("Can't update node")
			}

			if err := s.state.DeleteInterfaces(ctx, []string{iface.ID}); err != nil {
				return structs.NewInternalError("Can't delete interface")
			}
		}
	}

	return nil
}

func filterNodes(nodes []*structs.NodeListStub, filters structs.Filters) []*structs.NodeListStub {

	statusFilter := "*"
	if len(filters.Get("status")) > 0 {
		statusFilter = filters.Get("status")[0]
	}

	metaFilters := map[string]string{}
	for _, f := range filters.Get("meta") {
		kv := strings.Split(f, ":")
		if len(kv) != 2 {
			continue
		}
		metaFilters[kv[0]] = kv[1]
	}

	out := []*structs.NodeListStub{}

	for _, n := range nodes {

		matchStatusFilter := true
		if matched, _ := path.Match(statusFilter, n.Status); !matched {
			matchStatusFilter = false
		}

		matchAllMetaFilters := true
		for k, v := range metaFilters {
			if metaValue, found := n.Meta[k]; found {
				if matched, _ := path.Match(v, metaValue); !matched {
					matchAllMetaFilters = false
				}
				continue
			}
			matchAllMetaFilters = false
		}
		if matchStatusFilter && matchAllMetaFilters {
			out = append(out, n)
		}
	}

	return out
}
