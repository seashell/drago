package drago

import (
	application "github.com/seashell/drago/drago/application"
	inmem "github.com/seashell/drago/drago/inmem"
)

func (s *Server) setupApplication() error {

	// backend, err := etcd.NewBackend(s.etcdClient)
	// if err != nil {
	// 	return err
	// }

	// aclStateRepo := etcd.NewACLStateRepositoryAdapter(backend)
	// aclTokenRepo := etcd.NewACLTokenRepositoryAdapter(backend)
	// aclPolicyRepo := etcd.NewACLPolicyRepositoryAdapter(backend)

	backend := inmem.NewBackend()

	aclStateRepo := inmem.NewACLStateRepositoryAdapter(backend)
	aclTokenRepo := inmem.NewACLTokenRepositoryAdapter(backend)
	aclPolicyRepo := inmem.NewACLPolicyRepositoryAdapter(backend)

	s.services.tokens = application.NewACLTokenService(aclTokenRepo)
	s.services.policies = application.NewACLPolicyService(aclPolicyRepo)
	s.services.acl = application.NewACLService(aclStateRepo, aclTokenRepo, aclPolicyRepo)

	return nil
}
