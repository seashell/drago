import gql from 'graphql-tag'

export const GET_HOST_WITH_INTERFACES_AND_LINKS = gql`
  query hosts($id: String) {
    result: getHost(id: $id) @rest(path: "hosts/{args.id}", type: "Host") {
      id @export(as: "hostId")
      name
      advertiseAddress
      interfaces
        @rest(type: "PaginatedResult", path: "interfaces?hostId={exportVariables.hostId}") {
        page
        perPage
        pageCount
        totalCount
        items @type(name: "Interface") {
          id
          name
          advertiseAddress
        }
      }
      links @rest(type: "PaginatedResult", path: "links?hostId={exportVariables.hostId}") {
        page
        perPage
        pageCount
        totalCount
        items @type(name: "Link") {
          id
          fromInterface
          toInterface
          allowedIps
          persistentKeepalive
        }
      }
    }
  }
`

export const GET_HOSTS = gql`
  query getHosts($networkId: String) {
    result: hosts(networkId: $networkId) @rest(type: "PaginatedResult", path: "hosts?{args}") {
      page
      perPage
      pageCount
      totalCount
      items @type(name: "Host") {
        id
        name
        advertiseAddress
      }
    }
  }
`

export const GET_HOST = gql`
  query getHost($id: String) {
    result: getHost(id: $id) @rest(path: "hosts/{args.id}", type: "Host") {
      id
      name
      advertiseAddress
    }
  }
`

export const CREATE_HOST = gql`
  mutation createHost($name: String!, $advertiseAddress: String) {
    createHost(input: { name: $name, advertiseAddress: $advertiseAddress })
      @rest(method: "POST", path: "hosts", type: "Host") {
      id
      name
      advertiseAddress
    }
  }
`

export const UPDATE_HOST = gql`
  mutation updateHost($id: Int!, $name: String!, $advertiseAddress: String) {
    updateHost(id: $id, input: { name: $name, advertiseAddress: $advertiseAddress })
      @rest(method: "PATCH", path: "hosts/{args.id}", type: "Host") {
      id
      name
      advertiseAddress
    }
  }
`

export const DELETE_HOST = gql`
  mutation deleteHost($id: String!) {
    deleteHost(id: $id) @rest(method: "DELETE", path: "hosts/{args.id}", type: "Host") {
      id
    }
  }
`
