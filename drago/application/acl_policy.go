package application

import (
	"context"
	"errors"

	structs "github.com/seashell/drago/drago/application/structs"
	domain "github.com/seashell/drago/drago/domain"
)

const (
	errPolicyNotFound = "policy not found"
)

var (
	// ErrPolicyNotFound ...
	ErrPolicyNotFound = errors.New(errPolicyNotFound)
)

// ACLPolicyService ...
type ACLPolicyService interface {
	List(context.Context, *structs.ACLPolicyListInput) (*structs.ACLPolicyListOutput, error)
	GetByName(context.Context, *structs.ACLPolicyGetInput) (*structs.ACLPolicyGetOutput, error)
	Upsert(context.Context, *structs.ACLPolicyUpsertInput) (*structs.ACLPolicyUpsertOutput, error)
	Delete(context.Context, *structs.ACLPolicyDeleteInput) (*structs.ACLPolicyDeleteOutput, error)
}

type aclPolicyService struct {
	repo domain.ACLPolicyRepository
}

// NewACLPolicyService ...
func NewACLPolicyService(pr domain.ACLPolicyRepository) ACLPolicyService {
	return &aclPolicyService{
		repo: pr,
	}
}

func (s *aclPolicyService) GetByName(ctx context.Context, in *structs.ACLPolicyGetInput) (*structs.ACLPolicyGetOutput, error) {

	t, err := s.repo.GetByName(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	out := &structs.ACLPolicyGetOutput{
		ACLPolicy: structs.ACLPolicy{
			Name:      t.Name,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		},
	}

	return out, nil
}

func (s *aclPolicyService) Upsert(ctx context.Context, in *structs.ACLPolicyUpsertInput) (*structs.ACLPolicyUpsertOutput, error) {

	_, err := s.repo.Upsert(ctx, &domain.ACLPolicy{
		Name: in.Name,
	})
	if err != nil {
		return nil, err
	}

	return &structs.ACLPolicyUpsertOutput{}, nil
}

func (s *aclPolicyService) Delete(ctx context.Context, in *structs.ACLPolicyDeleteInput) (*structs.ACLPolicyDeleteOutput, error) {
	_, err := s.repo.DeleteByName(ctx, in.Name)
	if err != nil {
		return nil, err
	}
	return &structs.ACLPolicyDeleteOutput{}, nil
}

func (s *aclPolicyService) List(ctx context.Context, in *structs.ACLPolicyListInput) (*structs.ACLPolicyListOutput, error) {
	policies, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	items := []*structs.ACLPolicyListItem{}
	for _, p := range policies {
		items = append(items, &structs.ACLPolicyListItem{
			Name:      p.Name,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}

	out := &structs.ACLPolicyListOutput{
		Items: items,
	}

	return out, nil
}
