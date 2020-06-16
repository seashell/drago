import React, { useState, useEffect, useMemo } from 'react'
import styled from 'styled-components'
import { navigate, useLocation, useParams } from '@reach/router'
import { Portal } from 'react-portal'

import { useQuery } from 'react-apollo'
import { GET_LINKS } from '_graphql/actions/links'
import { GET_INTERFACES } from '_graphql/actions/interfaces'

import Box from '_components/box'
import Text from '_components/text'
import ErrorState from '_components/error-state'
import EmptyState from '_components/empty-state'
import { Dragon as Spinner } from '_components/spinner'

import InterfaceCard from './interface-card'
import LinkCard from './link-card'

import Graph from './graph'

const Container = styled.div`
  display: flex;
  flex-direction: column;
  grid-column: span 12;
`

const StyledInterfaceCard = styled(InterfaceCard)`
  position: absolute;
  bottom: 32px;
  right: 8px;
  background: #fff;
  z-index: 99;
`

const StyledLinkCard = styled(LinkCard)`
  position: absolute;
  bottom: 32px;
  right: 8px;
  background: #fff;
  z-index: 99;
`

const Topology = () => {
  const location = useLocation()
  const urlParams = useParams()

  const [interfaces, setInterfaces] = useState([])
  const [links, setLinks] = useState([])

  const [selectedInterfaceID, setSelectedInterfaceID] = useState(undefined)
  const [selectedLinkID, setSelectedLinkID] = useState(undefined)

  const getInterfacesQuery = useQuery(GET_INTERFACES, {
    variables: { networkId: urlParams.networkId },
    onCompleted: data => {
      if (data === undefined) return
      setInterfaces(data.result.items)
    },
  })

  const getLinksQuery = useQuery(GET_LINKS, {
    variables: {},
    onCompleted: data => {
      if (data === undefined) return
      setLinks(data.result.items)
    },
  })

  useEffect(() => {
    window.scrollTo(0, 0)
    getLinksQuery.refetch()
    getInterfacesQuery.refetch()
  }, [location])

  const handleNodeClick = iface => {
    if (iface !== null) {
      navigate(`/interfaces/${iface.id}`)
    }
  }

  const handleNodeHover = n => {
    if (n !== null) {
      if (n.type === 'interface') {
        setSelectedInterfaceID(n.id)
        return
      }
    }
    setSelectedInterfaceID(undefined)
  }

  const handleLinkHover = l => {
    if (l !== null) {
      if (l.source.type === 'interface' && l.target.type === 'interface') {
        setSelectedLinkID(l.id)
        return
      }
    }
    setSelectedLinkID(undefined)
  }

  const gNodes = interfaces
    .map(iface => ({
      type: 'interface',
      id: iface.id,
      data: iface,
    }))
    .concat(
      interfaces
        .map(iface => ({
          type: 'host',
          id: iface.parentHost.id,
          data: iface.parentHost,
        }))
        .filter((host, index, self) => index === self.findIndex(h => h.id === host.id))
    )

  const gLinks = links
    .map(link => ({
      id: link.id,
      source: link.fromInterfaceId,
      target: link.toInterfaceId,
      label: link.allowedIps,
      data: link,
    }))
    .concat(
      interfaces.map(iface => ({
        id: `${iface.id}+${iface.parentHost.id}`,
        target: iface.id,
        source: iface.parentHost.id,
      }))
    )

  const MemoizedGraph = useMemo(
    () => (
      <Graph
        links={gLinks}
        nodes={gNodes}
        onNodeHovered={handleNodeHover}
        onLinkHovered={handleLinkHover}
        onNodeClicked={handleNodeClick}
      />
    ),
    [links, interfaces]
  )

  const hoveredInterface = interfaces.find(h => h.id === selectedInterfaceID) || {
    data: { name: '', ipAddress: '', listenPort: '' },
  }
  const hoveredLink = links.find(l => l.id === selectedLinkID) || { data: {} }

  const isError = getInterfacesQuery.error || getLinksQuery.error
  const isLoading = getInterfacesQuery.loading || getLinksQuery.loading
  const isEmpty = !isLoading && interfaces.length === 0

  if (isLoading) {
    return <Spinner />
  }

  return isError ? (
    <ErrorState />
  ) : isEmpty ? (
    <EmptyState description="Oops! It seems that you don't have any hosts yet registered in this network." />
  ) : (
    <Container>
      <Box mb={3} width="100%">
        <Text textStyle="title">Network topology</Text>
      </Box>

      {isLoading && <Spinner />}
      {isError ? <ErrorState /> : isEmpty ? <EmptyState /> : MemoizedGraph}

      {selectedLinkID && (
        <Portal>
          <StyledLinkCard
            sourceName={hoveredLink.fromInterface.name}
            sourceAddress={hoveredLink.fromInterface.ipAddress}
            targetName={hoveredLink.toInterface.name}
            targetAddress={hoveredLink.toInterface.ipAddress}
            allowedIPs={hoveredLink.allowedIps}
            persistentKeepalive={hoveredLink.persistentKeepalive}
          />
        </Portal>
      )}
      {selectedInterfaceID && (
        <Portal>
          <StyledInterfaceCard
            name={hoveredInterface.name}
            address={hoveredInterface.ipAddress}
            listenPort={hoveredInterface.listenPort}
          />
        </Portal>
      )}
    </Container>
  )
}

Topology.propTypes = {}

export default Topology
