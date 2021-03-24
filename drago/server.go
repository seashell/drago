package drago

import (
	"context"
	"fmt"
	"sync"

	auth "github.com/seashell/drago/drago/auth"
	state "github.com/seashell/drago/drago/state"
	inmem "github.com/seashell/drago/drago/state/inmem"
	structs "github.com/seashell/drago/drago/structs"
	acl "github.com/seashell/drago/pkg/acl"
	log "github.com/seashell/drago/pkg/log"
	rpc "github.com/seashell/drago/pkg/rpc"
)

var (
	// AnonymousACLToken is used when no secret is provided,
	// and the request is made anonymously.
	AnonymousACLToken = &structs.ACLToken{
		ID:       "anonymous",
		Name:     "Anonymous Token",
		Type:     structs.ACLTokenTypeClient,
		Policies: []string{"anonymous"},
	}
)

// Server is the Drago server
type Server struct {
	config *Config

	logger      log.Logger
	rpcServer   *rpc.Server
	rpcClient   *rpc.Client
	state       state.Repository
	authHandler auth.AuthorizationHandler

	services struct {
		ACL         *ACLService
		Nodes       *NodeService
		Networks    *NetworkService
		Interfaces  *InterfaceService
		Connections *ConnectionService
		Status      *StatusService
	}

	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

// NewServer creates a new Drago server from the
// configuration, potentially returning an error
func NewServer(config *Config) (*Server, error) {

	s := &Server{
		config:     config,
		logger:     config.Logger.WithName("server"),
		shutdownCh: make(chan struct{}),
	}

	var err error

	err = s.setupACLModel()
	if err != nil {
		s.logger.Errorf("Error setting up ACL model: %s", err.Error())
	}

	err = s.setupApplication()
	if err != nil {
		s.logger.Errorf("Error setting up application modules: %s", err.Error())
	}

	err = s.setupRPCServer()
	if err != nil {
		s.logger.Errorf("Error setting up rpc server: %s", err.Error())
	}

	return s, nil
}

// Stats is used to return statistics for the server
func (s *Server) Stats() map[string]map[string]string {

	stats := map[string]map[string]string{
		"drago": {
			"server": "true",
			"peers":  "",
		},
	}

	return stats
}

// Shutdown tears down the server
func (s *Server) Shutdown() error {
	s.shutdownLock.Lock()
	defer s.shutdownLock.Unlock()

	if s.shutdown {
		s.logger.Infof("already shutdown")
		return nil
	}
	s.logger.Infof("shutting down")

	//s.etcdServer.Close()

	s.shutdown = true
	close(s.shutdownCh)

	return nil
}

func (s *Server) setupApplication() error {

	s.state = inmem.NewStateRepository(s.logger)

	// Setup default policies
	ctx := context.TODO()
	for _, p := range s.defaultACLPolicies() {
		err := s.state.UpsertACLPolicy(ctx, p)
		if err != nil {
			return err
		}
	}

	// err := mock.PopulateRepository(s.state)
	// if err != nil {
	// 	return fmt.Errorf("failed to populate repository with mock data: %v", err)
	// }

	s.authHandler = auth.NewAuthorizationHandler(
		s.config.ACL.Model,
		s.secretResolver(),
		s.policyResolver(),
	)

	nodeService, err := NewNodeService(s.config, s.logger, s.state, s.authHandler)
	if err != nil {
		return fmt.Errorf("failed to create node service: %v", err)
	}

	s.services.Nodes = nodeService
	s.services.ACL = NewACLService(s.config, s.logger, s.state, s.authHandler)
	s.services.Networks = NewNetworkService(s.config, s.logger, s.state, s.authHandler)
	s.services.Interfaces = NewInterfaceService(s.config, s.logger, s.state, s.authHandler)
	s.services.Connections = NewConnectionService(s.config, s.logger, s.state, s.authHandler)

	s.services.Status = NewStatusService(s.config, s.state, s.authHandler)

	return nil
}

// setupACLModel defines an ACL model containing resource types, associated
// capabilities, and aliases which can be used by the application.
func (s *Server) setupACLModel() error {

	model := acl.NewModel()

	model.Resource("token").
		Capabilities(ACLTokenWrite, ACLTokenRead, ACLTokenList).
		Alias("read", ACLTokenRead, ACLTokenList).
		Alias("write", ACLTokenWrite, ACLTokenRead, ACLTokenList)

	model.Resource("policy").
		Capabilities(ACLPolicyWrite, ACLPolicyRead, ACLPolicyList).
		Alias("read", ACLPolicyRead, ACLPolicyList).
		Alias("write", ACLPolicyWrite, ACLPolicyRead, ACLPolicyList)

	model.Resource("network").
		Capabilities(NetworkWrite, NetworkRead, NetworkList).
		Alias("read", NetworkRead, NetworkList).
		Alias("write", NetworkWrite, NetworkRead, NetworkList)

	model.Resource("node").
		Capabilities(NodeWrite, NodeRead, NodeList).
		Alias("read", NodeRead, NodeList).
		Alias("write", NodeWrite, NodeRead, NodeList)

	model.Resource("interface").
		Capabilities(InterfaceWrite, InterfaceRead, InterfaceList).
		Alias("read", InterfaceRead, InterfaceList).
		Alias("write", InterfaceWrite, InterfaceRead, InterfaceList)

	model.Resource("connection").
		Capabilities(ConnectionWrite, ConnectionRead, ConnectionList).
		Alias("read", ConnectionRead, ConnectionList).
		Alias("write", ConnectionWrite, ConnectionRead, ConnectionList)

	s.config.ACL.Model = model

	return nil
}

// returns the ACL policies to be loaded by default into the ACL when the Drago server starts.
func (s *Server) defaultACLPolicies() []*structs.ACLPolicy {
	return []*structs.ACLPolicy{
		{
			Name:        "anonymous",
			Description: "Default policy utilized when no access token is provided",
			Rules: []*structs.ACLPolicyRule{
				{Resource: "token", Path: "*", Capabilities: []string{ACLTokenList}},
				{Resource: "policy", Path: "*", Capabilities: []string{ACLPolicyList}},
				{Resource: "network", Path: "*", Capabilities: []string{NetworkList}},
				{Resource: "node", Path: "*", Capabilities: []string{NodeList}},
				{Resource: "interface", Path: "*", Capabilities: []string{InterfaceList}},
				{Resource: "connection", Path: "*", Capabilities: []string{ConnectionList}},
			},
		},
	}
}

// returns an acl.SecretResolverFunc
func (s *Server) secretResolver() acl.SecretResolverFunc {
	return func(ctx context.Context, secret string) (acl.Token, error) {

		var err error
		var t *structs.ACLToken

		if secret == "" {
			t = AnonymousACLToken
		} else {
			t, err = s.state.ACLTokenBySecret(ctx, secret)
			if err != nil {
				return nil, err
			}
			if t == nil {
				return nil, fmt.Errorf("token not found")
			}
		}

		return auth.NewToken(
			t.Type == structs.ACLTokenTypeManagement,
			t.Policies,
		), nil
	}
}

// returns an acl.SecretResolverFunc
func (s *Server) policyResolver() acl.PolicyResolverFunc {
	return func(ctx context.Context, policy string) (acl.Policy, error) {
		pol, err := s.state.ACLPolicyByName(ctx, policy)
		if err != nil {
			return nil, err
		}

		res := auth.NewPolicy(pol.Name, []acl.Rule{})
		for _, r := range pol.Rules {
			res.AddRule(auth.NewRule(r.Resource, r.Path, r.Capabilities))
		}
		return res, nil
	}
}

func (s *Server) setupRPCServer() error {

	config := &rpc.ServerConfig{
		Logger:      s.logger,
		BindAddress: fmt.Sprintf("%s:%d", s.config.BindAddr, s.config.Ports.RPC),
		Receivers: map[string]interface{}{
			"ACL":        s.services.ACL,
			"Node":       s.services.Nodes,
			"Interface":  s.services.Interfaces,
			"Connection": s.services.Connections,
			"Network":    s.services.Networks,
			"Status":     s.services.Status,
		},
	}

	rpcServer, err := rpc.NewServer(config)
	if err != nil {
		return err
	}

	s.rpcServer = rpcServer

	return nil
}
