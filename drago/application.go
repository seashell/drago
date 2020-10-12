package drago

import (
	"context"

	application "github.com/seashell/drago/drago/application"
	structs "github.com/seashell/drago/drago/application/structs"
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

	s.setupDefaultACLPolicies()

	return nil
}

func (s *Server) setupDefaultACLPolicies() error {
	defaultPolicies := []*domain.ACLPolicy{
		{Name: "anonymous"},
		{Name: "manager"},
	}
	for _, p := range defaultPolicies {
		s.services.ACLPolicies.Upsert(context.Background(), &structs.ACLPolicyUpsertInput{
			Name:        p.Name,
			Description: p.Description,
			Rules:       "",
		})
	}
	return nil
}
