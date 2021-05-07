package drago

import (
	"context"
	"time"

	structs "github.com/seashell/drago/drago/structs"
	uuid "github.com/seashell/drago/pkg/uuid"
)

const (
	ACLTokenList  = "list"
	ACLTokenRead  = "read"
	ACLTokenWrite = "write"
)

// GetToken returns a Token entity by ID
func (s *ACLService) GetToken(args *structs.ACLTokenSpecificRequest, out *structs.SingleACLTokenResponse) error {

	if !s.config.ACL.Enabled {
		return structs.ErrACLDisabled
	}

	ctx := context.TODO()

	// Check if authorized
	if err := s.authHandler.Authorize(ctx, args.AuthToken, "token", args.ACLTokenID, ACLTokenRead); err != nil {
		return structs.ErrPermissionDenied
	}

	t, err := s.state.ACLTokenByID(ctx, args.ACLTokenID)
	if err != nil {
		return structs.ErrNotFound
	}

	out.ACLToken = t

	return nil
}

// UpsertToken creates or updates a new Token entity
func (s *ACLService) UpsertToken(args *structs.ACLTokenUpsertRequest, out *structs.ACLTokenUpsertResponse) error {

	if !s.config.ACL.Enabled {
		return structs.ErrACLDisabled
	}

	ctx := context.TODO()

	// Check if authorized
	if err := s.authHandler.Authorize(ctx, args.AuthToken, "token", "", ACLTokenWrite); err != nil {
		return structs.ErrPermissionDenied
	}

	t := args.ACLToken

	err := t.Validate()
	if err != nil {
		return structs.ErrInvalidInput
	}

	// Generate a new ID and secret if the token is new
	if t.ID == "" {
		t.ID = uuid.Generate()
		t.Secret = uuid.Generate()
		t.CreatedAt = time.Now()
	} else {
		old, err := s.state.ACLTokenByID(ctx, t.ID)
		if err != nil {
			return structs.ErrNotFound
		}
		t = old.Merge(t)
	}

	t.UpdatedAt = time.Now()

	err = s.state.UpsertACLToken(ctx, t)
	if err != nil {
		return structs.ErrInternal
	}

	out.ACLToken = t

	return nil
}

// DeleteToken deletes a token entity from the repository
func (s *ACLService) DeleteToken(args *structs.ACLTokenDeleteRequest, out *structs.GenericResponse) error {

	if !s.config.ACL.Enabled {
		return structs.ErrACLDisabled
	}

	ctx := context.TODO()

	// Check if authorized
	if err := s.authHandler.Authorize(ctx, args.AuthToken, "token", "", ACLTokenWrite); err != nil {
		return structs.ErrPermissionDenied
	}

	err := s.state.DeleteACLTokens(ctx, args.ACLTokenIDs)
	if err != nil {
		return structs.ErrInternal
	}

	return nil
}

// ListTokens retrieves all token entities in the repository
func (s *ACLService) ListTokens(args *structs.ACLTokenListRequest, out *structs.ACLTokenListResponse) error {

	if !s.config.ACL.Enabled {
		return structs.ErrACLDisabled
	}

	ctx := context.TODO()

	// Check if authorized
	if err := s.authHandler.Authorize(ctx, args.AuthToken, "token", "", ACLTokenList); err != nil {
		return structs.ErrPermissionDenied
	}

	tokens, err := s.state.ACLTokens(ctx)
	if err != nil {
		return structs.ErrInternal
	}

	out.Items = nil

	for _, t := range tokens {
		out.Items = append(out.Items, t.Stub())
	}

	return nil
}
