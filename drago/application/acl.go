package application

import (
	"context"
	"errors"

	structs "github.com/seashell/drago/drago/application/structs"
	domain "github.com/seashell/drago/drago/domain"
)

const (
	// aclBootstrapResetFileName is the name of the file in the data dir containing the reset index.
	aclBootstrapResetFileName = "acl-bootstrap-reset"
)

const (
	errPermissionDenied    = "permission denied"
	errAlreadyBootstrapped = "acl already bootstrapped"
	errInvalidSecret       = "invalid secret"
	errACLDisabled         = "acl disabled"
)

var (
	// ErrAlreadyBootstrapped ...
	ErrAlreadyBootstrapped = errors.New(errAlreadyBootstrapped)
	// ErrPermissionDenied ...
	ErrPermissionDenied = errors.New(errPermissionDenied)
	// ErrInvalidSecret ...
	ErrInvalidSecret = errors.New(errInvalidSecret)
	// ErrACLDisabled ...
	ErrACLDisabled = errors.New(errACLDisabled)
)

type aclService struct {
	config *Config
}

// NewACLService ...
func NewACLService(config *Config) ACLService {
	return &aclService{config}
}

// Bootstrap
func (s *aclService) Bootstrap(ctx context.Context, in *structs.ACLBootstrapInput) (*structs.ACLBootstrapOutput, error) {

	if !s.config.ACLEnabled {
		return nil, structs.NewError(ErrACLDisabled)
	}

	if s.isBootstrapped(ctx) {
		return nil, structs.NewError(ErrAlreadyBootstrapped)
	}

	id, err := s.config.ACLTokenRepository.Create(ctx, &domain.ACLToken{
		Name:     "Root Token",
		Type:     domain.ACLTokenTypeManagement,
		Policies: nil,
	})
	if err != nil {
		return nil, err
	}

	t, err := s.config.ACLTokenRepository.GetByID(ctx, *id)
	if err != nil {
		return nil, err
	}

	// Update ACL state
	state := s.state()
	state.RootTokenID = t.ID
	state.RootTokenResetIndex++

	err = s.config.ACLStateRepository.Set(ctx, state)
	if err != nil {
		return nil, err
	}

	return &structs.ACLBootstrapOutput{
		ACLToken: structs.ACLToken{
			ID:        t.ID,
			Secret:    t.Secret,
			Name:      t.Name,
			Type:      t.Type,
			Policies:  t.Policies,
			CreatedAt: t.CreatedAt,
		},
	}, nil
}

func (s *aclService) ResolveToken(ctx context.Context, in *structs.ACLResolveTokenInput) (*structs.ACLResolveTokenOutput, error) {
	t, err := s.config.ACLTokenRepository.FindBySecret(ctx, in.Secret)
	if err != nil {
		return nil, err
	}

	if t == nil {
		return nil, ErrInvalidSecret
	}

	return &structs.ACLResolveTokenOutput{
		ACLToken: structs.ACLToken{
			ID:        t.ID,
			Secret:    t.Secret,
			Name:      t.Name,
			Type:      t.Type,
			Policies:  t.Policies,
			CreatedAt: t.CreatedAt,
		},
	}, nil
}

func (s *aclService) isBootstrapped(ctx context.Context) bool {
	return s.state().RootTokenID != ""
}

// Lazy state persistence
func (s *aclService) state() *domain.ACLState {
	ctx := context.TODO()
	state, err := s.config.ACLStateRepository.Get(ctx)
	if errors.Is(err, domain.ErrNotFound) {
		state = &domain.ACLState{}
		s.config.ACLStateRepository.Set(ctx, state)
	}
	return state
}
