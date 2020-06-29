package state

import (
	"github.com/seashell/drago/api"
)


type StateDB interface{

	Name()	(string)

	GetHostSettings()					(*api.HostSettings, error)
	PutHostSettings(*api.HostSettings)	(error)
}