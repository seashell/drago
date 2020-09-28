package application

import domain "github.com/seashell/drago/drago/domain"

type Drago struct {
	aclStateRepository  domain.ACLStateRepository
	aclTokenRepository  domain.ACLPolicyRepository
	aclPolicyRepository domain.ACLPolicyRepository
}

func New() (*Drago, error) {
	return nil, nil
}
