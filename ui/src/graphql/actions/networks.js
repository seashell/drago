import gql from 'graphql-tag'

export const GET_NETWORKS = gql`
  query getNetworks {
    result: networks @rest(type: "PaginatedResult", path: "networks") {
      page
      perPage
      pageCount
      totalCount
      items @type(name: "Network") {
        id @export(as: "networkId")
        name
        ipAddressRange
        createdAt
        updatedAt
        hosts
          @rest(
            type: "PaginatedResult"
            path: "hosts?networkId={exportVariables.networkId}&page=1&perPage=1"
          ) {
          count: totalCount
        }
      }
    }
  }
`

export const CREATE_NETWORK = gql`
  mutation createNetwork($name: String!, $ipAddressRange: String!) {
    createNetwork(input: { name: $name, ipAddressRange: $ipAddressRange })
      @rest(method: "POST", path: "networks", type: "Network") {
      id
    }
  }
`

export const DELETE_NETWORK = gql`
  mutation deleteNetwork($id: Int!) {
    deleteNetwork(id: $id) @rest(method: "DELETE", path: "networks/{args.id}", type: "Network") {
      id
    }
  }
`
