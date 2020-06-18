import gql from 'graphql-tag'

export const GET_LINKS = gql`
  query getLinks($networkId: String!, $fromHostId: String!, $fromInterfaceId: String!) {
    result: links(
      networkId: $networkId
      fromHostId: $fromHostId
      fromInterfaceId: $fromInterfaceId
    ) @rest(method: "GET", type: "PaginatedResult", path: "links?{args}") {
      page
      perPage
      pageCount
      totalCount
      items @type(name: "Link") {
        id
        fromInterfaceId @export(as: "fromInterfaceId")
        toInterfaceId @export(as: "toInterfaceId")
        allowedIps
        persistentKeepalive
        fromInterface
          @rest(
            method: "GET"
            type: "Interface"
            path: "interfaces/{exportVariables.fromInterfaceId}"
          ) {
          id
          name
          ipAddress
          hostId @export(as: "hostId")
          host @rest(method: "GET", type: "Host", path: "hosts/{exportVariables.hostId}") {
            name
          }
        }
        toInterface
          @rest(
            method: "GET"
            type: "Interface"
            path: "interfaces/{exportVariables.toInterfaceId}"
          ) {
          id
          name
          ipAddress
          hostId @export(as: "hostId")
          host @rest(method: "GET", type: "Host", path: "hosts/{exportVariables.hostId}") {
            name
          }
        }
      }
    }
  }
`

export const GET_LINK = gql`
  query getLink($id: String) {
    result: getLink(id: $id) @rest(path: "links/{args.id}", type: "Link") {
      id
      fromInterfaceId @export(as: "fromInterfaceId")
      toInterfaceId @export(as: "toInterfaceId")
      allowedIps
      persistentKeepalive
      fromInterface
        @rest(
          method: "GET"
          type: "Interface"
          path: "interfaces/{exportVariables.fromInterfaceId}"
        ) {
        id
        name
        ipAddress
        hostId @export(as: "hostId")
        host @rest(method: "GET", type: "Host", path: "hosts/{exportVariables.hostId}") {
          name
        }
      }
      toInterface
        @rest(
          method: "GET"
          type: "Interface"
          path: "interfaces/{exportVariables.toInterfaceId}"
        ) {
        id
        name
        ipAddress
        hostId @export(as: "hostId")
        host @rest(method: "GET", type: "Host", path: "hosts/{exportVariables.hostId}") {
          name
        }
      }
    }
  }
`

export const UPDATE_LINK = gql`
  mutation updateLink(
    $id: String!
    $fromInterfaceId: String!
    $toInterfaceId: String!
    $allowedIps: [String]!
    $persistenKeepalive: Int!
  ) {
    updateLink(
      id: $id
      input: {
        fromInterfaceId: $fromInterfaceId
        toInterfaceId: $toInterfaceId
        allowedIps: $allowedIps
        persistentKeepalive: $persistentKeepalive
      }
    ) @rest(method: "PATCH", path: "links/{args.id}", type: "Link") {
      id
      fromInterfaceId
      toInterfaceId
      allowedIps
      persistentKeepalive
    }
  }
`

export const CREATE_LINK = gql`
  mutation createLink(
    $fromInterfaceId: String!
    $toInterfaceId: String!
    $allowedIps: [String]!
    $persistenKeepalive: Int!
  ) {
    createLink(
      input: {
        fromInterfaceId: $fromInterfaceId
        toInterfaceId: $toInterfaceId
        allowedIps: $allowedIps
        persistentKeepalive: $persistentKeepalive
      }
    ) @rest(method: "POST", path: "links", type: "Link") {
      id
      fromInterfaceId
      toInterfaceId
      allowedIps
      persistentKeepalive
    }
  }
`

export const DELETE_LINK = gql`
  mutation deleteLink($id: Int!) {
    deleteLink(id: $id) @rest(method: "DELETE", path: "links/{args.id}", type: "Link") {
      id
    }
  }
`
