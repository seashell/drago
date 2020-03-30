import React, { useEffect, useRef, useState } from 'react'
import { navigate } from '@reach/router'

import styled from 'styled-components'
import { ForceGraph2D } from 'react-force-graph'
import { GET_HOSTS, GET_LINKS } from '_graphql/actions'
import { useQuery } from 'react-apollo'
import { Dragon } from '_components/spinner'
import HostCard from './host-card'
import LinkCard from './link-card'

const Container = styled.div`
  display: flex;
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  max-height: 100vh;
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

const Topology = () => {
  const ref = useRef(null)

  const [hosts, setHosts] = useState([])
  const [links, setLinks] = useState([])

  const [hoveredNodeID, setHoveredNodeID] = useState(undefined)
  const [hoveredLinkID, setHoveredLinkID] = useState(undefined)

  const getHostsQuery = useQuery(GET_HOSTS, {
    onCompleted: data => {
      setHosts(
        data.result.items.map(host => ({
          id: parseInt(host.id, 10),
          name: host.name,
          hostObj: host,
          isHover: false,
        }))
      )
    },
  })

  const getLinksQuery = useQuery(GET_LINKS, {
    onCompleted: data => {
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
  })

  const graphData = {
    nodes: hosts,
    links,
  }

  useEffect(() => {
    const fg = ref.current
    fg.d3Force('link').distance(60)
    // fg.d3Force('charge').strength(-800)
  })

  const nodeCanvasObject = (node, ctx, globalScale) => {
    const fontSize = 1 + 1.4 * globalScale

    ctx.fillStyle = '#cccccc'
    ctx.font = `${fontSize}px Roboto`

    if (node.hostObj.advertiseAddress !== null) {
      ctx.fillText(`${node.hostObj.advertiseAddress}`, node.x + 2, node.y - 1.7 * globalScale)
    }

    if (node.isHover) {
      ctx.beginPath()
      ctx.arc(node.x, node.y, 6 * 1.4, 0, 2 * Math.PI, false)
      ctx.fillStyle = '#cccccc'
      ctx.fill()
    }
  }

  const linkCanvasObject = (link, ctx, globalScale) => {}

  const handleNodeHover = (node, prevNode) => {
    const elem = document.getElementById('wrapper')
    elem.style.cursor = null

    setHoveredNodeID(undefined)

    if (node) {
      node.isHover = true
      setHoveredNodeID(node.id)
      elem.style.cursor = 'pointer'
    } else {
      prevNode.isHover = false
    }

    if (prevNode) {
      prevNode.isHover = false
    }
  }

  const handleLinkHover = (link, prevLink) => {
    const elem = document.getElementById('wrapper')
    elem.style.cursor = null

    setHoveredLinkID(undefined)

    if (link) {
      link.isHover = true
      setHoveredLinkID(link.id)
      elem.style.cursor = 'pointer'
    } else {
      prevLink.isHover = false
    }
    if (prevLink) {
      prevLink.isHover = false
    }
  }

  const handleNodeClick = n => {
    navigate(`/hosts/${n.id}`)
  }

  const hoveredNode = hosts.find(h => h.id === hoveredNodeID)
  const hoveredLink = links.find(l => l.id === hoveredLinkID)

  return (
    <Container>
      {(getHostsQuery.loading || getLinksQuery.loading) && <Dragon />}
      <div id="wrapper">
        {hoveredNodeID !== undefined && (
          <StyledHostCard
            name={hoveredNode.hostObj.name}
            address={hoveredNode.hostObj.address}
            advertiseAddress={hoveredNode.hostObj.advertiseAddress}
            listenPort={hoveredNode.hostObj.listenPort}
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
        <ForceGraph2D
          ref={ref}
          graphData={graphData}
          nodeRelSize={5}
          nodeLabel={null}
          nodeColor={n => n.color || '#333333'}
          nodeCanvasObjectMode={() => 'before'}
          nodeCanvasObject={nodeCanvasObject}
          onNodeHover={handleNodeHover}
          onNodeClick={handleNodeClick}
          linkCanvasObjectMode={() => 'after'}
          linkCanvasObject={linkCanvasObject}
          linkWidth={l => (l.isHover ? 5 : 1)}
          onLinkHover={handleLinkHover}
          linkDirectionalArrowLength={3.5}
          linkDirectionalArrowRelPos={1}
          linkCurvature={0.1}
          linkDirectionalParticles={1}
          // linkDirectionalParticleSpeed={0.05}
          linkDirectionalParticleColor={() => '#dddddd'}
        />
      </div>
    </Container>
  )
}

export default Topology
