package drago

import (
	"context"
	"fmt"
	"sync"

	handler "github.com/seashell/drago/drago/adapter/http"
	middleware "github.com/seashell/drago/drago/adapter/http/middleware"
	auth "github.com/seashell/drago/drago/auth"
	"github.com/seashell/drago/drago/mock"
	state "github.com/seashell/drago/drago/state"
	inmem "github.com/seashell/drago/drago/state/inmem"
	structs "github.com/seashell/drago/drago/structs"
	acl "github.com/seashell/drago/pkg/acl"
	http "github.com/seashell/drago/pkg/http"
	log "github.com/seashell/drago/pkg/log"
	rpc "github.com/seashell/drago/pkg/rpc"
	"github.com/shurcooL/go-goon"
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

	logger log.Logger

	httpServer *http.Server
	rpcServer  *rpc.Server

	rpcClient *rpc.Client

	state state.Repository

	authHandler auth.AuthorizationHandler

	services struct {
		ACL        *ACLService
		Nodes      *NodeService
		Networks   *NetworkService
		Interfaces *InterfaceService
		Status     *StatusService
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

	err = s.setupRPCClient()
	if err != nil {
		s.logger.Errorf("Error setting up rpc client: %s", err.Error())
	}

	err = s.setupHTTPServer()
	if err != nil {
		s.logger.Errorf("Error setting up http server: %s", err.Error())
	}

	return s, nil
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

	err := mock.PopulateWithData(s.state)
	if err != nil {
		return fmt.Errorf("failed to populate repository with mock data: %v", err)
	}

	fmt.Println("1")

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

	s.config.ACL.Model = model

	return nil
}

// returns the ACL policies to be loaded by default into the ACL when the Drago server starts.
func (s *Server) defaultACLPolicies() []*structs.ACLPolicy {
	return []*structs.ACLPolicy{
		{
			Name: "anonymous",
			Rules: []*structs.ACLPolicyRule{
				{"token", "*", []string{}},
				{"policy", "*", []string{}},
			},
		},
		{
			Name:  "manager",
			Rules: []*structs.ACLPolicyRule{},
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

		goon.Dump(t)

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

func (s *Server) setupHTTPServer() error {

	config := &http.Config{
		Logger:      s.logger,
		BindAddress: fmt.Sprintf("%s:%d", s.config.BindAddr, s.config.Ports.HTTP),
		Handlers: map[string]http.Handler{
			"/api/nodes/":        handler.NewNodeHandler(s.rpcClient),
			"/api/interfaces/":   handler.NewInterfaceHandler(s.rpcClient),
			"/api/networks/":     handler.NewNetworkHandler(s.rpcClient),
			"/api/acl/policies/": handler.NewACLPolicyHandler(s.rpcClient),
			"/api/acl/tokens/":   handler.NewACLTokenHandler(s.rpcClient),
			"/api/acl/":          handler.NewACLHandler(s.rpcClient),
			"/status":            handler.NewStatusHandler(s.rpcClient),
		},
		Middleware: []http.Middleware{
			middleware.CORS(),
			middleware.Logging(s.logger),
		},
	}

	if s.config.UI {
		//config.Handlers["/ui/"] = handler.NewFilesystemHandlerAdapter(ui.Bundle)
		config.Handlers["/"] = handler.NewFallthroughHandler("/ui/")
	}

	httpServer, err := http.NewServer(config)
	if err != nil {
		return err
	}

	s.httpServer = httpServer

	return nil
}

func (s *Server) setupRPCServer() error {

	config := &rpc.ServerConfig{
		Logger:      s.logger,
		BindAddress: fmt.Sprintf("%s:%d", s.config.BindAddr, s.config.Ports.RPC),
		Receivers: map[string]interface{}{
			"ACL":       s.services.ACL,
			"Node":      s.services.Nodes,
			"Interface": s.services.Interfaces,
			"Network":   s.services.Networks,
			"Status":    s.services.Status,
		},
	}

	rpcServer, err := rpc.NewServer(config)
	if err != nil {
		return err
	}

	s.rpcServer = rpcServer

	return nil
}

func (s *Server) setupRPCClient() error {

	config := &rpc.ClientConfig{
		Logger:  s.logger,
		Address: fmt.Sprintf("%s:%d", s.config.BindAddr, s.config.Ports.RPC),
	}

	rpcClient, err := rpc.NewClient(config)
	if err != nil {
		return err
	}

	s.rpcClient = rpcClient

	return nil
}
