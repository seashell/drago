package application

import (
	"context"
	"errors"

	structs "github.com/seashell/drago/drago/application/structs"
	domain "github.com/seashell/drago/drago/domain"
)

const (
	errPermissionDenied    = "permission denied"
	errAlreadyBootstrapped = "acl already bootstrapped"
	errInvalidSecret       = "invalid secret"
)

var (
	// ErrAlreadyBootstrapped ...
	ErrAlreadyBootstrapped = errors.New(errAlreadyBootstrapped)
	// ErrPermissionDenied ...
	ErrPermissionDenied = errors.New(errPermissionDenied)
	// ErrInvalidSecret ...
	ErrInvalidSecret = errors.New(errInvalidSecret)
)

// ACLService ...
type ACLService interface {
	Bootstrap(ctx context.Context, in *structs.ACLBootstrapInput) (*structs.ACLBootstrapOutput, error)
	ResolveToken(ctx context.Context, in *structs.ACLResolveTokenInput) (*structs.ACLResolveTokenOutput, error)
}

type aclService struct {
	stateRepo  domain.ACLStateRepository
	tokenRepo  domain.ACLTokenRepository
	policyRepo domain.ACLPolicyRepository
}

// NewACLService ...
func NewACLService(sr domain.ACLStateRepository, tr domain.ACLTokenRepository, pr domain.ACLPolicyRepository) ACLService {
	return &aclService{
		stateRepo:  sr,
		tokenRepo:  tr,
		policyRepo: pr,
	}
}

// Bootstrap
func (s *aclService) Bootstrap(ctx context.Context, in *structs.ACLBootstrapInput) (*structs.ACLBootstrapOutput, error) {

	if s.isBootstrapped(ctx) {
		return nil, structs.NewError(ErrAlreadyBootstrapped)
	}

	id, err := s.tokenRepo.Create(ctx, &domain.ACLToken{
		Name:     "Root Token",
		Type:     domain.ACLTokenTypeManagement,
		Policies: nil,
	})
	if err != nil {
		return nil, err
	}

	t, err := s.tokenRepo.GetByID(ctx, *id)
	if err != nil {
		return nil, err
	}

	// Update ACL state
	state := s.state()
	state.RootTokenID = t.ID
	state.RootTokenResetIndex++

	err = s.stateRepo.Set(ctx, state)
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
	t, err := s.tokenRepo.FindBySecret(ctx, in.Secret)
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
	state, err := s.stateRepo.Get(ctx)
	if errors.Is(err, domain.ErrNotFound) {
		state = &domain.ACLState{}
		s.stateRepo.Set(ctx, state)
	}
	return state
}
