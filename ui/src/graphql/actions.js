import gql from 'graphql-tag'

const GET_NODES = gql`
  query getNodes {
    result @rest(type: "NodesPayload", path: "nodes") {
      count
      items @type(name: "Node") {
        id
        name
        publicKey
        advertiseAddr
        interface @type(name: "Interface") {
          address
          listenPort
        }
      }
    }
  }
`

const GET_NODE = gql`
  query getNode($id: String) {
    result: getNode(id: $id) @rest(path: "nodes/{args.id}", type: "Node") {
      id
      name
      publicKey
      advertiseAddr
      interface @type(name: "Interface") {
        address
        listenPort
      }
    }
  }
`

const CREATE_NODE = gql`
  mutation createNode($name: String!, $address: String!) {
    createNode(input: { name: $name, address: $address })
      @rest(method: "POST", path: "nodes", type: "Node") {
      id
      name
      interface @type(name: "Interface") {
        address
      }
    }
  }
`

const UPDATE_NODE = gql`
  mutation updateNode($id: Int!, $name: String!, $address: String!) {
    updateNode(id: $id, input: { name: $name, address: $address })
      @rest(method: "PUT", path: "nodes/{args.id}", type: "Node") {
      id
      name
      interface @type(name: "Interface") {
        address
      }
    }
  }
`

const DELETE_NODE = gql`
  mutation deleteNode($id: Int!) {
    deleteNode(id: $id) @rest(method: "DELETE", path: "nodes/{args.id}", type: "Node") {
      id
    }
  }
`

export { GET_NODES, GET_NODE, CREATE_NODE, UPDATE_NODE, DELETE_NODE }
