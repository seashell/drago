package drago

import (
	auth "github.com/seashell/drago/drago/auth"
	state "github.com/seashell/drago/drago/state"
	structs "github.com/seashell/drago/drago/structs"
)

// StatusService is used to check on server status
type StatusService struct {
	config      *Config
	state       state.Repository
	authHandler auth.AuthorizationHandler
}

// NewStatusService ...
func NewStatusService(config *Config, state state.Repository, authHandler auth.AuthorizationHandler) *StatusService {
	return &StatusService{
		config:      config,
		state:       state,
		authHandler: authHandler,
	}
}

// Ping is used to check for connectivity
func (s *StatusService) Ping(args structs.GenericRequest, out *structs.GenericResponse) error {
	return nil
}

// Version returns the version of the server
func (s *StatusService) Version(in structs.GenericRequest, out *structs.StatusVersionResponse) error {
	out.Version = s.config.Version.VersionNumber()
	return nil
}
