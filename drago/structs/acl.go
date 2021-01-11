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
	WriteRequest
	ResetIndex uint64
}

// ResolveACLTokenRequest :
type ResolveACLTokenRequest struct {
	QueryOptions
	Secret string
}

// ResolveACLTokenResponse :
type ResolveACLTokenResponse struct {
	Response
	ACLToken *ACLToken
}
