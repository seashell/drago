package drago

import (
	"context"

	application "github.com/seashell/drago/drago/application"
	domain "github.com/seashell/drago/drago/domain"
	inmem "github.com/seashell/drago/drago/infrastructure/storage/inmem"
)

func (s *Server) setupApplication() error {

	// etcdClient, err := clientv3.New(clientv3.Config{
	// 	Endpoints:        s.config.Etcd.InitialAdvertiseClientURLs,
	// 	AutoSyncInterval: time.Second * 5,
	// 	DialTimeout:      5 * time.Second,
	// })
	// if err != nil {
	// 	return err
	// }

	// backend, err := etcd.NewBackend(s.etcdClient)
	// if err != nil {
	// 	return err
	// }
	// aclStateRepo := etcd.NewACLStateRepositoryAdapter(backend)
	// aclTokenRepo := etcd.NewACLTokenRepositoryAdapter(backend)
	// aclPolicyRepo := etcd.NewACLPolicyRepositoryAdapter(backend)

	backend := inmem.NewBackend()
	aclStateRepository := inmem.NewACLStateRepositoryAdapter(backend)
	aclTokenRepository := inmem.NewACLTokenRepositoryAdapter(backend)
	aclPolicyRepository := inmem.NewACLPolicyRepositoryAdapter(backend)
	networkRepository := inmem.NewNetworkRepositoryAdapter(backend)
	hostRepository := inmem.NewHostRepositoryAdapter(backend)

	// Setup default policies
	for _, p := range s.defaultACLPolicies() {
		_, err := aclPolicyRepository.Upsert(context.TODO(), p)
		if err != nil {
			return err
		}
	}

	s.services = application.New(&application.Config{
		ACLEnabled:          s.config.ACL.Enabled,
		ACLStateRepository:  aclStateRepository,
		ACLTokenRepository:  aclTokenRepository,
		ACLPolicyRepository: aclPolicyRepository,
		NetworkRepository:   networkRepository,
		HostRepository:      hostRepository,
		AuthHandler: NewAuthorizationHandlerAdapter(
			aclTokenRepository,
			aclPolicyRepository,
		),
	})

	return nil
}

func (s *Server) defaultACLPolicies() []*domain.ACLPolicy {
	return []*domain.ACLPolicy{
		{
			Name: "anonymous",
			Rules: []*domain.ACLPolicyRule{
				{"policy", "", []string{application.ACLPolicyList}},
			},
		},
		{
			Name:  "manager",
			Rules: []*domain.ACLPolicyRule{},
		},
	}
}
