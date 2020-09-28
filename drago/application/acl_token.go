package application

import (
	"context"
	"errors"

	structs "github.com/seashell/drago/drago/application/structs"
	domain "github.com/seashell/drago/drago/domain"
	uuid "github.com/seashell/drago/pkg/uuid"
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

// ACLTokenService ...
type ACLTokenService interface {
	List(context.Context, *structs.ACLTokenListInput) (*structs.ACLTokenListOutput, error)
	FindBySecret(context.Context, *structs.ACLTokenGetInput) (*structs.ACLTokenGetOutput, error)
	GetByID(context.Context, *structs.ACLTokenGetInput) (*structs.ACLTokenGetOutput, error)
	Create(context.Context, *structs.ACLTokenCreateInput) (*structs.ACLTokenCreateOutput, error)
	Delete(context.Context, *structs.ACLTokenDeleteInput) (*structs.ACLTokenDeleteOutput, error)
}

type aclTokenService struct {
	repo domain.ACLTokenRepository
}

// NewACLTokenService ...
func NewACLTokenService(tr domain.ACLTokenRepository) ACLTokenService {
	return &aclTokenService{
		repo: tr,
	}
}

// GetByID returns a Token entity by ID
func (s *aclTokenService) GetByID(ctx context.Context, in *structs.ACLTokenGetInput) (*structs.ACLTokenGetOutput, error) {

	t, err := s.repo.GetByID(ctx, in.ID)
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

// FindBySecret returns a Token entity by ID
func (s *aclTokenService) FindBySecret(ctx context.Context, in *structs.ACLTokenGetInput) (*structs.ACLTokenGetOutput, error) {

	t, err := s.repo.FindBySecret(ctx, in.Secret)
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

// Create creates a new Token entity
func (s *aclTokenService) Create(ctx context.Context, in *structs.ACLTokenCreateInput) (*structs.ACLTokenCreateOutput, error) {

	if in.Type != domain.ACLTokenTypeClient && in.Type != domain.ACLTokenTypeManagement {
		return nil, structs.NewError(ErrInvalidTokenType, in.Type)
	}

	if in.Type == domain.ACLTokenTypeManagement && !(in.Policies == nil || len(in.Policies) == 0) {
		return nil, structs.NewError(ErrInvalidTokenPolicies, in.Policies)
	}

	id, err := s.repo.Create(ctx, &domain.ACLToken{
		Name:     in.Name,
		Type:     in.Type,
		Policies: in.Policies,
		Secret:   uuid.Generate(),
	})
	if err != nil {
		return nil, err
	}

	t, err := s.repo.GetByID(ctx, *id)
	if err != nil {
		return nil, err
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
	_, err := s.repo.DeleteByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	return &structs.ACLTokenDeleteOutput{}, nil
}

// List retrieves all token entities in the repository
func (s *aclTokenService) List(ctx context.Context, in *structs.ACLTokenListInput) (*structs.ACLTokenListOutput, error) {

	tokens, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	items := []*structs.ACLTokenListItem{}
	for _, t := range tokens {
		items = append(items, &structs.ACLTokenListItem{
			ID:        t.ID,
			Name:      t.Name,
			Type:      t.Type,
			Policies:  t.Policies,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		})
	}

	out := &structs.ACLTokenListOutput{
		Items: items,
	}

	return out, nil
}
