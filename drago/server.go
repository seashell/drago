package drago

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strings"
	"sync"

	handler "github.com/seashell/drago/drago/adapter/http"
	middleware "github.com/seashell/drago/drago/adapter/http/middleware"
	auth "github.com/seashell/drago/drago/auth"
	state "github.com/seashell/drago/drago/state"
	inmem "github.com/seashell/drago/drago/state/inmem"
	structs "github.com/seashell/drago/drago/structs"
	acl "github.com/seashell/drago/pkg/acl"
	http "github.com/seashell/drago/pkg/http"
	rpc "github.com/seashell/drago/pkg/rpc"
	"github.com/shurcooL/go-goon"
	embed "go.etcd.io/etcd/embed"
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

	etcdServer *embed.Etcd
	httpServer *http.Server
	rpcServer  *rpc.Server

	rpcClient *rpc.Client

	state state.Repository

	authHandler auth.AuthorizationHandler

	services struct {
		ACL      *ACLService
		Networks *NetworkService
		Status   *StatusService
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
		shutdownCh: make(chan struct{}),
	}

	var err error

	//err = s.setupEtcdServer()
	//if err != nil {
	//	s.logger.Errorf("Error setting up etcd server: %s", err.Error())
	//}

	//err = s.setupEtcdClient()
	//if err != nil {
	//	s.logger.Errorf("Error setting up etcd client: %s", err.Error())
	//}

	err = s.setupACLModel()
	if err != nil {
		s.config.Logger.Errorf("Error setting up ACL model: %s", err.Error())
	}

	err = s.setupApplication()
	if err != nil {
		s.config.Logger.Errorf("Error setting up application modules: %s", err.Error())
	}

	err = s.setupRPCServer()
	if err != nil {
		s.config.Logger.Errorf("Error setting up rpc server: %s", err.Error())
	}

	err = s.setupRPCClient()
	if err != nil {
		s.config.Logger.Errorf("Error setting up rpc client: %s", err.Error())
	}

	err = s.setupHTTPServer()
	if err != nil {
		s.config.Logger.Errorf("Error setting up http server: %s", err.Error())
	}

	return s, nil
}

// State returns a repository containing the server state
func (s *Server) State() state.Repository {
	return s.state
}

// Shutdown tears down the server
func (s *Server) Shutdown() error {
	s.shutdownLock.Lock()
	defer s.shutdownLock.Unlock()

	if s.shutdown {
		s.config.Logger.Infof("already shutdown")
		return nil
	}
	s.config.Logger.Infof("shutting down")

	s.etcdServer.Close()

	s.shutdown = true
	close(s.shutdownCh)

	return nil
}

func (s *Server) setupApplication() error {

	s.state = inmem.NewStateRepository(s.config.Logger)

	// Setup default policies
	ctx := context.TODO()
	for _, p := range s.defaultACLPolicies() {
		err := s.State().UpsertACLPolicy(ctx, p)
		if err != nil {
			return err
		}
	}

	s.authHandler = auth.NewAuthorizationHandler(
		s.config.ACL.Model,
		s.secretResolver(),
		s.policyResolver(),
	)

	s.services.ACL = NewACLService(s.config, s.State(), s.authHandler)
	s.services.Networks = NewNetworkService(s.config, s.State(), s.authHandler)
	s.services.Status = NewStatusService(s.config, s.State(), s.authHandler)

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
			t, err = s.State().ACLTokenBySecret(ctx, secret)
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
		pol, err := s.State().ACLPolicyByName(ctx, policy)
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

func (s *Server) setupEtcdServer() error {

	cfg := embed.NewConfig()

	// Advertise peer URLs
	apURLs, err := parseUrls(s.config.Etcd.InitialAdvertisePeerURLs)
	if err != nil {
		return err
	}

	// Listen peer URLs
	lpURLs, err := parseUrls(s.config.Etcd.ListenPeerURLs)
	if err != nil {
		return err
	}

	// Advertise client URLs
	acURLs, err := parseUrls(s.config.Etcd.InitialAdvertiseClientURLs)
	if err != nil {
		return err
	}

	// Listen client URLs
	lcURLs, err := parseUrls(s.config.Etcd.ListenClientURLs)
	if err != nil {
		return err
	}

	cfg.Name = s.config.Etcd.Name
	cfg.Dir = path.Join(s.config.DataDir, "/etcd")
	cfg.WalDir = path.Join(s.config.DataDir, "/etcd", "/wal")
	cfg.Logger = "zap"

	cfg.APUrls = apURLs
	cfg.LPUrls = lpURLs
	cfg.ACUrls = acURLs
	cfg.LCUrls = lcURLs

	cfg.LogOutputs = []string{"stderr", path.Join(s.config.DataDir, "/etcd.log")}
	cfg.LogLevel = strings.ToLower(s.config.LogLevel)

	s.config.Logger.Infof("starting etcd server")

	etcdServer, err := embed.StartEtcd(cfg)
	if err != nil {
		return err
	}

	s.etcdServer = etcdServer

	return nil
}

func parseUrls(urls []string) ([]url.URL, error) {
	res := []url.URL{}
	for _, v := range urls {
		url, err := url.Parse(v)
		if err != nil {
			return nil, err
		}
		res = append(res, *url)
	}
	return res, nil
}

func (s *Server) setupHTTPServer() error {

	config := &http.Config{
		Logger:      s.config.Logger,
		BindAddress: fmt.Sprintf("%s:%d", s.config.BindAddr, s.config.Ports.HTTP),
		Handlers: map[string]http.Handler{
			"/api/acl/networks/": handler.NewNetworkHandler(s.rpcClient),
			"/api/acl/policies/": handler.NewACLPolicyHandler(s.rpcClient),
			"/api/acl/tokens/":   handler.NewACLTokenHandler(s.rpcClient),
			"/api/acl/":          handler.NewACLHandler(s.rpcClient),
			"/status":            handler.NewStatusHandler(s.rpcClient),
		},
		Middleware: []http.Middleware{
			middleware.CORS(),
			middleware.Logging(s.config.Logger),
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
		Logger:      s.config.Logger,
		BindAddress: fmt.Sprintf("%s:%d", s.config.BindAddr, s.config.Ports.RPC),
		Receivers: map[string]interface{}{
			"ACL":     s.services.ACL,
			"Network": s.services.Networks,
			"Status":  s.services.Status,
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
		Logger:  s.config.Logger,
		Address: fmt.Sprintf("%s:%d", s.config.BindAddr, s.config.Ports.RPC),
	}

	rpcClient, err := rpc.NewClient(config)
	if err != nil {
		return err
	}

	s.rpcClient = rpcClient

	return nil
}
