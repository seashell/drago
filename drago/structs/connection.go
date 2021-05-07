package structs

import (
	"sort"
	"time"
)

// Connection :
type Connection struct {
	ID        string
	NetworkID string
	NodeIDs   []string

	// PeerSettings contains the ID and the configurations to be applied
	// to each of the connected interfaces.
	PeerSettings map[string]*PeerSettings

	// If the connection is going from a NAT-ed peer to a public peer,
	// the node behind the NAT must regularly send an outgoing ping to
	// keep the bidirectional connection alive in the NAT router's
	// connection table.
	PersistentKeepalive *int

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Validate :
func (c *Connection) Validate() error {
	return nil
}

// ConnectedInterfaceIDs :
func (c *Connection) ConnectedInterfaceIDs() []string {
	ids := []string{}
	for k, _ := range c.PeerSettings {
		ids = append(ids, k)
	}
	sort.Strings(ids)
	return ids
}

// PeerSettingsByInterfaceID :
func (c *Connection) PeerSettingsByInterfaceID(s string) *PeerSettings {
	return c.PeerSettings[s]
}

// OtherPeerSettingsByInterfaceID : given the ID of one of the connected interfaces,
// returns the settings for the peer/interface at the other end of the connection.
func (c *Connection) OtherPeerSettingsByInterfaceID(s string) *PeerSettings {
	for ifaceID, settings := range c.PeerSettings {
		if ifaceID != s {
			return settings
		}
	}
	return nil
}

// ConnectsInterfaces : checks whether a Connection connects two
// interfaces whose indices are passed as arguments.
func (c *Connection) ConnectsInterfaces(a, b string) bool {
	if c.ConnectsInterface(a) && c.ConnectsInterface(b) {
		return true
	}
	return false
}

// ConnectsInterface : checks whether a connection connects
// an interface whose index is passed as argument.
func (c *Connection) ConnectsInterface(s string) bool {
	if _, ok := c.PeerSettings[s]; ok {
		return true
	}
	return false
}

// Merge :
func (c *Connection) Merge(in *Connection) *Connection {

	result := *c

	if in.ID != "" {
		result.ID = in.ID
	}
	if in.PeerSettings != nil {
		if result.PeerSettings == nil {
			result.PeerSettings = in.PeerSettings
		} else {
			for k := range in.PeerSettings {
				if _, ok := result.PeerSettings[k]; ok {
					result.PeerSettings[k] = result.PeerSettings[k].Merge(in.PeerSettings[k])
				} else {
					result.PeerSettings[k] = in.PeerSettings[k]
				}
			}
		}
	}
	if in.PersistentKeepalive != nil {
		result.PersistentKeepalive = in.PersistentKeepalive
	}

	return &result
}

// Stub :
func (c *Connection) Stub() *ConnectionListStub {

	peers := []string{}
	for k := range c.PeerSettings {
		peers = append(peers, k)
	}

	return &ConnectionListStub{
		ID:                  c.ID,
		NetworkID:           c.NetworkID,
		NodeIDs:             c.NodeIDs,
		Peers:               peers,
		PeerSettings:        c.PeerSettings,
		PersistentKeepalive: c.PersistentKeepalive,
		BytesTransferred:    0,
		CreatedAt:           c.CreatedAt,
		UpdatedAt:           c.UpdatedAt,
	}
}

// ConnectionListStub :
type ConnectionListStub struct {
	ID                  string
	NetworkID           string
	NodeIDs             []string
	Peers               []string
	PeerSettings        map[string]*PeerSettings
	PersistentKeepalive *int
	BytesTransferred    uint64
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// PeerSettings :
type PeerSettings struct {
	InterfaceID  string
	RoutingRules *RoutingRules
}

// Merge :
func (r *PeerSettings) Merge(in *PeerSettings) *PeerSettings {
	result := *r
	if in.InterfaceID != "" {
		result.InterfaceID = in.InterfaceID
	}
	if in.RoutingRules != nil {
		result.RoutingRules = r.RoutingRules.Merge(in.RoutingRules)
	}
	return &result
}

// RoutingRules :
type RoutingRules struct {
	// AllowedIPs defines the IP ranges for which traffic will be routed/accepted.
	// Example: If AllowedIPs = [192.0.2.3/32, 192.168.1.1/24], the node
	// will accept traffic for itself (192.0.2.3/32), and for all nodes in the
	// local network (192.168.1.1/24).
	AllowedIPs []string
}

// Merge :
func (r *RoutingRules) Merge(in *RoutingRules) *RoutingRules {
	result := *r
	if in.AllowedIPs != nil {
		result.AllowedIPs = in.AllowedIPs
	}
	return &result
}

// ConnectionSpecificRequest :
type ConnectionSpecificRequest struct {
	ConnectionID string

	QueryOptions
}

// SingleConnectionResponse :
type SingleConnectionResponse struct {
	Connection *Connection

	Response
}

// ConnectionUpsertRequest :
type ConnectionUpsertRequest struct {
	Connection *Connection

	WriteRequest
}

// ConnectionDeleteRequest :
type ConnectionDeleteRequest struct {
	ConnectionIDs []string

	WriteRequest
}

// ConnectionListRequest :
type ConnectionListRequest struct {
	InterfaceID string
	NodeID      string
	NetworkID   string

	QueryOptions
}

// ConnectionListResponse :
type ConnectionListResponse struct {
	Items []*ConnectionListStub

	Response
}
