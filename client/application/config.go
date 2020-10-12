package application

import (
	"context"

	domain "github.com/seashell/drago/client/domain"
	structs "github.com/seashell/drago/drago/application/structs"
)

// NetworkInterfaceController provides network configuration capabilities.
type NetworkInterfaceController interface {
	CreateInterface(iface *domain.Interface) error
	ListInterfaces() ([]*domain.Interface, error)
	DeleteInterface(name string) error
	DeleteAllInterfaces() error
}

// DragoGateway is a gateway/client for acessing Drago's remote API.
type DragoGateway interface {
	Agent() DragoAgentGateway
}

// DragoAgentGateway is a gateway/client for acessing Drago's remote agent API.
type DragoAgentGateway interface {
	SynchronizeSelf(ctx context.Context, state *structs.HostSynchronizeInput) (*structs.HostSynchronizeOutput, error)
}

// DragoTokenGateway is a gateway/client for acessing Drago's remote tokens API.
type DragoTokenGateway interface {
	// ...
}

// Config contains configurations for the Drago client application
// services.
type Config struct {
	DragoGateway        DragoGateway
	StateRepository     domain.Repository
	InterfaceController NetworkInterfaceController
}
