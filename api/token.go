package api

import "context"

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

func (t *Tokens) Create(ctx context.Context, token *Token) (*string, error) {
	receiver := struct {
		Token *string `json:"secret"`
	}{}

	err := t.client.createResource(tokensPath, token, &receiver)
	if err != nil {
		return nil, err
	}

	return receiver.Token, nil
}

func (t *Tokens) Update(ctx context.Context, token *Token) (*string, error) {
	return nil, nil
}

func (t *Tokens) Delete(ctx context.Context, token *Token) error {
	return nil
}

func (t *Tokens) List(ctx context.Context) ([]*Token, error) {
	return []*Token{}, nil
}
