package storage

import (
	"github.com/seashell/drago/api"
)

// StateRepository :
type StateRpository interface {
	Name() string
	GetHostSettings() (*api.HostSettings, error)
	PutHostSettings(*api.HostSettings) error
}
