package structs

// ACLState ...
type ACLState struct {
	RootTokenID         string
	RootTokenResetIndex int
}

// Update ...
func (s *ACLState) Update(rootID string) error {
	s.RootTokenID = rootID
	s.RootTokenResetIndex++
	return nil
}

// ACLBootstrapRequest :
type ACLBootstrapRequest struct {
	ResetIndex uint64

	WriteRequest
}

// ResolveACLTokenRequest :
type ResolveACLTokenRequest struct {
	Secret string

	QueryOptions
}

// ResolveACLTokenResponse :
type ResolveACLTokenResponse struct {
	ACLToken *ACLToken

	Response
}
