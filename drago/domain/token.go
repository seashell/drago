package domain

import "time"

const (
	// Client token type
	TokenTypeClient = "client"
	// Management token type
	TokenTypeManagement = "management"
)

// Token :
type Token struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Name      *string   `json:"name"`
	Secret    string    `json:"secret"`
	Poilicies []string  `json:"policies"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TokenRepository : Token repository interface
type TokenRepository interface {
	GetByID(string) (*Token, error)
	Create(n *Token) (*string, error)
	Update(n *Token) (*string, error)
	DeleteByID(string) (*string, error)
	FindAll(pageInfo PageInfo) ([]*Token, *Page, error)
}

func (t *Token) Merge(b *Token) *Token {

	if b == nil {
		return t
	}

	result := *t

	return &result
}
