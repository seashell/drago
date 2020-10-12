package application

import (
	"context"
	"errors"

	structs "github.com/seashell/drago/drago/application/structs"
	domain "github.com/seashell/drago/drago/domain"
)

const (
	ACLPolicyList  = "list"
	ACLPolicyRead  = "read"
	ACLPolicyWrite = "write"
)

const (
	errPolicyNotFound = "policy not found"
)

var (
	// ErrPolicyNotFound ...
	ErrPolicyNotFound = errors.New(errPolicyNotFound)
)

type aclPolicyService struct {
	config *Config
}

// NewACLPolicyService ...
func NewACLPolicyService(config *Config) ACLPolicyService {
	return &aclPolicyService{config}
}

func (s *aclPolicyService) GetByName(ctx context.Context, in *structs.ACLPolicyGetInput) (*structs.ACLPolicyGetOutput, error) {

	// Check if authorized
	if err := s.authorize(ctx, in.Subject, in.Name, ACLPolicyRead); err != nil {
		return nil, err
	}

	p, err := s.config.ACLPolicyRepository.GetByName(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	out := &structs.ACLPolicyGetOutput{
		ACLPolicy: structs.ACLPolicy{
			Name:      p.Name,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		},
	}

	return out, nil
}

func (s *aclPolicyService) Upsert(ctx context.Context, in *structs.ACLPolicyUpsertInput) (*structs.ACLPolicyUpsertOutput, error) {

	// Check if authorized
	if err := s.authorize(ctx, in.Subject, "", ACLPolicyWrite); err != nil {
		return nil, err
	}

	_, err := s.config.ACLPolicyRepository.Upsert(ctx, &domain.ACLPolicy{
		Name: in.Name,
	})
	if err != nil {
		return nil, err
	}
	return &structs.ACLPolicyUpsertOutput{}, nil
}

func (s *aclPolicyService) Delete(ctx context.Context, in *structs.ACLPolicyDeleteInput) (*structs.ACLPolicyDeleteOutput, error) {

	// Check if authorized
	if err := s.authorize(ctx, in.Subject, "", ACLPolicyWrite); err != nil {
		return nil, err
	}

	_, err := s.config.ACLPolicyRepository.DeleteByName(ctx, in.Name)
	if err != nil {
		return nil, err
	}
	return &structs.ACLPolicyDeleteOutput{}, nil
}

func (s *aclPolicyService) List(ctx context.Context, in *structs.ACLPolicyListInput) (*structs.ACLPolicyListOutput, error) {

	// Check if authorized
	if err := s.authorize(ctx, in.Subject, "", ACLPolicyList); err != nil {
		return nil, err
	}

	policies, err := s.config.ACLPolicyRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	out := &structs.ACLPolicyListOutput{
		Items: []*structs.ACLPolicyListItem{},
	}
	for _, p := range policies {
		out.Items = append(out.Items, &structs.ACLPolicyListItem{
			Name:      p.Name,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}

	return out, nil
}

func (s *aclPolicyService) authorize(ctx context.Context, sub, id, op string) error {
	if s.config.ACLEnabled {
		if err := s.config.AuthHandler.Authorize(ctx, sub, ResourceACLPolicy, id, op); err != nil {
			return err
		}
	}
	return nil
}
