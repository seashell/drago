import { useMutation, useQuery } from '@apollo/client'
import { useLocation, useParams } from '@reach/router'
import React, { useEffect, useState } from 'react'
import styled from 'styled-components'
import { icons } from '_assets/'
import BackLink from '_components/back-link'
import Box from '_components/box'
import Button from '_components/button'
import { useConfirmationDialog } from '_components/confirmation-dialog'
import EmptyState from '_components/empty-state'
import Icon from '_components/icon'
import List from '_components/list'
import { Dragon as Spinner } from '_components/spinner'
import Text from '_components/text'
import {
  CREATE_CONNECTION,
  CREATE_INTERFACE,
  DELETE_CONNECTION,
  DELETE_INTERFACE,
  UPDATE_CONNECTION,
  UPDATE_INTERFACE,
} from '_graphql/mutations'
import { GET_CONNECTIONS, GET_INTERFACES, GET_NODE } from '_graphql/queries'
import ConnectPeerModal from '_modals/connect-peer'
import JoinNetworkModal from '_modals/join-network'
import ConnectionCard from './connection-card'
import InterfaceCard from './interface-card'

const Container = styled(Box)`
  flex-direction: column;
`

const SectionTitle = styled(Text)`
  color: ${(props) => props.theme.colors.neutralLight};
  font-weight: 600;
  font-size: 0.76rem;
  letter-spacing: 0.06rem;
  text-transform: uppercase;
`

const StyledIcon = styled(Icon)`
  border-radius: 4px;
  background: ${(props) => props.theme.colors.neutralLighter};
  align-items: center;
  justify-content: center;
`

const MetaItemContainer = styled(Box)`
  height: 48px;
  align-items: center;
  :not(:last-child) {
    border-bottom: 1px solid ${(props) => props.theme.colors.neutralLighter};
  }
`

const renderMetaItem = (key, value) => (
  <MetaItemContainer key={key}>
    <Text textStyle="body" width="240px" color="neutralDark">
      {key}
    </Text>
    <Text textStyle="body" color="neutralDark">
      {value}
    </Text>
  </MetaItemContainer>
)

function randomInt(min, max) {
  min = Math.ceil(min)
  max = Math.floor(max)
  return Math.floor(Math.random() * (max - min + 1)) + min
}

const ClientDetails = () => {
  const location = useLocation()
  const { nodeId } = useParams()
  const { confirm } = useConfirmationDialog()

  const [selectedInterfaceId, setSelectedInterfaceId] = useState()
  const [selectedConnectionId, setSelectedConnectionId] = useState()

  const [isJoinNetworkModalOpen, setJoinNetworkModalOpen] = useState(false)
  const [isConnectPeerModalOpen, setConnectPeerModalOpen] = useState(false)

  const getNodeQuery = useQuery(GET_NODE, {
    variables: { id: nodeId },
  })

  const getNodeInterfacesQuery = useQuery(GET_INTERFACES, {
    variables: { nodeId, networkId: '' },
  })

  const getNodeConnectionsQuery = useQuery(GET_CONNECTIONS, {
    variables: { nodeId, interfaceId: '', networkId: '' },
  })

  const [createConnection, createConnectionMutation] = useMutation(CREATE_CONNECTION)
  const [createInterface, createInterfaceMutation] = useMutation(CREATE_INTERFACE)
  const [updateInterface, updateInterfaceMutation] = useMutation(UPDATE_INTERFACE)
  const [updateConnection, updateConnectionMutation] = useMutation(UPDATE_CONNECTION)
  const [deleteInterface, deleteInterfaceMutation] = useMutation(DELETE_INTERFACE)
  const [deleteConnection, deleteConnectionMutation] = useMutation(DELETE_CONNECTION)

  useEffect(() => {
    window.scrollTo(0, 0)
    getNodeQuery.refetch()
    getNodeInterfacesQuery.refetch()
    getNodeConnectionsQuery.refetch()
  }, [location])

  useEffect(() => {
    if (location.state != null) {
      setTimeout(() => {
        setSelectedConnectionId(location.state.connectionId)
      }, 300)
    }
  }, [location])

  const handleJoinNetwork = (id) => {
    createInterface({
      variables: {
        nodeId,
        networkId: id,
      },
    })
      .then(() => {
        getNodeInterfacesQuery.refetch()
      })
      .catch(() => {})
  }

  const handleConnectToPeer = (networkId, peerInterfaceId) => {
    const nodeInterface = getNodeInterfacesQuery.data.result.find(
      (el) => el.NetworkID === networkId
    )

    createConnection({
      variables: {
        connection: {
          peerSettings: {
            [nodeInterface.ID]: {},
            [peerInterfaceId]: {},
          },
        },
      },
    })
      .then(() => {
        getNodeConnectionsQuery.refetch()
      })
      .catch(() => {})
  }

  const handleInterfaceChange = (id, values) => {
    updateInterface({
      variables: {
        id,
        address: values.address,
        listenPort: values.listenPort,
        dns: values.dns,
        mtu: values.mtu,
      },
    })
      .then(() => {
        getNodeInterfacesQuery.refetch()
      })
      .catch(() => {})
  }

  const handleInterfaceDelete = (id) => {
    confirm({
      title: 'Are you sure?',
      details:
        'This will remove the node from the associated network, and destroy all of its connections.',
      isDestructive: true,
      onConfirm: () => {
        deleteInterface({ variables: { id } })
          .then(() => {
            getNodeInterfacesQuery.refetch()
          })
          .catch(() => {})
      },
    })
  }

  const handleConnectionChange = (id, values) => {
    const connection = {
      id,
      PeerSettings: {
        [values.interfaceId]: {
          InterfaceID: values.interfaceId,
          RoutingRules: {
            AllowedIPs: values.allowedIPs,
          },
        },
      },
      persistentKeepalive: values.persistentKeepalive,
    }

    updateConnection({
      variables: {
        id,
        connection,
      },
    })
      .then(() => {
        getNodeConnectionsQuery.refetch()
      })
      .catch(() => {})
  }

  const handleConnectionDelete = (id) => {
    confirm({
      title: 'Are you sure?',
      details: `Depending on your network's topology, this might interfere in the node's ability to communicate with others.`,
      isDestructive: true,
      onConfirm: () => {
        deleteConnection({ variables: { id } })
          .then(() => {
            getNodeConnectionsQuery.refetch()
          })
          .catch(() => {})
      },
    })
  }

  const handleJoinNetworkButtonClick = () => {
    setJoinNetworkModalOpen(true)
  }

  const handleConnectToPeerButtonClick = () => {
    setConnectPeerModalOpen(true)
  }

  const handleInterfaceCardClick = (id) => {
    setSelectedInterfaceId(selectedInterfaceId !== id ? id : undefined)
  }

  const handleConnectionCardClick = (id) => {
    setSelectedConnectionId(selectedConnectionId !== id ? id : undefined)
  }

  const isLoading =
    getNodeQuery.loading ||
    getNodeInterfacesQuery.loading ||
    getNodeConnectionsQuery.loading ||
    createInterfaceMutation.loading ||
    createConnectionMutation.loading ||
    deleteInterfaceMutation.loading ||
    deleteConnectionMutation.loading

  if (isLoading) {
    return <Spinner />
  }

  const node = getNodeQuery.data ? getNodeQuery.data.result : { Interfaces: [], Meta: {} }
  const interfaces = getNodeInterfacesQuery.data ? getNodeInterfacesQuery.data.result : []
  const connections = interfaces.reduce(
    (acc, el) =>
      acc.concat(
        (getNodeConnectionsQuery.data ? getNodeConnectionsQuery.data.result : []).map((c) => ({
          Network: el.Network,
          FromInterfaceID: el.ID,
          FromInterfaceAddress: el.Address,
          FromNodeID: el.NodeID,
          ToInterfaceID: c.Peers.find((p) => p !== el.ID),
          ...c,
        }))
      ),
    []
  )
  
  return (
    <Container pb="72px">
      <JoinNetworkModal
        isOpen={isJoinNetworkModalOpen}
        onJoin={handleJoinNetwork}
        onClose={() => setJoinNetworkModalOpen(false)}
      />
      <ConnectPeerModal
        isOpen={isConnectPeerModalOpen}
        onJoin={handleConnectToPeer}
        onClose={() => setConnectPeerModalOpen(false)}
      />
      <BackLink text="Clients" to={`/ui/clients`} mb={3} />
      <Box alignItems="center" mb={4}>
        <StyledIcon mr="12px" p={2} icon={<icons.Host />} size="48px" color="neutralDarker" />
        <Box width="100%">
          <Box flexDirection="column">
            <Text textStyle="title" mb={1}>
              {node.Name}
            </Text>
            <Text textStyle="detail">{node.ID}</Text>
          </Box>
          {interfaces.length > 0 && (
            <Button variant="primary" ml="auto" onClick={handleJoinNetworkButtonClick}>
              Join Network
            </Button>
          )}
        </Box>
      </Box>

      <Box px={3} py={2} border="discrete" alignItems="center">
        <SectionTitle width="160px">Client details</SectionTitle>
        <Box mr={4}>
          <Text textStyle="strong" fontSize="12px" mr={2}>
            Status
          </Text>
          <Text textStyle="detail" color={node.Status === 'ready' ? 'success' : 'danger'}>
            {node.Status}
          </Text>
        </Box>

        <Box mr={4}>
          <Text textStyle="strong" fontSize="12px" mr={2}>
            Advertise Address
          </Text>
          <Text textStyle="detail" mr={3}>
            {node.AdvertiseAddress ? node.AdvertiseAddress : 'N/A'}
          </Text>
        </Box>

      </Box>

      <Box alignItems="center" mt={4} mb={3}>
        <Text textStyle="subtitle">Interfaces</Text>
      </Box>
      <List>
        {interfaces.map((el) => (
          <InterfaceCard
            key={el.ID}
            id={el.ID}
            showSpinner={updateInterfaceMutation.loading}
            hasPublicKey={el.HasPublicKey}
            publicKey={el.PublicKey}
            name={el.Name}
            connectionsCount={el.ConnectionsCount}
            address={el.Address}
            listenPort={el.ListenPort}
            dns={el.DNS}
            mtu={el.MTU}
            table={el.Table}
            preUp={el.PreUp}
            postUp={el.PostUp}
            preDown={el.PreDown}
            postDown={el.PostDown}
            networkId={el.NetworkID}
            networkName={el.Network.Name}
            networkAddressRange={el.Network.AddressRange}
            onChange={(id, values) => handleInterfaceChange(id, values)}
            onDelete={() => handleInterfaceDelete(el.ID)}
            onClick={() => handleInterfaceCardClick(el.ID)}
            isExpanded={selectedInterfaceId === el.ID}
          />
        ))}
      </List>

      {interfaces.length === 0 && (
        <EmptyState
          title="No networks."
          description="By joining a network clients can communicate with each other."
          image={<Icon my={4} icon={<icons.Network />} size="48px" color="neutral" />}
          extra={
            <Box alignItems="center" mt="24px">
              <Button variant="primary" onClick={handleJoinNetworkButtonClick}>
                Join Network
              </Button>
            </Box>
          }
        />
      )}

      <Box alignItems="flex-end" mt={4} mb={3}>
        <Text textStyle="subtitle">Connections</Text>
      </Box>
      <List>
        {connections.map((el) => (
          <ConnectionCard
            key={el.ID}
            id={el.ID}
            networkName={el.Network.Name}
            fromInterfaceId={el.FromInterfaceID}
            toInterfaceId={el.ToInterfaceID}
            fromInterfaceAddress={el.FromInterfaceAddress}
            fromNodeID={el.FromNodeID}
            createdAt={el.CreatedAt}
            updatedAt={el.UpdatedAt}
            showSpinner={updateConnectionMutation.loading}
            onChange={handleConnectionChange}
            onDelete={() => handleConnectionDelete(el.ID)}
            onClick={() => handleConnectionCardClick(el.ID)}
            isExpanded={selectedConnectionId === el.ID}
          />
        ))}
      </List>

      {interfaces.length > 0 && connections.length === 0 && (
        <EmptyState
          title="No connections."
          description="Connect with peers in your networks to interact wth them."
          image={<Icon my={4} icon={<icons.Connection />} size="48px" color="neutral" />}
          extra={
            <Box alignItems="center" mt="24px">
              <Button variant="primary" onClick={handleConnectToPeerButtonClick}>
                Connect to Peer
              </Button>
            </Box>
          }
        />
      )}
      {interfaces.length === 0 && (
        <EmptyState
          title="No connections."
          description="Nodes need to join a network before they can connect to peers."
          image={<Icon my={4} icon={<icons.Connection />} size="48px" color="neutral" />}
          extra={
            <Box alignItems="center" mt="24px">
              <Button variant="primary" onClick={handleJoinNetworkButtonClick}>
                Join Network
              </Button>
            </Box>
          }
        />
      )}

      <Text textStyle="subtitle" mt={4} mb={3}>
        Meta
      </Text>
      {Object.keys(node.Meta).length === 0 && (
        <EmptyState
          title="No metadata."
          description="Node does not contain any metadata."
          image={<Icon my={4} icon={<icons.Label />} size="48px" color="neutral" />}
        />
      )}
      {Object.keys(node.Meta).map((k) => renderMetaItem(k, node.Meta[k]))}
    </Container>
  )
}

ClientDetails.propTypes = {}

export default ClientDetails
