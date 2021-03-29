package drago

import (
	"context"
	"fmt"
	"time"

	auth "github.com/seashell/drago/drago/auth"
	state "github.com/seashell/drago/drago/state"
	structs "github.com/seashell/drago/drago/structs"
	log "github.com/seashell/drago/pkg/log"
	uuid "github.com/seashell/drago/pkg/uuid"
)

const (
	ConnectionList  = "list"
	ConnectionRead  = "read"
	ConnectionWrite = "write"
)

type ConnectionService struct {
	config      *Config
	logger      log.Logger
	state       state.Repository
	authHandler auth.AuthorizationHandler
}

// NewConnectionService ...
func NewConnectionService(config *Config, logger log.Logger, state state.Repository, authHandler auth.AuthorizationHandler) *ConnectionService {
	return &ConnectionService{
		config:      config,
		logger:      logger,
		state:       state,
		authHandler: authHandler,
	}
}

// GetConnection returns a Connection entity by ID
func (s *ConnectionService) GetConnection(args *structs.ConnectionSpecificRequest, out *structs.SingleConnectionResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "connection", args.ConnectionID, ConnectionRead); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	n, err := s.state.ConnectionByID(ctx, args.ConnectionID)
	if err != nil {
		return structs.ErrNotFound
	}

	out.Connection = n

	return nil
}

// ListConnections retrieves all connection entities in the repository
func (s *ConnectionService) ListConnections(args *structs.ConnectionListRequest, out *structs.ConnectionListResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "connection", "", ConnectionList); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	out.Items = nil

	var err error
	var connections []*structs.Connection

	if args.InterfaceID != "" {
		if connections, err = s.state.ConnectionsByInterfaceID(ctx, args.InterfaceID); err != nil {
			return structs.ErrInternal
		}
	} else if args.NodeID != "" {
		if connections, err = s.state.ConnectionsByNodeID(ctx, args.NodeID); err != nil {
			return structs.ErrInternal
		}
	} else if args.NetworkID != "" {
		if connections, err = s.state.ConnectionsByNetworkID(ctx, args.NetworkID); err != nil {
			return structs.ErrInternal
		}
	} else {
		if connections, err = s.state.Connections(ctx); err != nil {
			return structs.ErrInternal
		}
	}

	for _, c := range connections {
		shouldAppend := true
		if args.NetworkID != "" && c.NetworkID != args.NetworkID {
			shouldAppend = false
		}
		if args.InterfaceID != "" && !c.ConnectsInterface(args.InterfaceID) {
			shouldAppend = false
		}
		if shouldAppend {
			out.Items = append(out.Items, c.Stub())
		}
	}

	return nil
}

// UpsertConnection upserts a new Connection entity
func (s *ConnectionService) UpsertConnection(args *structs.ConnectionUpsertRequest, out *structs.GenericResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "connection", "", ConnectionWrite); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	c := args.Connection

	err := c.Validate()
	if err != nil {
		return structs.NewInvalidInputError(err.Error())
	}

	// If the connection already exists, we simply merge the new values into the existing struct.
	// Otherwise, we generate a new ID and set the protected attributes in preparation for inserting
	// the new struct into the repository.
	if c.ID != "" {
		old, err := s.state.ConnectionByID(ctx, c.ID)
		if err != nil {
			return structs.ErrNotFound // connection does not exist
		}
		c = old.Merge(c)
	} else {
		c.ID = uuid.Generate()
		c.NodeIDs = []string{}
		c.CreatedAt = time.Now()
		c.ModifyIndex = 0
	}

	connectedInterfaceIDs := c.ConnectedInterfaceIDs()

	if len(connectedInterfaceIDs) != 2 {
		return structs.NewInternalError("A connection must specify exactly two interfaces")
	}
	if connectedInterfaceIDs[0] == connectedInterfaceIDs[1] {
		return structs.NewInternalError("Can't connect an interface to itself")
	}
	// Make sure interfaces are not already connected
	if conn, err := s.state.ConnectionByInterfaceIDs(ctx, connectedInterfaceIDs[0], connectedInterfaceIDs[1]); err == nil {
		if conn.ID != c.ID {
			return structs.NewInternalError("Interfaces already connected")
		}
	}

	// Make sure both peer settings are correctly initialized
	for _, id := range connectedInterfaceIDs {

		// Initialize PeerSettings
		if c.PeerSettings[id] == nil {
			c.PeerSettings[id] = &structs.PeerSettings{InterfaceID: id, RoutingRules: &structs.RoutingRules{AllowedIPs: []string{}}}
		}

		// Initialize InterfaceID
		if c.PeerSettings[id].InterfaceID == "" {
			c.PeerSettings[id].InterfaceID = id
		} else if c.PeerSettings[id].InterfaceID != id {
			return structs.NewInternalError("Peer ID mismatch")
		}

		// Initialize RoutingRules, if necessary
		if c.PeerSettings[id].RoutingRules == nil {
			c.PeerSettings[id].RoutingRules = &structs.RoutingRules{AllowedIPs: []string{}}
		}
	}

	// Make sure both peer interfaces exist
	ifaces := []*structs.Interface{}
	for _, id := range connectedInterfaceIDs {
		if iface, err := s.state.InterfaceByID(ctx, id); err == nil {
			ifaces = append(ifaces, iface)
			continue
		}
		return structs.NewInternalError(fmt.Sprintf("Interface %s does not exist", id))
	}

	if ifaces[0].NetworkID != ifaces[1].NetworkID {
		return structs.NewInternalError("Interfaces are not in the same network")
	}

	// Assign network ID in case we're creating a new connection
	c.NetworkID = ifaces[0].NetworkID

	c.UpdatedAt = time.Now()
	c.ModifyIndex++

	// TODO: wrap in a transaction

	for _, iface := range ifaces {
		iface.UpsertConnection((c.ID))
		if err = s.state.UpsertInterface(ctx, iface); err != nil {
			return structs.ErrInternal
		}

		if node, err := s.state.NodeByID(ctx, iface.NodeID); err == nil {
			c.NodeIDs = append(c.NodeIDs, iface.NodeID)
			node.UpsertConnection((c.ID))
			if err = s.state.UpsertNode(ctx, node); err != nil {
				return structs.ErrInternal
			}
		}
	}

	if network, err := s.state.NetworkByID(ctx, c.NetworkID); err == nil {
		network.UpsertConnection((c.ID))
		if err = s.state.UpsertNetwork(ctx, network); err != nil {
			return structs.ErrInternal
		}
	}

	err = s.state.UpsertConnection(ctx, c)
	if err != nil {
		return structs.ErrInternal
	}

	return nil
}

// DeleteConnection deletes a connection entity from the repository
func (s *ConnectionService) DeleteConnection(args *structs.ConnectionDeleteRequest, out *structs.GenericResponse) error {

	ctx := context.TODO()

	// Check if authorized
	if s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "connection", "", ConnectionWrite); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	for _, connID := range args.ConnectionIDs {
		if conn, err := s.state.ConnectionByID(ctx, connID); err == nil {

			var nodes []*structs.Node
			var ifaces []*structs.Interface
			var network *structs.Network

			for _, nodeID := range conn.NodeIDs {
				if node, err := s.state.NodeByID(ctx, nodeID); err == nil {
					nodes = append(nodes, node)
				}
			}

			for _, ifaceID := range conn.ConnectedInterfaceIDs() {
				if iface, err := s.state.InterfaceByID(ctx, ifaceID); err == nil {
					ifaces = append(ifaces, iface)
				}
			}

			if network, err = s.state.NetworkByID(ctx, conn.NetworkID); err != nil {
				return structs.NewInternalError(err.Error())
			}

			for _, node := range nodes {
				node.RemoveConnection(connID)
				if err := s.state.UpsertNode(ctx, node); err != nil {
					return structs.ErrInternal // could not update node
				}
			}

			for _, iface := range ifaces {
				iface.RemoveConnection(connID)
				if err = s.state.UpsertInterface(ctx, iface); err != nil {
					return structs.ErrInternal // could not update interface
				}
			}

			network.RemoveConnection(connID)
			if err := s.state.UpsertNetwork(ctx, network); err != nil {
				return structs.ErrInternal // could not update network
			}
		}
	}

	// Remove connections
	if err := s.state.DeleteConnections(ctx, args.ConnectionIDs); err != nil {
		return structs.ErrInternal
	}

	return nil
}
