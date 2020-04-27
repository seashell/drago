import React, { useEffect, useRef } from 'react'
import styled from 'styled-components'
import PropTypes from 'prop-types'

import { ForceGraph2D } from 'react-force-graph'
import Box from '_components/box'
import Text from '_components/text'

import { icons } from '_assets/'

const Container = styled.div.attrs({
  id: 'wrapper',
})`
  width: 100%;
`

const EmptyStateContainer = styled(Box).attrs({
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

const EmptyState = () => (
  <EmptyStateContainer>
    <icons.EmptyStateCube />
    <Text textStyle="description" mt={4}>
      Oops! It seems that there are no registered hosts.
    </Text>
  </EmptyStateContainer>
)

const Graph = ({ nodes, links, onNodeClicked, onNodeHovered, onLinkHovered }) => {
  const ref = useRef(null)

  useEffect(() => {
    const fg = ref.current
    if (fg !== null) {
      fg.d3Force('link').distance(60)
    }
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
    if (node) {
      node.isHover = true
      elem.style.cursor = 'pointer'
    } else {
      prevNode.isHover = false
    }
    if (prevNode) {
      prevNode.isHover = false
    }
    onNodeHovered(node)
  }

  const handleLinkHover = (link, prevLink) => {
    const elem = document.getElementById('wrapper')
    elem.style.cursor = null
    if (link) {
      link.isHover = true
      elem.style.cursor = 'pointer'
    } else {
      prevLink.isHover = false
    }
    if (prevLink) {
      prevLink.isHover = false
    }
    onLinkHovered(link)
  }

  const handleNodeClick = n => {
    onNodeClicked(n)
  }

  return (
    <Container>
      {nodes.length === 0 ? (
        <EmptyState />
      ) : (
        <ForceGraph2D
          ref={ref}
          graphData={{ nodes, links }}
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
      )}
    </Container>
  )
}

Graph.propTypes = {
  nodes: PropTypes.arrayOf(PropTypes.object).isRequired,
  links: PropTypes.arrayOf(PropTypes.object).isRequired,
  onNodeClicked: PropTypes.func.isRequired,
  onNodeHovered: PropTypes.func.isRequired,
  onLinkHovered: PropTypes.func.isRequired,
}

export default Graph
