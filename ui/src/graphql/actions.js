import gql from 'graphql-tag'

const GET_NETWORKS = gql`
  query getNetworks {
    result @rest(type: "NetworksPayload", path: "networks") {
      page
      perPage
      pageCount
      totalCount
      items @type(name: "Network") {
        id
        name
        ipAddressRange
        createdAt
        updatedAt
      }
    }
  }
`

const CREATE_NETWORK = gql`
  mutation createNetwork($name: String!, $ipAddressRange: String!) {
    createNetwork(input: { name: $name, ipAddressRange: $ipAddressRange })
      @rest(method: "POST", path: "networks", type: "Network") {
      id
    }
  }
`

const DELETE_NETWORK = gql`
  mutation deleteNetwork($id: Int!) {
    deleteNetwork(id: $id) @rest(method: "DELETE", path: "networks/{args.id}", type: "Network") {
      id
    }
  }
`

const GET_HOSTS = gql`
  query getHosts($networkId: String) {
    result: getHosts(networkId: $networkId)
      @rest(type: "HostsPayload", path: "hosts?networkId={args.networkId}") {
      page
      perPage
      pageCount
      totalCount
      items @type(name: "Host") {
        id
        name
        publicKey
        advertiseAddress
        ipAddress
        listenPort
        lastSeen
      }
    }
  }
`

const GET_HOST = gql`
  query getHost($networkId: String, $id: String) {
    result: getHost(networkId: $networkId, id: $id)
      @rest(path: "hosts/{args.id}?networkId={args.networkId}", type: "Host") {
      id
      name
      publicKey
      advertiseAddress
      ipAddress
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
    $networkId: String
    $name: String!
    $ipAddress: String!
    $advertiseAddress: String
  ) {
    createHost(
      input: {
        networkId: $networkId
        name: $name
        ipAddress: $ipAddress
        advertiseAddress: $advertiseAddress
      }
    ) @rest(method: "POST", path: "hosts?networkId={args.networkId}", type: "Host") {
      id
      name
      ipAddress
      advertiseAddress
      listenPort
    }
  }
`

const UPDATE_HOST = gql`
  mutation updateHost(
    $id: Int!
    $name: String!
    $ipAddress: String!
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
    $networkId: String
  ) {
    updateHost(
      id: $id
      networkId: $networkId
      input: {
        name: $name
        ipAddress: $ipAddress
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
    ) @rest(method: "PATCH", path: "hosts/{args.id}?networkId={args.networkId}", type: "Host") {
      id
      name
      publicKey
      advertiseAddress
      ipAddress
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
  mutation deleteHost($networkId: String!, $id: String!) {
    deleteHost(networkId: $networkId, id: $id)
      @rest(method: "DELETE", path: "hosts/{args.id}?networkId={args.networkId}", type: "Host") {
      id
    }
  }
`

const CREATE_LINK = gql`
  mutation createLink($networkId: String, $from: String!, $to: String!) {
    createLink(
      networkId: $networkId
      input: {
        fromHostId: $from
        toHostId: $to
        allowedIPs: $allowedIPs
        persistentKeepalive: $persistentKeepalive
      }
    ) @rest(method: "POST", path: "links?networkId={args.networkId}", type: "Link") {
      id
      from
      to
      allowedIPs
      persistentKeepalive
    }
  }
`

const DELETE_LINK = gql`
  mutation deleteLink($networkId: String!, $id: String!) {
    deleteLink(networkId: $networkId, id: $id)
      @rest(method: "DELETE", path: "links/{args.id}?networkId={args.networkId}", type: "Link") {
      id
    }
  }
`

const GET_LINKS_FROM_HOST = gql`
  query getLinksFromHost($networkId: String!, $hostId: String!) {
    result: getLinksFromHost(networkId: $networkId, hostId: $hostId)
      @rest(type: "LinksPayload", path: "links?networkId={args.networkId}&hostId={args.hostId}") {
      page
      perPage
      pageCount
      totalCount
      items @type(name: "Link") {
        id
        toHost
        fromHost
        allowedIps
        persistentKeepalive
      }
    }
  }
`

const GET_ALL_LINKS = gql`
  query getAllLinks($networkId: String!, $hostId: String!) {
    result: getAllLinks(networkId: $networkId, hostId: $hostId)
      @rest(type: "LinksPayload", path: "links?networkId={args.networkId}") {
      page
      perPage
      pageCount
      totalCount
      items @type(name: "Link") {
        id
        toHost
        fromHost
        allowedIps
        persistentKeepalive
      }
    }
  }
`

export {
  GET_NETWORKS,
  CREATE_NETWORK,
  DELETE_NETWORK,
  GET_HOSTS,
  GET_HOST,
  CREATE_HOST,
  UPDATE_HOST,
  DELETE_HOST,
  GET_ALL_LINKS,
  GET_LINKS_FROM_HOST,
  CREATE_LINK,
  DELETE_LINK,
}
