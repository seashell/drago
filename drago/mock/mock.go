package mock

import (
	"context"

	"github.com/seashell/drago/drago/state"
	"github.com/seashell/drago/drago/structs"
	"github.com/seashell/drago/pkg/uuid"
)

func PopulateWithData(repo state.Repository) error {

	ctx := context.TODO()

	nid := "8cbc8089-e294-3fab-9f79-84ea6700c431"
	sid := "84dd48eb-5f4d-f8aa-6bb6-bed687d9ed56"

	repo.UpsertNode(ctx, &structs.Node{ID: nid, Name: "eduardo-gram", SecretID: sid, Status: structs.NodeStatusInit})
	repo.UpsertNode(ctx, &structs.Node{ID: uuid.Generate(), Name: "node-1", SecretID: uuid.Generate(), Status: structs.NodeStatusDown})
	repo.UpsertNode(ctx, &structs.Node{ID: uuid.Generate(), Name: "node-2", SecretID: uuid.Generate(), Status: structs.NodeStatusDown})
	repo.UpsertNode(ctx, &structs.Node{ID: uuid.Generate(), Name: "node-3", SecretID: uuid.Generate(), Status: structs.NodeStatusDown})

	netID := "8579e9cc-787b-4e57-b37f-088ed4f491f2"

	repo.UpsertNetwork(ctx, &structs.Network{ID: netID, Name: "network-1", AddressRange: "192.168.0.0/16"})
	repo.UpsertNetwork(ctx, &structs.Network{ID: uuid.Generate(), Name: "network-2", AddressRange: "10.1.1.0/24"})
	repo.UpsertNetwork(ctx, &structs.Network{ID: uuid.Generate(), Name: "network-3", AddressRange: "192.0.0.0/8"})

	repo.UpsertInterface(ctx, &structs.Interface{ID: "c01648a1-b675-455a-8e5b-29db18be6663", NodeID: nid, NetworkID: netID, Name: "wg0", Address: "192.168.0.1/24", ListenPort: 9999})
	repo.UpsertInterface(ctx, &structs.Interface{ID: "69731e10-4b70-4261-abfe-84e21d008ab7", NodeID: nid, NetworkID: netID, Name: "wg1", Address: "192.168.0.2/24", ListenPort: 9898})

	// repo.UpsertLink(ctx, &structs.Link{ID: uuid.Generate(), NodeID: nid, NetworkID: netID, Name: "wg1", Address: "192.168.0.2"})

	return nil
}
