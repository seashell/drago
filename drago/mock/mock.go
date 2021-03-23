package mock

import (
	"context"

	"github.com/seashell/drago/drago/state"
	"github.com/seashell/drago/drago/structs"
	"github.com/seashell/drago/pkg/util"
	"github.com/seashell/drago/pkg/uuid"
)

// AuthHandler :
type AuthHandler struct{}

// Authorize : authorize anything
func (h *AuthHandler) Authorize() error {
	return nil
}

// PopulateRepository :
func PopulateRepository(repo state.Repository) error {

	ctx := context.TODO()

	nodeID1 := "8cbc8089-e294-3fab-9f79-84ea6700c431"
	nodeID2 := "8dd4c160-3034-47f0-94b5-7de423d36424"

	repo.UpsertACLToken(ctx, &structs.ACLToken{ID: uuid.Generate(), Type: structs.ACLTokenTypeManagement, Name: "Bootstrap Token", Secret: "abc"})
	repo.UpsertACLToken(ctx, &structs.ACLToken{ID: uuid.Generate(), Type: structs.ACLTokenTypeClient, Name: "client-12343234", Secret: "xyz", Policies: []string{"read", "write"}})

	// Create nodes
	node1 := &structs.Node{ID: nodeID1, Name: "test-node", SecretID: "84dd48eb-5f4d-f8aa-6bb6-bed687d9ed56", Meta: map[string]string{"location": "uk"}, Status: structs.NodeStatusInit, AdvertiseAddress: util.StrToPtr("192.168.100.7")}
	node2 := &structs.Node{ID: nodeID2, Name: "other-test-node", SecretID: "c9a1b3dc-bbc9-49dc-83ed-77ff8abb1f30", Meta: map[string]string{"location": "us"}, Status: structs.NodeStatusInit}

	// Create networks
	net1 := &structs.Network{ID: "8579e9cc-787b-4e57-b37f-088ed4f491f2", Name: "network-1", AddressRange: "192.168.0.0/16"}

	// Populate repository with interfaces
	ifaceID1 := "c01648a1-b675-455a-8e5b-29db18be6663"
	ifaceID2 := "618969bc-60b8-4018-8bf4-d2f4fdce43ae"

	iface1 := &structs.Interface{
		ID:        ifaceID1,
		NodeID:    nodeID1,
		NetworkID: net1.ID,
		Name:      util.StrToPtr("wg0"),
		PublicKey: util.StrToPtr("uNAObp9zCLkivCIv/mKvgNUVtgVRoDegtLnaGtVeQWo="),
		Address:   util.StrToPtr("192.168.0.200/24")}
	iface2 := &structs.Interface{
		ID:        ifaceID2,
		NodeID:    nodeID2,
		NetworkID: net1.ID,
		Name:      util.StrToPtr("wg0"),
		Address:   util.StrToPtr("192.168.0.2/24")}

	conn1 := &structs.Connection{
		ID:        "14b62335-ba2b-4a05-8c6d-29b4e11f86b6",
		NetworkID: net1.ID,
		NodeIDs:   []string{iface1.NodeID, iface2.NodeID},
		PeerSettings: map[string]*structs.PeerSettings{
			ifaceID1: {
				InterfaceID: ifaceID1,
				RoutingRules: &structs.RoutingRules{
					AllowedIPs: []string{"dhjadbasj"},
				},
			},
			ifaceID2: {
				InterfaceID: ifaceID2,
				RoutingRules: &structs.RoutingRules{
					AllowedIPs: []string{},
				},
			},
		},
	}

	iface1.UpsertConnection(conn1.ID)
	iface2.UpsertConnection(conn1.ID)

	node1.UpsertInterface(ifaceID1)
	node2.UpsertInterface(ifaceID2)

	node1.UpsertConnection(conn1.ID)
	node2.UpsertConnection(conn1.ID)

	net1.UpsertInterface(ifaceID1)
	net1.UpsertInterface(ifaceID2)
	net1.UpsertConnection(conn1.ID)

	// Commit to the repository
	repo.UpsertNetwork(ctx, net1)

	repo.UpsertNode(ctx, node1)
	repo.UpsertNode(ctx, node2)

	repo.UpsertInterface(ctx, iface1)
	repo.UpsertInterface(ctx, iface2)

	repo.UpsertConnection(ctx, conn1)

	return nil
}
