package drago

import (
	"sync"

	application "github.com/seashell/drago/drago/application"
	http "github.com/seashell/drago/drago/http"
	rpc "github.com/seashell/drago/drago/rpc"
	log "github.com/seashell/drago/pkg/log"
	clientv3 "go.etcd.io/etcd/clientv3"
	embed "go.etcd.io/etcd/embed"
)

// Server is the Drago server
type Server struct {
	config *Config
	logger log.Logger

	httpServer *http.Server
	rpcServer  *rpc.Server

	etcdServer *embed.Etcd
	etcdClient *clientv3.Client

	services struct {
		acl      application.ACLService
		tokens   application.ACLTokenService
		policies application.ACLPolicyService
	}

	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

// NewServer is used to create a new Drago server from the
// configuration, potentially returning an error
func NewServer(config *Config) (*Server, error) {

	s := &Server{
		config:     config,
		logger:     config.Logger,
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

	err = s.setupApplication()
	if err != nil {
		s.logger.Errorf("Error setting up application modules: %s", err.Error())
	}

	err = s.setupACLSystem()
	if err != nil {
		s.logger.Errorf("Error setting up acl system: %s", err.Error())
	}

	err = s.setupHTTPServer()
	if err != nil {
		s.logger.Errorf("Error setting up http server: %s", err.Error())
	}

	err = s.setupRPCServer()
	if err != nil {
		s.logger.Errorf("Error setting up rpc server: %s", err.Error())
	}

	return s, nil
}

// Shutdown is used to tear down the server
func (s *Server) Shutdown() error {
	s.shutdownLock.Lock()
	defer s.shutdownLock.Unlock()

	if s.shutdown {
		s.logger.Infof("already shutdown")
		return nil
	}
	s.logger.Infof("shutting down")

	s.etcdServer.Close()
	s.etcdClient.Close()

	s.shutdown = true
	close(s.shutdownCh)

	return nil
}
