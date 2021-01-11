package auth

// Token implements the acl.Token interface
type Token struct {
	privileged bool
	policies   []string
}

// NewToken :
func NewToken(privileged bool, policies []string) *Token {
	return &Token{privileged, policies}
}

// Policies returns a slice of policies associated with
// the token.
func (t *Token) Policies() []string {
	return t.policies
}

// IsPrivileged returns true if the token has privileged access,
// and false otherwise.
func (t *Token) IsPrivileged() bool {
	return t.privileged
}
