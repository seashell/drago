import gql from 'graphql-tag'

export const GET_INTERFACES = gql`
  query getInterfaces($networkId: String, $hostId: String) {
    result: interfaces(networkId: $networkId, hostId: $hostId)
      @rest(type: "PaginatedResult", path: "interfaces?{args}") {
      page
      perPage
      pageCount
      totalCount
      items @type(name: "Interface") {
        id @export(as: "interfaceId")
        name
        hostId @export(as: "hostId")
        networkId
        ipAddress
        listenPort
        table
        dns
        preUp
        postUp
        preDown
        postDown
        publicKey
        links
          @rest(
            type: "PaginatedResult"
            path: "links?fromInterfaceId={exportVariables.interfaceId}&page=1&perPage=1"
          ) {
          count: totalCount
        }
        parentHost @rest(type: "Host", path: "hosts/{exportVariables.hostId}") {
          id
          name
          advertiseAddress
        }
      }
    }
  }
`

export const GET_INTERFACE = gql`
  query getInterface($id: String) {
    result: getInterface(id: $id) @rest(path: "interfaces/{args.id}", type: "Interface") {
      id
      name
      hostId
      networkId
      ipAddress
      listenPort
      table
      dns
      preUp
      postUp
      preDown
      postDown
      publicKey
    }
  }
`

export const CREATE_INTERFACE = gql`
  mutation createInterface(
    $hostId: String!
    $networkId: String!
    $name: String!
    $ipAddress: String!
    $listenPort: Int!
  ) {
    createInterface(
      input: {
        hostId: $hostId
        networkId: $networkId
        name: $name
        ipAddress: $ipAddress
        listenPort: $listenPort
      }
    ) @rest(method: "POST", path: "interfaces", type: "Interface") {
      id
    }
  }
`

export const UPDATE_INTERFACE = gql`
  mutation updateInterface(
    $id: String!
    $name: String!
    $networkId: String!
    $ipAddress: String!
    $listenPort: Int!
  ) {
    updateInterface(
      input: {
        id: $id
        name: $name
        networkId: $networkId
        ipAddress: $ipAddress
        listenPort: $listenPort
      }
    ) @rest(method: "PATCH", path: "interfaces/{args.id}", type: "Interface") {
      id
      name
      networkId
      ipAddress
      listenPort
    }
  }
`

export const DELETE_INTERFACE = gql`
  mutation deleteInterface($id: Int!) {
    deleteInterface(id: $id)
      @rest(method: "DELETE", path: "interfaces/{args.id}", type: "Interface") {
      id
    }
  }
`
