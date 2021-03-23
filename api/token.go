package api

import (
	"context"

	"github.com/seashell/drago/drago/structs"
)

const (
	tokensPath = "/api/tokens"
)

type Token struct {
	ID        *string   `json:"id,omitempty"`
	Secret    *string   `json:"secret,omitempty"`
	Name      *string   `json:"name,omitempty"`
	Type      *string   `json:"type,omitempty"`
	Policies  *[]string `json:"policies,omitempty"`
	IssuedAt  *int64    `json:"issuedAt,omitempty"`
	ExpiresAt *int64    `json:"expiresAt,omitempty"`
	NotBefore *int64    `json:"notBefore,omitempty"`
}

// Tokens is a client to the tokens API
type Tokens struct {
	client *Client
}

// Hosts returns a handle on the hosts endpoints.
func (c *Client) Tokens() *Tokens {
	return &Tokens{client: c}
}

func (t *Tokens) Get(ctx context.Context) (*Token, error) {
	return &Token{}, nil
}

func (t *Tokens) Create(req *structs.ACLTokenUpsertRequest) (*structs.SingleACLTokenResponse, error) {

	var resp structs.SingleACLTokenResponse
	err := t.client.createResource(tokensPath, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (t *Tokens) Update(token *Token) (*string, error) {
	return nil, nil
}

func (t *Tokens) Delete(token *Token) error {
	return nil
}

func (t *Tokens) List(ctx context.Context) ([]*Token, error) {
	return []*Token{}, nil
}
