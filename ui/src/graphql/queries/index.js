import gql from 'graphql-tag'

export const GET_SELF_TOKEN = gql`
  query getSelf {
    result: user @rest(path: "/api/acl/tokens/self", type: "Token") {
      ID
      Type
      Name
      Secret
      Policies
      CreatedAt
      UpdatedAt
    }
  }
`

export const GET_NETWORKS = gql`
  query getNetworks {
    result: networks @rest(type: "Network", path: "/api/networks/") {
      ID
      Name
      AddressRange
      InterfacesCount
      ConnectionsCount
      CreatedAt
      UpdatedAt
    }
  }
`

export const GET_NETWORK = gql`
  query getNetwork($id: Int!) {
    result: getNetwork(id: $id) @rest(type: "Network", path: "/api/networks/{args.id}") {
      ID
      Name
      AddressRange
      InterfacesCount
      ConnectionsCount
      CreatedAt
      UpdatedAt
    }
  }
`

export const GET_NETWORK_WITH_INTERFACES = gql`
  query getNetwork($id: Int!) {
    result: getNode(id: $id) @rest(type: "Network", path: "/api/networks/{args.id}") {
      ID @export(as: "networkId")
      Name
      AddressRange
      CreatedAt
      UpdatedAt
      Interfaces
        @rest(type: "Interface", path: "/api/interfaces/?network={exportVariables.networkId}") {
        ID
        Name
        HasPublicKey
        NodeID @export(as: "nodeId")
        Node @rest(type: "Node", path: "/api/nodes/{exportVariables.nodeId}/") {
          ID
          Name
          Status
          AdvertiseAddress
        }
        NetworkID
        ConnectionsCount
        Address
      }
      CreatedAt
      UpdatedAt
    }
  }
`

export const GET_NODES = gql`
  query getNodes {
    result: nodes @rest(type: "Node", path: "/api/nodes/") {
      ID
      Name
      Status
      InterfacesCount
      ConnectionsCount
      CreatedAt
      UpdatedAt
    }
  }
`

export const GET_NODE = gql`
  query getNode($id: Int!) {
    result: getNode(id: $id) @rest(type: "Node", path: "/api/nodes/{args.id}") {
      ID
      Name
      AdvertiseAddress
      Meta
      Status
      CreatedAt
      UpdatedAt
    }
  }
`

export const GET_NODE_WITH_INTERFACES = gql`
  query getNode($id: Int!) {
    result: getNode(id: $id) @rest(type: "Node", path: "/api/nodes/{args.id}") {
      ID @export(as: "nodeId")
      Name
      Address
      Status
      Meta
      Interfaces @rest(type: "Interface", path: "/api/interfaces/?node={exportVariables.nodeId}") {
        ID
        Name
        Address
        NodeID
        ConnectionsCount
        NetworkID @export(as: "networkId")
        Network @rest(type: "Network", path: "/api/networks/{exportVariables.networkId}") {
          ID
          Name
          AddressRange
        }
      }
      CreatedAt
      UpdatedAt
    }
  }
`

// TODO: migrate interface-node aggregation to the
// server to avoid multiple requests
export const GET_PEERS = gql`
  query getPeers($networkId: String!) {
    result: getPeers(networkId: $networkId)
      @rest(type: "Interface", path: "/api/interfaces/?network={args.networkId}") {
      ID
      Name
      Address
      NetworkID
      NodeID @export(as: "nodeId")
      Node @rest(type: "Node", path: "/api/nodes/{exportVariables.nodeId}") {
        ID
        Name
        Address
      }
      CreatedAt
      UpdatedAt
    }
  }
`

export const GET_PEER = gql`
  query getPeer($interfaceId: String!) {
    result: getPeer(interfaceId: $interfaceId)
      @rest(type: "Interface", path: "/api/interfaces/{args.interfaceId}") {
      ID
      Name
      Address
      NetworkID
      NodeID @export(as: "nodeId")
      Node @rest(type: "Node", path: "/api/nodes/{exportVariables.nodeId}") {
        ID
        Name
        Address
      }
      CreatedAt
      UpdatedAt
    }
  }
`

export const GET_INTERFACES = gql`
  query getInterfaces($nodeId: String!, $networkId: String!) {
    result: getInterfaces(nodeId: $nodeId, networkId: $networkId)
      @rest(
        type: "Interface"
        path: "/api/interfaces/?node={args.nodeId}&network={args.networkId}"
      ) {
      ID
      PublicKey
      HasPublicKey
      Name
      Address
      ListenPort
      DNS
      NodeID
      NetworkID @export(as: "networkId")
      ConnectionsCount
      CreatedAt
      UpdatedAt
      Network @rest(type: "Network", path: "/api/networks/{exportVariables.networkId}") {
        ID
        Name
        AddressRange
      }
    }
  }
`

export const GET_INTERFACE = gql`
  query getInterface($interfaceId: String!) {
    result: getInterface(interfaceId: $interfaceId)
      @rest(type: "Interface", path: "/api/interfaces/{args.interfaceId}") {
      ID
      Name
      Address
      PublicKey
      ListenPort
      NodeID
      NetworkID
      ConnectionsCount
      HasPublicKey
      CreatedAt
      UpdatedAt
    }
  }
`

export const GET_CONNECTIONS = gql`
  query getConnections($interfaceId: String!, $nodeId: String!, $networkId: String!) {
    result: getConnections(interfaceId: $interfaceId, nodeId: $nodeId, networkId: $networkId)
      @rest(
        type: "Connection"
        path: "/api/connections/?interface={args.interfaceId}&node={args.nodeId}&network={args.networkId}"
      ) {
      ID
      Peers
      PersistentKeepalive
      NetworkID @export(as: "networkId")
      CreatedAt
      UpdatedAt
      Network @rest(type: "Network", path: "/api/networks/{exportVariables.networkId}") {
        ID
        Name
        AddressRange
      }
    }
  }
`

export const GET_CONNECTION = gql`
  query getConnection($connectionId: String!) {
    result: getConnection(connectionId: $connectionId)
      @rest(type: "Connection", path: "/api/connections/{args.connectionId}") {
      ID
      PeerSettings
      PersistentKeepalive
      NetworkID @export(as: "networkId")
      CreatedAt
      UpdatedAt
    }
  }
`
