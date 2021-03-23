package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/seashell/drago/drago/state"
	"github.com/seashell/drago/drago/structs/config"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/embed"
)

const (
	defaultNamespace = "default"
	defaultPrefix    = "/registry"

	resourceTypeACLPolicy  = "policy"
	resourceTypeACLToken   = "token"
	resourceTypeNetwork    = "network"
	resourceTypeNode       = "node"
	resourceTypeInterface  = "interface"
	resourceTypeConnection = "connection"

	transactionContextKey = "etcdtxn"
)

type Config struct {
	DataDir  string
	LogLevel string
	config.EtcdConfig
}

// StateRepository implements StateRepository
type StateRepository struct {
	server *embed.Etcd
	client *clientv3.Client
	config *Config
}

// NewStateRepository :
func NewStateRepository(config *Config) (*StateRepository, error) {

	r := &StateRepository{
		config: config,
	}

	err := r.setupEtcdServer()
	if err != nil {
		return nil, fmt.Errorf("Error setting up etcd server: %s", err.Error())
	}

	err = r.setupEtcdClient()
	if err != nil {
		return nil, fmt.Errorf("Error setting up etcd client: %s", err.Error())
	}

	return r, nil
}

// Name returns the name identifying the state repository.
func (r *StateRepository) Name() string {
	return "etcd"
}

type transaction struct {
	txn clientv3.Txn
}

func (t transaction) Commit() (interface{}, error) {
	out, err := t.txn.Commit()
	return out, err
}

// Transaction
func (r *StateRepository) Transaction(ctx context.Context) state.Transaction {
	return transaction{r.client.Txn(ctx)}
}

func (r *StateRepository) setupEtcdClient() error {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:        r.config.InitialAdvertiseClientURLs,
		AutoSyncInterval: time.Second * 5,
		DialTimeout:      5 * time.Second,
	})
	if err != nil {
		return err
	}

	r.client = etcdClient

	return nil
}

func (r *StateRepository) setupEtcdServer() error {

	cfg := embed.NewConfig()

	// Advertise peer URLs
	apURLs, err := parseUrls(r.config.InitialAdvertisePeerURLs)
	if err != nil {
		return err
	}

	// Listen peer URLs
	lpURLs, err := parseUrls(r.config.ListenPeerURLs)
	if err != nil {
		return err
	}

	// Advertise client URLs
	acURLs, err := parseUrls(r.config.InitialAdvertiseClientURLs)
	if err != nil {
		return err
	}

	// Listen client URLs
	lcURLs, err := parseUrls(r.config.ListenClientURLs)
	if err != nil {
		return err
	}

	cfg.Name = r.config.Name
	cfg.Dir = path.Join(r.config.DataDir, "/etcd")
	cfg.WalDir = path.Join(r.config.DataDir, "/etcd", "/wal")
	cfg.Logger = "zap"

	cfg.APUrls = apURLs
	cfg.LPUrls = lpURLs
	cfg.ACUrls = acURLs
	cfg.LCUrls = lcURLs

	cfg.LogOutputs = []string{"stderr", path.Join(r.config.DataDir, "/etcd.log")}
	cfg.LogLevel = strings.ToLower(r.config.LogLevel)

	etcdServer, err := embed.StartEtcd(cfg)
	if err != nil {
		return err
	}

	r.server = etcdServer

	return nil
}

func strToPtr(s string) *string {
	return &s
}

func resourceKey(resourceType, resourceID string) string {
	key := fmt.Sprintf("%s/%s/%s/%s", defaultPrefix, resourceType, defaultNamespace, resourceID)
	return key
}

func encodeValue(in interface{}) string {
	encoded, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	return string(encoded)
}

func decodeValue(data []byte, out interface{}) error {
	return json.Unmarshal(data, out)
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
