package state

import (
	"github.com/seashell/drago/api"
)

// StateDB :
type StateDB interface {
	Name() string
	GetHostSettings() (*api.HostSettings, error)
	PutHostSettings(*api.HostSettings) error
}
