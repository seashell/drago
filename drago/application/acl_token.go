package application

import (
	"context"
	"errors"

	structs "github.com/seashell/drago/drago/application/structs"
	domain "github.com/seashell/drago/drago/domain"
	uuid "github.com/seashell/drago/pkg/uuid"
)

const (
	ACLTokenList  = "list"
	ACLTokenRead  = "read"
	ACLTokenWrite = "write"
)

const (
	errTokenNotFound        = "token not found"
	errInvalidTokenType     = "invalid token type"
	errInvalidTokenPolicies = "invalid token policies"
)

var (
	// ErrTokenNotFound ...
	ErrTokenNotFound = errors.New(errTokenNotFound)

	// ErrInvalidTokenType ...
	ErrInvalidTokenType = errors.New(errInvalidTokenType)

	// ErrInvalidTokenPolicies ...
	ErrInvalidTokenPolicies = errors.New(errInvalidTokenPolicies)
)

type aclTokenService struct {
	config *Config
}

// NewACLTokenService ...
func NewACLTokenService(config *Config) ACLTokenService {
	return &aclTokenService{config}
}

// GetByID returns a Token entity by ID
func (s *aclTokenService) GetByID(ctx context.Context, in *structs.ACLTokenGetInput) (*structs.ACLTokenGetOutput, error) {

	// Check if authorized
	if err := s.authorize(ctx, in.Subject, in.ID, ACLTokenRead); err != nil {
		return nil, ErrUnauthorized
	}

	t, err := s.config.ACLTokenRepository.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	out := &structs.ACLTokenGetOutput{
		ACLToken: structs.ACLToken{
			ID:        t.ID,
			Secret:    t.Secret,
			Type:      t.Type,
			Name:      t.Name,
			Policies:  t.Policies,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		},
	}

	return out, nil
}

// GetBySecret returns a Token entity by ID
func (s *aclTokenService) GetBySecret(ctx context.Context, in *structs.ACLTokenGetInput) (*structs.ACLTokenGetOutput, error) {

	t, err := s.config.ACLTokenRepository.FindBySecret(ctx, in.Secret)
	if err != nil {
		return nil, err
	}

	// Check if authorized
	if err := s.authorize(ctx, in.Subject, t.ID, ACLTokenRead); err != nil {
		return nil, ErrUnauthorized
	}

	out := &structs.ACLTokenGetOutput{
		ACLToken: structs.ACLToken{
			ID:        t.ID,
			Secret:    t.Secret,
			Type:      t.Type,
			Name:      t.Name,
			Policies:  t.Policies,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		},
	}

	return out, nil
}

// Create creates a new Token entity
func (s *aclTokenService) Create(ctx context.Context, in *structs.ACLTokenCreateInput) (*structs.ACLTokenCreateOutput, error) {

	// Check if authorized
	if err := s.authorize(ctx, in.Subject, "", ACLTokenWrite); err != nil {
		return nil, ErrUnauthorized
	}

	if in.Type != domain.ACLTokenTypeClient && in.Type != domain.ACLTokenTypeManagement {
		return nil, structs.NewError(ErrInvalidTokenType, in.Type)
	}

	if in.Type == domain.ACLTokenTypeManagement && !(in.Policies == nil || len(in.Policies) == 0) {
		return nil, structs.NewError(ErrInvalidTokenPolicies, in.Policies)
	}

	id, err := s.config.ACLTokenRepository.Create(ctx, &domain.ACLToken{
		Name:     in.Name,
		Type:     in.Type,
		Policies: in.Policies,
		Secret:   uuid.Generate(),
	})
	if err != nil {
		return nil, err
	}

	t, err := s.config.ACLTokenRepository.GetByID(ctx, *id)
	if err != nil {
		return nil, ErrUnauthorized
	}

	out := &structs.ACLTokenCreateOutput{
		ACLToken: structs.ACLToken{
			ID:        t.ID,
			Secret:    t.Secret,
			Name:      t.Name,
			Type:      t.Type,
			Policies:  t.Policies,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		},
	}

	return out, nil
}

// Delete deletes a token entity from the repository
func (s *aclTokenService) Delete(ctx context.Context, in *structs.ACLTokenDeleteInput) (*structs.ACLTokenDeleteOutput, error) {

	// Check if authorized
	if err := s.authorize(ctx, in.Subject, "", ACLTokenWrite); err != nil {
		return nil, ErrUnauthorized
	}

	_, err := s.config.ACLTokenRepository.DeleteByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	return &structs.ACLTokenDeleteOutput{}, nil
}

// List retrieves all token entities in the repository
func (s *aclTokenService) List(ctx context.Context, in *structs.ACLTokenListInput) (*structs.ACLTokenListOutput, error) {

	// Check if authorized
	if err := s.authorize(ctx, in.Subject, "", ACLTokenList); err != nil {
		return nil, ErrUnauthorized
	}

	tokens, err := s.config.ACLTokenRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	out := &structs.ACLTokenListOutput{
		Items: []*structs.ACLTokenListItem{},
	}

	for _, t := range tokens {
		out.Items = append(out.Items, &structs.ACLTokenListItem{
			ID:        t.ID,
			Name:      t.Name,
			Type:      t.Type,
			Policies:  t.Policies,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		})
	}

	return out, nil
}

func (s *aclTokenService) authorize(ctx context.Context, sub, id, op string) error {
	if s.config.ACLEnabled {
		if err := s.config.AuthHandler.Authorize(ctx, sub, ResourceACLToken, id, op); err != nil {
			return err
		}
	}
	return nil
}
