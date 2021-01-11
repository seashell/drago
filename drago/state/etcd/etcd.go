package etcd

import (
	"encoding/json"
	"fmt"

	"go.etcd.io/etcd/clientv3"
)

const (
	defaultNamespace = "default"
	defaultPrefix    = "/registry"

	resourceTypeACLPolicy = "policy"
	resourceTypeACLToken  = "token"
	resourceTypeNetwork   = "network"
)

// StateRepository implements StateRepository
type StateRepository struct {
	client *clientv3.Client
}

// NewStateRepository :
func NewStateRepository() (*StateRepository, error) {

	var client *clientv3.Client
	// etcdClient, err := clientv3.New(clientv3.Config{
	// 	Endpoints:        s.config.Etcd.InitialAdvertiseClientURLs,
	// 	AutoSyncInterval: time.Second * 5,
	// 	DialTimeout:      5 * time.Second,
	// })
	// if err != nil {
	// 	return err
	// }

	r := &StateRepository{client: client}
	return r, nil
}

// Name ...
func (b *StateRepository) Name() string {
	return "etcd"
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
