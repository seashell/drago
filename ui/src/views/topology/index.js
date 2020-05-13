import React, { useState, useMemo, useEffect } from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'
import { navigate, useLocation } from '@reach/router'
import { Portal } from 'react-portal'

import { useQuery } from 'react-apollo'
import { GET_HOSTS, GET_ALL_LINKS } from '_graphql/actions'

import { illustrations } from '_assets/'
import Box from '_components/box'
import Text from '_components/text'
import Button from '_components/button'
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
  border: 'none',
  height: '300px',
})`
  svg {
    height: 300px;
    width: auto;
  }
  padding: 20px;
  flex-direction: column;
  align-items: center;
  justify-content: center;
`

const ErrorState = () => (
  <ErrorStateContainer>
    <illustrations.Error />
    <Text textStyle="description" mt={4}>
      Oops! It seems that an error has occurred.
    </Text>
  </ErrorStateContainer>
)

const EmptyStateContainer = styled(Box).attrs({
  border: 'none',
  height: '300px',
})`
  svg {
    height: 300px;
    width: auto;
  }
  padding: 20px;
  flex-direction: column;
  align-items: center;
  justify-content: center;
`

const EmptyState = () => (
  <EmptyStateContainer>
    <illustrations.Empty />
    <Text textStyle="description" mt={4}>
      Oops! It seems that there are no hosts registered in this network.
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

const Topology = ({ networkId }) => {
  const location = useLocation()

  const [hosts, setHosts] = useState([])
  const [links, setLinks] = useState([])

  const [selectedHostID, setSelectedHostID] = useState(undefined)
  const [selectedLinkID, setSelectedLinkID] = useState(undefined)

  const getHostsQuery = useQuery(GET_HOSTS, {
    variables: { networkId },
    onCompleted: data => {
      if (data === undefined) return
      setHosts(
        data.result.items.map(host => ({
          id: host.id,
          hostObj: host,
          isHover: false,
        }))
      )
    },
    onError: () => {},
  })

  const getLinksQuery = useQuery(GET_ALL_LINKS, {
    variables: { networkId },
    onCompleted: data => {
      if (data === undefined) return
      setLinks(
        data.result.items.map(link => ({
          id: link.id,
          source: link.fromHost,
          target: link.toHost,
          linkObj: link,
          isHover: false,
        }))
      )
    },
    onError: () => {},
  })

  useEffect(() => {
    window.scrollTo(0, 0)
  }, [location])

  const handleNodeClick = n => {
    if (n !== null) {
      navigate(`/networks/${networkId}/hosts/${n.id}`)
    }
  }

  const handleNodeHover = n => {
    setSelectedHostID(n != null ? n.id : undefined)
  }

  const handleLinkHover = l => {
    setSelectedLinkID(l != null ? l.id : undefined)
  }

  const handleListViewButtonClick = () => {
    navigate(`/networks/${networkId}/hosts`)
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
    [handleNodeClick, hosts, links]
  )

  const hoveredNode = hosts.find(h => h.id === selectedHostID) || {
    hostObj: { name: '', address: '' },
  }
  const hoveredLink = links.find(l => l.id === selectedLinkID) || { linkObj: {} }

  return (
    <Container>
      <Box mb={3} width="100%">
        <Text textStyle="title">Overlay topology</Text>
        <Button
          onClick={handleListViewButtonClick}
          variant="primaryInverted"
          borderRadius={3}
          width="100px"
          height="40px"
          ml="auto"
          style={{ zIndex: '99' }}
        >
          List view
        </Button>
      </Box>

      {isLoading && <Spinner />}
      {isError ? <ErrorState /> : isEmpty ? <EmptyState /> : MemoizedGraph}

      {selectedLinkID && (
        <Portal>
          <StyledLinkCard
            sourceName={hoveredLink.source.hostObj.name}
            sourceAddress={hoveredLink.source.hostObj.ipAddress}
            targetName={hoveredLink.target.hostObj.name}
            targetAddress={hoveredLink.target.hostObj.ipAddress}
            allowedIPs={hoveredLink.linkObj.allowedIps}
            persistentKeepalive={hoveredLink.linkObj.persistentKeepalive}
          />
        </Portal>
      )}
      {selectedHostID && (
        <Portal>
          <StyledHostCard
            name={hoveredNode.hostObj.name}
            address={hoveredNode.hostObj.ipAddress}
            advertiseAddress={hoveredNode.hostObj.advertiseAddress}
            listenPort={hoveredNode.hostObj.listenPort}
          />
        </Portal>
      )}
    </Container>
  )
}

Topology.propTypes = {
  networkId: PropTypes.string.isRequired,
}

export default Topology
