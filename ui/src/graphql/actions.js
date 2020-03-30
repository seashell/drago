import gql from 'graphql-tag'

const GET_HOSTS = gql`
  query getHosts {
    result @rest(type: "HostsPayload", path: "hosts") {
      count
      items @type(name: "Host") {
        id
        name
        publicKey
        advertiseAddress
        address
        listenPort
        lastSeen
      }
    }
  }
`

const GET_HOST = gql`
  query getHost($id: String) {
    result: getHost(id: $id) @rest(path: "hosts/{args.id}", type: "Host") {
      id
      name
      publicKey
      advertiseAddress
      address
      listenPort
      lastSeen
      table
      dns
      preUp
      postUp
      preDown
      postDown
      publicKey
      jwt
      links @type(name: "LinksPayload") {
        count
        items @type(name: "Link") {
          id
          to
          allowedIPs
          persistentKeepalive
        }
      }
    }
  }
`

const CREATE_HOST = gql`
  mutation createHost(
    $name: String!
    $address: String!
    $advertiseAddress: String
    $listenPort: String
  ) {
    createHost(
      input: {
        name: $name
        address: $address
        advertiseAddress: $advertiseAddress
        listenPort: $listenPort
      }
    ) @rest(method: "POST", path: "hosts", type: "Host") {
      id
      name
      address
      advertiseAddress
      listenPort
    }
  }
`

const UPDATE_HOST = gql`
  mutation updateHost(
    $id: Int!
    $name: String!
    $address: String!
    $advertiseAddress: String
    $listenPort: String
    $publicKey: String
    $table: String
    $mtu: String
    $dns: String
    $preUp: String
    $postUp: String
    $preDown: String
    $postDown: String
  ) {
    updateHost(
      id: $id
      input: {
        name: $name
        address: $address
        advertiseAddress: $advertiseAddress
        listenPort: $listenPort
        publicKey: $publicKey
        table: $table
        mtu: $mtu
        dns: $dns
        preUp: $preUp
        postUp: $postUp
        preDown: $preDown
        postDown: $postDown
      }
    ) @rest(method: "PUT", path: "hosts/{args.id}", type: "Host") {
      id
      name
      publicKey
      advertiseAddress
      address
      listenPort
      table
      mtu
      dns
      preUp
      postUp
      preDown
      postDown
      publicKey
    }
  }
`

const DELETE_HOST = gql`
  mutation deleteHost($id: Int!) {
    deleteHost(id: $id) @rest(method: "DELETE", path: "hosts/{args.id}", type: "Host") {
      id
    }
  }
`

const CREATE_LINK = gql`
  mutation createLink($from: String!, $to: String!) {
    createLink(
      input: {
        from: $from
        to: $to
        allowedIPs: $allowedIPs
        persistentKeepalive: $persistentKeepalive
      }
    ) @rest(method: "POST", path: "links", type: "Link") {
      id
      from
      to
      allowedIPs
      persistentKeepalive
    }
  }
`

const DELETE_LINK = gql`
  mutation deleteLink($id: Int!) {
    deleteLink(id: $id) @rest(method: "DELETE", path: "links/{args.id}", type: "Link") {
      id
    }
  }
`

const GET_LINKS = gql`
  query getHosts {
    result @rest(type: "LinksPayload", path: "links") {
      count
      items @type(name: "Link") {
        id
        to @type(name: "Host") {
          id
          name
          address
          advertiseAddress
        }
        from @type(name: "Host") {
          id
          name
          address
          advertiseAddress
        }
        allowedIPs
        persistentKeepalive
      }
    }
  }
`

export {
  GET_HOSTS,
  GET_HOST,
  CREATE_HOST,
  UPDATE_HOST,
  DELETE_HOST,
  GET_LINKS,
  CREATE_LINK,
  DELETE_LINK,
}
