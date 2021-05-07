package drago

import (
	"context"
	"time"

	structs "github.com/seashell/drago/drago/structs"
)

const (
	ACLPolicyList  = "list"
	ACLPolicyRead  = "read"
	ACLPolicyWrite = "write"
)

func (s *ACLService) GetPolicy(args *structs.ACLPolicySpecificRequest, out *structs.SingleACLPolicyResponse) error {

	if !s.config.ACL.Enabled {
		return structs.ErrACLDisabled
	}

	ctx := context.TODO()

	// Check if authorized
	if !s.config.ACL.Enabled {
		if err := s.authHandler.Authorize(ctx, args.AuthToken, "policy", args.Name, ACLPolicyRead); err != nil {
			return structs.ErrPermissionDenied
		}
	}

	p, err := s.state.ACLPolicyByName(ctx, args.Name)
	if err != nil {
		return structs.ErrNotFound
	}

	out.ACLPolicy = p

	return nil
}

func (s *ACLService) UpsertPolicy(args *structs.ACLPolicyUpsertRequest, out *structs.GenericResponse) error {

	if !s.config.ACL.Enabled {
		return structs.ErrACLDisabled
	}

	ctx := context.TODO()

	// Check if authorized
	if err := s.authHandler.Authorize(ctx, args.AuthToken, "policy", "", ACLPolicyWrite); err != nil {
		return structs.ErrPermissionDenied
	}

	p := args.ACLPolicy

	err := p.Validate()
	if err != nil {
		return structs.NewError(structs.ErrInvalidInput, err)
	}

	old, err := s.state.ACLPolicyByName(ctx, p.Name)
	if err != nil {
		p.CreatedAt = time.Now()
	} else {
		p = old.Merge(p)
	}

	p.UpdatedAt = time.Now()

	err = s.state.UpsertACLPolicy(ctx, p)
	if err != nil {
		return structs.ErrInternal
	}

	return nil
}

func (s *ACLService) DeletePolicies(args *structs.ACLPolicyDeleteRequest, out *structs.GenericResponse) error {

	if !s.config.ACL.Enabled {
		return structs.ErrACLDisabled
	}

	ctx := context.TODO()

	// Check if authorized
	if err := s.authHandler.Authorize(ctx, args.AuthToken, "policy", "", ACLPolicyWrite); err != nil {
		return structs.ErrPermissionDenied
	}

	err := s.state.DeleteACLPolicies(ctx, args.Names)
	if err != nil {
		return structs.ErrInternal
	}

	return nil
}

func (s *ACLService) ListPolicies(args *structs.ACLPolicyListRequest, out *structs.ACLPolicyListResponse) error {

	if !s.config.ACL.Enabled {
		return structs.ErrACLDisabled
	}

	ctx := context.TODO()

	// Check if authorized
	if err := s.authHandler.Authorize(ctx, args.AuthToken, "policy", "", ACLPolicyList); err != nil {
		return structs.ErrPermissionDenied
	}

	policies, err := s.state.ACLPolicies(ctx)
	if err != nil {
		return structs.ErrInternal
	}

	out.Items = nil

	for _, p := range policies {
		out.Items = append(out.Items, p.Stub())
	}

	return nil
}
