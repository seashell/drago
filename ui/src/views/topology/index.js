import React, { useEffect, useRef } from 'react'
import ReactDOM from 'react-dom'
import { navigate } from '@reach/router'

import styled from 'styled-components'
import { ForceGraph2D } from 'react-force-graph'
import { GET_HOSTS, GET_LINKS } from '_graphql/actions'
import { useQuery } from 'react-apollo'
import { Dragon } from '_components/spinner'

const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  grid-column: span 12;
`

const Topology = () => {
  const getHostsQuery = useQuery(GET_HOSTS)
  const getLinksQuery = useQuery(GET_LINKS)

  const ref = useRef(null)

  const hosts = getHostsQuery.loading
    ? []
    : getHostsQuery.data.result.items.map(host => ({
        id: parseInt(host.id, 10),
        label: host.name,
        hostObj: host,
      }))

  const links = getLinksQuery.loading
    ? []
    : getLinksQuery.data.result.items.map(link => ({
        id: parseInt(link.id, 10),
        source: link.from.id,
        target: link.to.id,
        sourceObj: link.from,
        targetObj: link.to,
      }))

  const graphData = {
    nodes: hosts,
    links,
  }

  useEffect(() => {
    const fg = ref.current
    // fg.d3Force('link').distance(80)
    // fg.d3Force('charge').strength(-800)
  })

  return (
    <Container>
      {(getHostsQuery.loading || getLinksQuery.loading) && <Dragon />}
      <ForceGraph2D
        ref={ref}
        graphData={graphData}
        nodeColor={n => n.color || '#333333'}
        nodeLabel={n => n.label}
        nodeRelSize={4}
        linkDirectionalArrowLength={3.5}
        linkDirectionalArrowRelPos={1}
        linkCurvature={0.25}
        nodeCanvasObjectMode={() => 'after'}
        nodeCanvasObject={(node, ctx, globalScale) => {
          const fontSize = 1 + 1.4 * globalScale
          ctx.fillStyle = '#ff0000'
          ctx.font = `${fontSize}px Roboto`
          if (node.hostObj.advertiseAddress !== null) {
            ctx.fillText(`${node.hostObj.advertiseAddress}`, node.x, node.y - 2 * globalScale)
          }
        }}
        linkCanvasObjectMode={() => 'after'}
        linkCanvasObject={(link, ctx, globalScale) => {
          const fontSize = 1.5 + 0.8 * globalScale
          ctx.fillStyle = '#aaaaaa'
          ctx.font = `${fontSize}px Roboto`

          ctx.fillText(
            `From ${link.sourceObj.name} (${link.sourceObj.address})`,
            link.source.x + 5,
            link.source.y
          )
          ctx.fillText(
            `To ${link.targetObj.name} (${link.targetObj.address})`,
            link.source.x + 5,
            link.source.y + 1.8 * globalScale
          )
        }}
        onNodeClick={n => navigate(`/hosts/${n.id}`)}
      />
    </Container>
  )
}

export default Topology
