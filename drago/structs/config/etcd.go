package config

// EtcdConfig : see https://etcd.io/docs/v3.4.0/op-guide/configuration
type EtcdConfig struct {
	// Human-readable name of this node in the etcd cluster
	Name string

	// List of URLs to listen on for peer traffic. This flag tells the
	// etcd to accept incoming requests from its peers on the specified
	// scheme://IP:port combinations. Defaults to [“http://localhost:2380”].
	ListenPeerURLs []string

	// List of URLs to listen on for client traffic. This flag tells the
	// etcd to accept incoming requests from the clients on the specified
	// scheme://IP:port combinations. Defaults to [“http://localhost:2379”].
	ListenClientURLs []string

	// List of this member’s peer URLs to advertise to the rest of the cluster.
	// These addresses are used for communicating etcd data around the cluster.
	// At least one must be routable to all cluster members. These URLs can
	// contain domain names. Defaults to [“http://localhost:2380”].
	InitialAdvertisePeerURLs []string

	// List of this member’s client URLs to advertise to the rest of the cluster.
	// These URLs can contain domain names. Defaults to [“http://localhost:2379”].
	InitialAdvertiseClientURLs []string

	// Initial cluster configuration for bootstrapping. Defaults to [“http://localhost:2380”].
	// The key is the value of Name for each node provided. The default uses 'default' for
	// the key because this is the default for the --name flag
	InitialCluster []string

	// Initial cluster state (“new” or “existing”). Set to new for all members
	// present during initial static or DNS bootstrapping. If this option is set
	// to existing, etcd will attempt to join the existing cluster. If the wrong
	// value is set, etcd will attempt to start but fail safely.
	InitialClusterState string

	// Defines whether this node should run on proxy mode. If enabled, this means that
	// the node will simply forward requests to an already existing etcd cluster, without
	// actually joining it. See https://etcd.io/docs/v2/proxy/.
	ProxyModeEnabled bool

	// Comma-separated white list of origins for CORS (cross-origin
	// resource sharing).
	CORS string
}

// DefaultEtcdConfig returns the canonical defaults for the Drago
// `etcd` configuration.
func DefaultEtcdConfig() *EtcdConfig {
	return &EtcdConfig{
		Name:                       "default",
		ListenPeerURLs:             []string{"http://localhost:2380"},
		ListenClientURLs:           []string{"http://localhost:2379"},
		InitialAdvertisePeerURLs:   []string{"http://localhost:2380"},
		InitialAdvertiseClientURLs: []string{"http://localhost:2379"},
		InitialCluster:             []string{"http://localhost:2380"},
		InitialClusterState:        "new",
		ProxyModeEnabled:           false,
		CORS:                       "",
	}
}
