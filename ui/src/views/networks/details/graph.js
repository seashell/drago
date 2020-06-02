/* eslint-disable no-alert */
/* eslint-disable react/jsx-no-bind */
import React, { useRef, useState, useCallback, useMemo, useEffect } from 'react'
import styled from 'styled-components'
import PropTypes from 'prop-types'
import * as d3 from 'd3-force-3d'

import { ForceGraph2D } from 'react-force-graph'

const Container = styled.div.attrs({
  id: 'wrapper',
})`
  position: absolute;
  background: transparent;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 0;
`

const Graph = ({ nodes, links, onNodeClicked, onNodeHovered, onLinkHovered }) => {
  const fgRef = useRef()

  const NODE_REL_SIZE = 2

  const [highlightedNodes, setHighlightedNodes] = useState(new Set())
  const [highlightedLinks, setHighlightedLinks] = useState(new Set())

  const [hoveredNode, setHoverNode] = useState(null)

  useEffect(() => () => {
    if (fgRef.current) {
      fgRef.current.d3Force('link').distance(l => {
        if (l.source.type === 'host' || l.target.type === 'host') {
          return 6
        }
        return 30
      })
      fgRef.current.d3Force('link').strength(l => {
        if (l.source.type === 'host' || l.target.type === 'host') {
          return 1
        }
        return 0.2
      })
      // See https://github.com/vasturiano/d3-force-3d
      // fgRef.current.d3Force('charge').strength(() => -0.1) // force applied to each node
      // fgRef.current.d3Force('collide', d3.forceCollide(3))
      fgRef.current.d3Force('charge', d3.forceManyBody().strength(-0.1))
    }
  })

  const updateHighlightedElements = () => {
    setHighlightedNodes(highlightedNodes)
    setHighlightedLinks(highlightedLinks)
  }

  const handleNodeHover = node => {
    highlightedNodes.clear()
    highlightedLinks.clear()
    if (node) {
      highlightedNodes.add(node)
      node.neighbors.forEach(neighbor => highlightedNodes.add(neighbor))
      node.links.forEach(link => highlightedLinks.add(link))
    }
    setHoverNode(node || null)
    updateHighlightedElements()
    onNodeHovered(node || null)
  }

  const handleLinkHover = link => {
    highlightedNodes.clear()
    highlightedLinks.clear()
    if (link) {
      highlightedLinks.add(link)
      highlightedNodes.add(link.source)
      highlightedNodes.add(link.target)
    }
    updateHighlightedElements()
    onLinkHovered(
      (link && link.source.type === 'interface' && link.target.type === 'interface'
        ? link
        : null) || null
    )
  }

  const handleNodeClick = node => {
    onNodeClicked(node)
  }

  const nodeCanvasObject = useCallback(
    (node, ctx, globalScale) => {
      if (node.type === 'host') {
        const fontSize = 4 + Math.log10(1.4 * globalScale)
        ctx.fillStyle = '#666666'
        ctx.font = `${fontSize}px Roboto`
        ctx.fillText(`${node.data.name}`, node.x + 8, node.y + 2 - Math.log10(1.7 * globalScale))
        if (node.data.advertiseAddress) {
          ctx.fillStyle = '#cccccc'
          ctx.font = `${fontSize - 1}px Roboto`
          ctx.fillText(
            `${node.data.advertiseAddress}`,
            node.x + 8,
            node.y + 6 - Math.log10(1.7 * globalScale)
          )
        }
      }
      if (node === hoveredNode || highlightedNodes.has(node)) {
        ctx.beginPath()
        ctx.arc(node.x, node.y, NODE_REL_SIZE * 1.8, 0, 2 * Math.PI, false)
        ctx.fillStyle = '#cccccc'
        ctx.fill()
      }
    },
    [hoveredNode]
  )

  const data = useMemo(() => {
    const gData = { nodes, links }
    gData.links.forEach(link => {
      const a = gData.nodes.find(el => el.id === link.source)
      const b = gData.nodes.find(el => el.id === link.target)

      if (a !== undefined) {
        if (!a.links) {
          a.links = []
        }
        if (!a.neighbors) {
          a.neighbors = []
        }
        a.links.push(link)
        a.neighbors.push(b)
      }

      if (b !== undefined) {
        if (!b.links) {
          b.links = []
        }
        if (!b.neighbors) {
          b.neighbors = []
        }
        b.links.push(link)
        b.neighbors.push(a)
      }
    })
    return gData
  }, [nodes, links])

  return (
    <Container>
      <ForceGraph2D
        ref={fgRef}
        graphData={data}
        nodeRelSize={NODE_REL_SIZE}
        nodeVal={n => (n.type === 'host' ? 6 : 1)}
        nodeColor={n => (n.type === 'host' ? '#333333' : '#555555')}
        nodeCanvasObject={nodeCanvasObject}
        nodeCanvasObjectMode={() => 'before'}
        linkWidth={link => {
          if (link.source.type === 'host' || link.target.type === 'host') {
            return 0.1
          }
          return highlightedLinks.has(link) ? 5 : 1
        }}
        linkDirectionalParticles={link =>
          (link.type === link.source.type) === 'host' || link.target.type === 'host' ? 0 : 4
        }
        linkDirectionalParticleWidth={link => (highlightedLinks.has(link) ? 4 : 0)}
        linkDirectionalArrowRelPos={1}
        linkDirectionalArrowLength={link =>
          link.source.type === 'host' || link.target.type === 'host' ? 0 : 3
        }
        linkCurvature={link =>
          link.source.type === 'host' || link.target.type === 'host' ? 0 : 0.25
        }
        onNodeHover={handleNodeHover}
        onLinkHover={handleLinkHover}
        onNodeClick={handleNodeClick}
        cooldownTicks={100}
      />
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
