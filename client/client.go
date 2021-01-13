package client

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	nic "github.com/seashell/drago/client/nic"
	state "github.com/seashell/drago/client/state"
	boltdb "github.com/seashell/drago/client/state/boltdb"
	structs "github.com/seashell/drago/drago/structs"
	log "github.com/seashell/drago/pkg/log"
	rpc "github.com/seashell/drago/pkg/rpc"
	"github.com/seashell/drago/pkg/uuid"
)

var (
	defaultRegistrationRetryInterval   = 5 * time.Second
	defaultReconciliationRetryInterval = 5 * time.Second
	defaultReconciliationInterval      = 2 * time.Second
	defaultFirstHeartbeatDelay         = 1 * time.Second
	defaultHeartbeatInterval           = 1 * time.Second
)

// Client is the Drago client
type Client struct {
	config *Config

	logger log.Logger

	rpc *rpc.Client

	state state.Repository

	niController     *nic.Controller
	niControllerLock sync.Mutex

	node     *structs.Node
	nodeLock sync.Mutex

	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

// New is used to create a new Drago client from the
// configuration, potentially returning an error
func New(config *Config) (*Client, error) {

	rand.Seed(time.Now().Unix())

	config = DefaultConfig().Merge(config)

	c := &Client{
		config:     config,
		logger:     config.Logger.WithName("client"),
		shutdownCh: make(chan struct{}),
	}

	err := c.setupFilesystem()
	if err != nil {
		return nil, fmt.Errorf("error setting up the filesystem: %v", err)
	}

	err = c.setupRepository()
	if err != nil {
		return nil, fmt.Errorf("error setting up application: %v", err)
	}

	err = c.setupNode()
	if err != nil {
		return nil, fmt.Errorf("error setting up node: %v", err)
	}

	// Setup the network controller
	nc, err := nic.NewController(&nic.Config{
		InterfacePrefix: "drago",
		WireguardPath:   c.config.WireguardPath,
	})
	if err != nil {
		return nil, fmt.Errorf("error setting network controller: %v", err)

	}
	c.niController = nc

	// Start goroutine for registering the client and issuing periodic heartbeats
	go c.registerAndHeartbeat()

	// Start goroutine for reconciling the client state
	go c.run()

	c.logger.Infof("started client node %s", c.NodeID())

	return c, nil
}

func (c *Client) Node() *structs.Node {
	return c.node
}

// NodeID returns the node ID for the given client
func (c *Client) NodeID() string {
	return c.node.ID
}

// NodeSecretID returns the node secret ID for the given client
func (c *Client) NodeSecretID() string {
	return c.node.ID
}

func (c *Client) setupNode() error {

	if c.node == nil {
		c.node = &structs.Node{}
	}

	id, err := c.getNodeID()
	if err != nil {
		return fmt.Errorf("could not retrieve node ID: %v", err)
	}

	secret, err := c.getNodeSecretID()
	if err != nil {
		return fmt.Errorf("could not retrieve node secret ID: %v", err)
	}

	c.node.ID = id
	c.node.SecretID = secret
	c.node.Status = structs.NodeStatusInit

	if c.node.Name == "" {
		c.node.Name, _ = os.Hostname()
	}
	if c.node.Meta == nil {
		c.node.Meta = make(map[string]string)
	}

	if c.node.Name == "" {
		c.node.Name = c.node.ID
	}
	c.node.Status = structs.NodeStatusInit

	return nil
}

func (c *Client) setupFilesystem() error {

	// Ensure the state dir exists
	if c.config.StateDir != "" {
		if err := os.MkdirAll(c.config.StateDir, 0700); err != nil {
			return fmt.Errorf("failed to create state dir: %s", err)
		}
	} else {
		// Otherwise make a temp directory to use.
		tmp, err := c.createTempDir("DragoClient")
		if err != nil {
			return fmt.Errorf("failed to create tmp dir for storing state: %s", err)
		}
		c.config.StateDir = tmp
	}

	c.logger.Infof("using state directory %s", c.config.StateDir)

	return nil
}

func (c *Client) setupRepository() error {

	repo, err := boltdb.NewStateRepository(path.Join(c.config.StateDir, "client.state"))
	if err != nil {
		return fmt.Errorf("failed to open state database: %v", err)
	}
	c.state = repo

	return nil
}

func (c *Client) setupRPCClient() error {

	rpcClient, err := rpc.NewClient(&rpc.ClientConfig{
		Logger:  c.logger,
		Address: c.config.Servers[0],
	})
	if err != nil {
		return err
	}

	c.rpc = rpcClient

	return nil
}

func (c *Client) registerAndHeartbeat() {

	c.tryToRegisterUntilSuccessful()

	heartbeatCh := time.After(defaultFirstHeartbeatDelay)

	for {
		select {
		case <-heartbeatCh:
		case <-c.shutdownCh:
			return
		}

		if err := c.updateNodeStatus(); err != nil {
			c.logger.Debugf("error updating node status: %v", err)
			c.tryToRegisterUntilSuccessful()
			heartbeatCh = time.After(defaultHeartbeatInterval)
		} else {
			heartbeatCh = time.After(defaultHeartbeatInterval)
		}

	}
}

type interfacesUpdate struct {
	interfaces []*structs.Interface
}

func (c *Client) run() {

	c.logger.Debugf("running client node")

	interfacesUpdateCh := make(chan *interfacesUpdate)
	go c.watchInterfaces(interfacesUpdateCh)

	for {
		select {
		case update := <-interfacesUpdateCh:
			c.shutdownLock.Lock()
			if c.shutdown {
				c.shutdownLock.Unlock()
				return
			}

			c.applyInterfaceUpdate(update)
			c.shutdownLock.Unlock()

		case <-c.shutdownCh:
			return
		}
	}
}

func (c *Client) applyInterfaceUpdate(i interfacesUpdate) {

	c.interfacesLock.RLock()
	existing := make(map[string]uint64, len(c.allocs))
	for id, ar := range c.allocs {
		existing[id] = ar.Alloc().AllocModifyIndex
	}
	c.allocLock.RUnlock()

}

func (c *Client) watchInterfaces(ch chan *interfacesUpdate) {

	c.logger.Debugf("watching interfaces")

	req := &structs.NodeSpecificRequest{
		ID: c.NodeID(),
	}

	for {
		var resp structs.NodeInterfacesResponse
		err := c.RPC("Node.GetInterfaces", req, &resp)
		if err != nil {
			c.logger.Debugf("error fetching interfaces: %v", err)
			retryCh := time.After(defaultReconciliationRetryInterval)
			select {
			case <-retryCh:
			case <-c.shutdownCh:
				return
			}
		}

		ch <- &interfacesUpdate{
			interfaces: resp.Items,
		}

		retryCh := time.After(c.config.ReconcileInterval)
		select {
		case <-c.shutdownCh:
			return
		case <-retryCh:
		}
	}
}

func (c *Client) tryToRegisterUntilSuccessful() {

	for {

		req := &structs.NodeRegisterRequest{
			Node: c.Node(),
		}

		c.logger.Infof("registering node...")

		var err error
		var resp structs.NodeUpdateResponse
		if err = c.RPC("Node.Register", req, &resp); err == nil {

			c.logger.Infof("node successfully registered")

			c.nodeLock.Lock()
			c.node.Status = structs.NodeStatusReady
			c.nodeLock.Unlock()

			return
		}

		c.logger.Errorf("error registering node: %v", err)

		retryCh := time.After(time.Duration(defaultRegistrationRetryInterval))

		select {
		case <-retryCh:
		case <-c.shutdownCh:
			return
		}
	}
}

func (c *Client) updateNodeStatus() error {

	req := &structs.NodeUpdateStatusRequest{
		ID:     c.NodeID(),
		Status: structs.NodeStatusReady,
	}

	var err error
	var resp structs.NodeUpdateResponse
	if err = c.RPC("Node.UpdateStatus", req, &resp); err != nil {
		return err
	}

	return nil
}

// Shutdown is used to tear down the client
func (c *Client) Shutdown() error {
	c.shutdownLock.Lock()
	defer c.shutdownLock.Unlock()

	if c.shutdown {
		c.logger.Infof("client already shutdown")
		return nil
	}
	c.logger.Infof("shutting down")

	c.shutdown = true
	close(c.shutdownCh)

	return nil
}

func (c *Client) createTempDir(pattern string) (string, error) {
	p, err := ioutil.TempDir("", pattern)
	if err != nil {
		return "", fmt.Errorf("could not create temporary directory: %v", err)
	}
	p, err = filepath.EvalSymlinks(p)
	if err != nil {
		return "", fmt.Errorf("could not retrieve path to StateDir: %v", err)
	}
	return p, nil
}

// RPC calls a RPC method on a remote server using the clients RPC client, establishing
// the connection if it's being used for the first time, of if it has been disconnected.
func (c *Client) RPC(method string, args interface{}, reply interface{}) error {

	if c.rpc == nil {
		err := c.setupRPCClient()
		if err != nil {
			return err
		}
	}

	if err := c.rpc.Call(method, args, reply); err != nil {
		c.rpc = nil
		return err
	}

	return nil
}

func (c *Client) getNodeID() (string, error) {

	id := uuid.Generate()

	if c.config.DevMode {
		return id, nil
	}

	path := filepath.Join(c.config.StateDir, "client-id")

	id, err := c.readFileLazy(path, id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (c *Client) getNodeSecretID() (string, error) {

	secret := uuid.Generate()

	path := filepath.Join(c.config.StateDir, "secret-id")

	secret, err := c.readFileLazy(path, secret)
	if err != nil {
		return "", err
	}

	return secret, nil
}
