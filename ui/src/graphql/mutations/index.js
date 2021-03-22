import gql from 'graphql-tag'

export const CREATE_NETWORK = gql`
  mutation createNetwork($name: String!, $addressRange: String!) {
    createNetwork(input: { Name: $name, AddressRange: $addressRange })
      @rest(method: "POST", path: "/api/networks/", type: "Network") {
      ID
    }
  }
`

export const DELETE_NETWORK = gql`
  mutation deleteNetwork($id: Int!) {
    deleteNetwork(id: $id)
      @rest(method: "DELETE", path: "/api/networks/{args.id}", type: "Network") {
      ID
    }
  }
`

export const CREATE_INTERFACE = gql`
  mutation createInterface($name: String!, $addressRange: String!) {
    createInterface(input: { NetworkID: $networkId, NodeID: $nodeId })
      @rest(method: "POST", path: "/api/interfaces/", type: "Interface") {
      ID
    }
  }
`

export const UPDATE_INTERFACE = gql`
  mutation updateInterface(
    $id: String!
    $address: String!
    $listenPort: Number!
    $dns: String!
    $mtu: String!
  ) {
    updateInterface(
      id: $id
      input: { id: $id, address: $address, listenPort: $listenPort, dns: $dns, mtu: $mtu }
    ) @rest(method: "POST", path: "/api/interfaces/{args.id}", type: "Interface") {
      ID
    }
  }
`

export const DELETE_INTERFACE = gql`
  mutation deleteInterface($id: Int!) {
    deleteInterface(id: $id)
      @rest(method: "DELETE", path: "/api/interfaces/{args.id}", type: "Interface") {
      ID
    }
  }
`

export const CREATE_CONNECTION = gql`
  mutation createConnection($connection: Connection!) {
    createConnection(input: $connection)
      @rest(method: "POST", path: "/api/connections/", type: "Connection") {
      ID
    }
  }
`

export const UPDATE_CONNECTION = gql`
  mutation updateConnection($id: String!, $connection: Connection!) {
    updateConnection(id: $id, input: $connection)
      @rest(method: "POST", path: "/api/connections/{args.id}", type: "Connection") {
      ID
    }
  }
`

export const DELETE_CONNECTION = gql`
  mutation deleteConnection($id: Int!) {
    deleteConnection(id: $id)
      @rest(method: "DELETE", path: "/api/connections/{args.id}", type: "Connection") {
      ID
    }
  }
`
