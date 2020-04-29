import React, { useState, useMemo } from 'react'
import styled from 'styled-components'
import { navigate } from '@reach/router'
import { Portal } from 'react-portal'

import { useQuery } from 'react-apollo'
import { GET_HOSTS, GET_LINKS } from '_graphql/actions'

import { icons } from '_assets/'
import Box from '_components/box'
import Text from '_components/text'
import { Dragon as Spinner } from '_components/spinner'

import HostCard from './host-card'
import LinkCard from './link-card'

import Graph from './graph'

const Container = styled.div`
  display: flex;
  flex-direction: column;
  grid-column: span 12;
  height: 100vh;
`

const ErrorStateContainer = styled(Box).attrs({
  border: 'discrete',
  height: '300px',
})`
  svg {
    height: 120px;
  }
  padding: 20px;
  flex-direction: column;
  align-items: center;
  justify-content: center;
`

const ErrorState = () => (
  <ErrorStateContainer>
    <icons.ErrorStateCube />
    <Text textStyle="description" mt={4}>
      Oops! It seems that an error has occurred.
    </Text>
  </ErrorStateContainer>
)

const EmptyStateContainer = styled(Box).attrs({
  border: 'discrete',
  height: '300px',
})`
  svg {
    height: 120px;
  }
  flex-direction: column;
  align-items: center;
  justify-content: center;
`

const EmptyState = () => (
  <EmptyStateContainer>
    <icons.EmptyStateCube />
    <Text textStyle="description" mt={4}>
      Oops! It seems that there are no registered hosts.
    </Text>
  </EmptyStateContainer>
)

const StyledHostCard = styled(HostCard)`
  position: absolute;
  top: 100px;
  left: 20px;
  background: #fff;
  z-index: 99;
`

const StyledLinkCard = styled(LinkCard)`
  position: absolute;
  top: 100px;
  left: 20px;
  background: #fff;
  z-index: 99;
`

const Topology = () => {
  const [hosts, setHosts] = useState([])
  const [links, setLinks] = useState([])

  const [selectedHostID, setSelectedHostID] = useState(undefined)
  const [selectedLinkID, setSelectedLinkID] = useState(undefined)

  const getHostsQuery = useQuery(GET_HOSTS, {
    onCompleted: data => {
      if (data === undefined) return
      setHosts(
        data.result.items.map(host => ({
          id: parseInt(host.id, 10),
          name: host.name,
          hostObj: host,
          isHover: false,
        }))
      )
    },
    onError: () => {},
  })

  const getLinksQuery = useQuery(GET_LINKS, {
    onCompleted: data => {
      if (data === undefined) return
      setLinks(
        data.result.items.map(link => ({
          id: parseInt(link.id, 10),
          source: link.from.id,
          target: link.to.id,
          sourceObj: link.from,
          targetObj: link.to,
          allowedIPs: link.allowedIPs,
          persistentKeepalive: link.persistentKeepalive,
          isHover: false,
        }))
      )
    },
    onError: () => {},
  })

  const handleNodeClick = n => {
    if (n !== null) {
      navigate(`/hosts/${n.id}`)
    }
  }

  const handleNodeHover = n => {
    setSelectedHostID(n != null ? n.id : undefined)
  }

  const handleLinkHover = l => {
    setSelectedLinkID(l != null ? l.id : undefined)
  }

  const isError = getHostsQuery.error || getLinksQuery.error
  const isLoading = getHostsQuery.loading || getLinksQuery.loading
  const isEmpty = !isLoading && hosts.length === 0

  const MemoizedGraph = useMemo(
    () => (
      <Graph
        nodes={hosts}
        links={links}
        onNodeHovered={handleNodeHover}
        onLinkHovered={handleLinkHover}
        onNodeClicked={handleNodeClick}
      />
    ),
    [hosts, links]
  )

  const hoveredNode = hosts.find(h => h.id === selectedHostID) || { hostObj: { name: '' } }
  const hoveredLink = links.find(l => l.id === selectedLinkID) || { sourceObj: {}, targetObj: {} }

  return (
    <Container>
      <Box mb={3} width="100%">
        <Text textStyle="title">Overlay topology</Text>
      </Box>

      {isLoading && <Spinner />}
      {isError ? <ErrorState /> : isEmpty ? <EmptyState /> : MemoizedGraph}

      {selectedLinkID && (
        <Portal>
          <StyledLinkCard
            sourceName={hoveredLink.sourceObj.name}
            sourceAddress={hoveredLink.sourceObj.address}
            targetName={hoveredLink.targetObj.name}
            targetAddress={hoveredLink.targetObj.address}
            allowedIPs={hoveredLink.allowedIPs}
            persistentKeepalive={hoveredLink.persistentKeepalive}
          />
        </Portal>
      )}
      {selectedHostID && (
        <Portal>
          <StyledHostCard
            name={hoveredNode.hostObj.name}
            address={hoveredNode.hostObj.address}
            advertiseAddress={hoveredNode.hostObj.advertiseAddress}
            listenPort={hoveredNode.hostObj.listenPort}
            lastSeen={hoveredNode.hostObj.lastSeen}
          />
        </Portal>
      )}
    </Container>
  )
}

export default Topology
