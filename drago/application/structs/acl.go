package structs

// ACLBootstrapInput :
type ACLBootstrapInput struct{}

// ACLBootstrapOutput :
type ACLBootstrapOutput struct {
	ACLToken
}

// ACLResolveTokenInput :
type ACLResolveTokenInput struct {
	Secret string
}

// ACLResolveTokenOutput :
type ACLResolveTokenOutput struct {
	ACLToken
}
