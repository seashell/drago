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

	"github.com/seashell/drago/agent/conn"
	nic "github.com/seashell/drago/client/nic"
	state "github.com/seashell/drago/client/state"
	boltdb "github.com/seashell/drago/client/state/boltdb"
	structs "github.com/seashell/drago/drago/structs"
	log "github.com/seashell/drago/pkg/log"
	uuid "github.com/seashell/drago/pkg/uuid"
)

var (
	defaultRegistrationRetryInterval   = 5 * time.Second
	defaultReconciliationRetryInterval = 5 * time.Second
	defaultReconciliationInterval      = 2 * time.Second
	defaultFirstHeartbeatDelay         = 1 * time.Second
	defaultHeartbeatInterval           = 5 * time.Second
)

// Client is the Drago client
type Client struct {
	config *Config

	logger log.Logger

	rpc conn.RPCConnection

	state state.Repository

	niController     nic.NetworkInterfaceController
	niControllerLock sync.Mutex

	node     *structs.Node
	nodeLock sync.Mutex

	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

// New is used to create a new Drago client from the
// configuration, potentially returning an error
func New(conn conn.RPCConnection, config *Config) (*Client, error) {

	rand.Seed(time.Now().Unix())

	config = DefaultConfig().Merge(config)

	c := &Client{
		config:     config,
		rpc:        conn,
		logger:     config.Logger.WithName("client"),
		shutdownCh: make(chan struct{}),
	}

	if err := c.setupState(); err != nil {
		return nil, fmt.Errorf("error setting up client state: %v", err)
	}

	if err := c.setupNode(); err != nil {
		return nil, fmt.Errorf("error setting up node: %v", err)
	}

	if err := c.setupNetworkController(); err != nil {
		return nil, fmt.Errorf("error setting up network controller: %v", err)
	}

	if err := c.setupInterfaces(); err != nil {
		return nil, fmt.Errorf("error setting up interfaces: %v", err)
	}

	go c.registerAndHeartbeat()

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

// Stats is used to return statistics for the server
func (c *Client) Stats() map[string]map[string]string {

	stats := map[string]map[string]string{
		"client": {
			"node_id": c.NodeID(),
		},
	}

	return stats
}

func (c *Client) setupNode() error {

	if c.node == nil {
		c.node = &structs.Node{}
	}

	id, err := c.readOrGenerateNodeID()
	if err != nil {
		return fmt.Errorf("could not retrieve node ID: %v", err)
	}

	secret, err := c.readOrGenerateNodeSecretID()
	if err != nil {
		return fmt.Errorf("could not retrieve node secret ID: %v", err)
	}

	c.node.ID = id
	c.node.SecretID = secret

	c.node.Name = c.config.Name
	c.node.Meta = c.config.Meta

	c.node.AdvertiseAddress = c.config.AdvertiseAddress
	c.node.Status = structs.NodeStatusInit

	if c.node.Name == "" {
		if hostname, _ := os.Hostname(); hostname != "" {
			c.node.Name = hostname
		} else {
			c.node.Name = c.node.ID
		}
	}
	if c.node.Meta == nil {
		c.node.Meta = make(map[string]string)
	}

	return nil
}

func (c *Client) setupState() error {

	// Ensure the state dir exists. If it was not was specified,
	// create a temporary directory to store the client state.
	if c.config.StateDir != "" {
		if err := os.MkdirAll(c.config.StateDir, 0700); err != nil {
			return fmt.Errorf("failed to create state dir: %s", err)
		}
	} else {
		tmp, err := c.createTempDir("DragoClient")
		if err != nil {
			return fmt.Errorf("failed to create tmp dir for storing state: %s", err)
		}
		c.config.StateDir = tmp
	}

	c.logger.Infof("using state directory %s", c.config.StateDir)

	repo, err := boltdb.NewStateRepository(path.Join(c.config.StateDir, "client.state"))
	if err != nil {
		return fmt.Errorf("failed to open state database: %v", err)
	}

	c.state = repo

	return nil
}

func (c *Client) setupNetworkController() error {

	nc, err := nic.NewController(&nic.Config{
		InterfacesPrefix: c.config.InterfacesPrefix,
		WireguardPath:    c.config.WireguardPath,
		KeyStore:         c.state, // TODO: improve how we store private keys (do we really need to store them?)
	})
	if err != nil {
		return err
	}

	c.niController = nc

	return nil
}

func (c *Client) setupInterfaces() error {

	current, err := c.niController.Interfaces()
	if err != nil {
		return fmt.Errorf("could not retrieve network interfaces from network controller: %v", err)
	}

	desired, err := c.state.Interfaces()
	if err != nil {
		return fmt.Errorf("could not retrieve network interfaces from state: %v", err)
	}

	c.reconcileInterfaces(current, desired)

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

		c.logger.Debugf("heartbeating (client -> server)")

		if err := c.updateNodeStatus(); err != nil {
			c.logger.Debugf("error updating node status: %v", err)
			c.tryToRegisterUntilSuccessful()
			heartbeatCh = time.After(defaultHeartbeatInterval)
		} else {
			heartbeatCh = time.After(defaultHeartbeatInterval)
		}

	}
}

func (c *Client) run() {

	c.logger.Debugf("running node")

	interfacesUpdateCh := make(chan []*structs.Interface)
	go c.watchInterfaces(interfacesUpdateCh)

	go c.synchronizeInterfaces()

	for {
		select {
		case desired := <-interfacesUpdateCh:
			c.shutdownLock.Lock()
			if c.shutdown {
				c.shutdownLock.Unlock()
				return
			}

			current, err := c.state.Interfaces()
			if err != nil {
				c.logger.Errorf("could not read interfaces from state repository: %v", err)
			}

			c.reconcileInterfaces(current, desired)

			c.shutdownLock.Unlock()
		case <-c.shutdownCh:
			return
		}
	}
}

func (c *Client) reconcileInterfaces(current, desired []*structs.Interface) {

	currentMap := map[string]*structs.Interface{}
	for _, i := range current {
		currentMap[i.ID] = i
	}

	desiredMap := map[string]*structs.Interface{}
	for _, i := range desired {
		desiredMap[i.ID] = i
	}

	diff := interfacesDiff(currentMap, desiredMap)

	c.logger.Debugf("interface updates: (created: %d, deleted: %d, updated: %d, unchanged: %d)",
		len(diff.created), len(diff.deleted), len(diff.updated), len(diff.unchanged))

	c.niControllerLock.Lock()
	defer c.niControllerLock.Unlock()

	// Delete old interfaces
	for _, id := range diff.deleted {
		if err := c.state.DeleteInterfaces([]string{id}); err != nil {
			c.logger.Warnf("could not persist interface deletion to the state: %v", err)
		}
		if err := c.niController.DeleteInterfaceByAlias(id); err != nil {
			c.logger.Warnf("could not delete interface: %v", err)
		}
	}

	// Create a new interface from scratch
	for _, id := range diff.created {

		iface := desiredMap[id]

		err := c.state.UpsertInterface(iface)
		if err != nil {
			c.logger.Warnf("could not persist interface: %v", err)
			continue
		}

		err = c.niController.CreateInterface(iface)
		if err != nil {
			c.logger.Warnf("could not create wireguard interface: %v", err)
		}
	}

	// Update an existing interface
	for _, id := range diff.updated {

		iface := desiredMap[id]

		err := c.state.UpsertInterface(iface)
		if err != nil {
			c.logger.Warnf("could not persist interface: %v", err)
			continue
		}

		if err := c.niController.UpdateInterface(iface); err != nil {
			c.logger.Warnf("could not update wireguard interface: %v", err)
		}
	}

}

func (c *Client) watchInterfaces(ch chan []*structs.Interface) {

	req := &structs.NodeSpecificRequest{
		NodeID: c.NodeID(),
	}

	for {

		c.logger.Debugf("updating interface configuration (server -> client)")

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
		} else {
			ch <- resp.Items
		}

		retryCh := time.After(c.config.ReconcileInterval)
		select {
		case <-c.shutdownCh:
			return
		case <-retryCh:
		}
	}
}

func (c *Client) synchronizeInterfaces() {

	for {

		c.logger.Debugf("updating interface status (client -> server)")

		interfaces, err := c.niController.Interfaces()
		if err != nil {
			c.logger.Warnf("could not retrieve interfaces with public key from controller")
		}

		if interfaces == nil {
			interfaces = []*structs.Interface{}
		}

		req := &structs.NodeInterfaceUpdateRequest{
			NodeID:     c.NodeID(),
			Interfaces: interfaces,
		}

		var resp structs.GenericResponse
		err = c.RPC("Node.UpdateInterfaces", req, &resp)
		if err != nil {
			c.logger.Debugf("error updating interfaces: %v", err)
			retryCh := time.After(randomDuration(defaultReconciliationInterval, 0*time.Second))
			select {
			case <-retryCh:
			case <-c.shutdownCh:
				return
			}
		}
		retryCh := time.After(randomDuration(defaultReconciliationInterval, 0*time.Second))
		select {
		case <-c.shutdownCh:
			return
		case <-retryCh:
		}
	}
}

func (c *Client) tryToRegisterUntilSuccessful() {

	for {

		c.logger.Debugf("registering node (client -> server)")

		req := &structs.NodeRegisterRequest{
			Node: c.Node(),
		}

		var err error
		var resp structs.NodeUpdateResponse
		if err = c.RPC("Node.Register", req, &resp); err == nil {

			c.nodeLock.Lock()
			c.node.Status = structs.NodeStatusReady
			c.nodeLock.Unlock()

			return
		}

		c.logger.Debugf("error registering node: %v", err)

		retryCh := time.After(randomDuration(defaultRegistrationRetryInterval, 0*time.Second))

		select {
		case <-retryCh:
		case <-c.shutdownCh:
			return
		}
	}
}

func (c *Client) updateNodeStatus() error {

	c.logger.Debugf("updating node status (client -> server)")

	req := &structs.NodeUpdateStatusRequest{
		NodeID:           c.NodeID(),
		Status:           structs.NodeStatusReady,
		AdvertiseAddress: c.Node().AdvertiseAddress,
		Meta:             c.node.Meta,
	}

	var err error
	var resp structs.NodeUpdateResponse
	if err = c.RPC("Node.UpdateStatus", req, &resp); err != nil {
		return err
	}

	return nil
}

func (c *Client) RPC(method string, args interface{}, reply interface{}) error {
	return c.rpc.Call(method, args, reply)
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

func (c *Client) readOrGenerateNodeID() (string, error) {

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

func (c *Client) readOrGenerateNodeSecretID() (string, error) {

	path := filepath.Join(c.config.StateDir, "secret-id")

	secret, err := c.readFileLazy(path, uuid.Generate())
	if err != nil {
		return "", err
	}

	return secret, nil
}

// Generates a random duration in the interval [mean-delta, mean+delta]
func randomDuration(mean time.Duration, delta time.Duration) time.Duration {
	t := mean.Milliseconds() + int64((rand.Float32()-0.5)*float32(delta.Milliseconds()))
	return time.Duration(t * int64(time.Millisecond))
}
