import gql from 'graphql-tag'

const GET_NODES = gql`
  query getNodes {
    result @rest(type: "NodesPayload", path: "nodes") {
      count
      items @type(name: "Node") {
        id
        label
        address
      }
    }
  }
`

const GET_NODE = gql`
  query getNode($nodeId: String) {
    result: getNode(nodeId: $nodeId) @rest(type: "Node", path: "nodes/{args.nodeId}") {
      id
      label
      address
    }
  }
`

const CREATE_NODE = gql`
  mutation test {
    data
  }
`

export { GET_NODES, GET_NODE, CREATE_NODE }
