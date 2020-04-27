import React, { useState } from 'react'
import { navigate } from '@reach/router'

import styled from 'styled-components'
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
`

const StyledHostCard = styled(HostCard)`
  position: absolute;
  top: 100px;
  left: 20px;
`

const StyledLinkCard = styled(LinkCard)`
  position: absolute;
  top: 100px;
  left: 20px;
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

const Topology = () => {
  const [hoveredNodeID, setHoveredNodeID] = useState(undefined)
  const [hoveredLinkID, setHoveredLinkID] = useState(undefined)

  const [hosts, setHosts] = useState([])
  const [links, setLinks] = useState([])

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

  const hoveredNode = hosts.find(h => h.id === hoveredNodeID)
  const hoveredLink = links.find(l => l.id === hoveredLinkID)

  const handleNodeClick = n => {
    navigate(`/hosts/${n.id}`)
  }

  const handleNodeHover = n => {
    setHoveredNodeID(n.id)
  }

  const handleLinkHover = l => {
    setHoveredLinkID(l.id)
  }

  return (
    <Container>
      <Box mb={3}>
        <Text textStyle="title">Overlay topology</Text>
      </Box>
      {(getHostsQuery.loading || getLinksQuery.loading) && <Spinner />}
      {hoveredNodeID !== undefined && (
        <StyledHostCard
          name={hoveredNode.hostObj.name}
          address={hoveredNode.hostObj.address}
          advertiseAddress={hoveredNode.hostObj.advertiseAddress}
          listenPort={hoveredNode.hostObj.listenPort}
          lastSeen={hoveredNode.hostObj.lastSeen}
        />
      )}
      {hoveredLinkID !== undefined && (
        <StyledLinkCard
          sourceName={hoveredLink.sourceObj.name}
          sourceAddress={hoveredLink.sourceObj.address}
          targetName={hoveredLink.targetObj.name}
          targetAddress={hoveredLink.targetObj.address}
          allowedIPs={hoveredLink.allowedIPs}
          persistentKeepalive={hoveredLink.persistentKeepalive}
        />
      )}
      {getHostsQuery.error || getLinksQuery.error ? (
        <ErrorState />
      ) : (
        <Graph
          nodes={hosts}
          links={links}
          onNodeHovered={handleNodeHover}
          onLinkHovered={handleLinkHover}
          onNodeClicked={handleNodeClick}
        />
      )}
    </Container>
  )
}

export default Topology
