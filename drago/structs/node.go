package structs

import "time"

const (
	NodeStatusInit  = "initializing"
	NodeStatusReady = "ready"
	NodeStatusDown  = "down"
)

// Node :
type Node struct {
	ID   string
	Name string
	Meta map[string]string
}

// NodeSpecificRequest :
type NodeSpecificRequest struct {
	QueryOptions
	ID string
}

// SingleNodeResponse :
type SingleNodeResponse struct {
	Response
	Node
}

// NodeRegisterRequest :
type NodeRegisterRequest struct {
	WriteRequest
	Node
}

// NodeUpdateStatusRequest:
type NodeUpdateStatusRequest struct {
	ID     string
	Status string

	Response
}

// NodeUpdateResponse is used to update nodes
type NodeUpdateResponse struct {
	HeartbeatTTL time.Duration

	Response
}

type NodeClientInterfacesResponse struct {
}

type NodeClientPeersResponse struct {
}
