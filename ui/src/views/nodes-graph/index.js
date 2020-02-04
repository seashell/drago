import React, { useEffect, useRef } from 'react'
import ReactDOM from 'react-dom'
import { navigate } from '@reach/router'

import styled from 'styled-components'
import { ForceGraph2D } from 'react-force-graph'
import { GET_NODES } from '_graphql/actions'
import { useQuery } from 'react-apollo'

const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  grid-column: span 12;
`

const NodesGraph = () => {
  const { loading, error, data } = useQuery(GET_NODES)

  const ref = useRef(null)

  const nodes = loading
    ? []
    : data.result.items.map(node => ({ id: parseInt(node.id, 10), label: node.label }))

  const links = []
  if (!loading) {
    nodes.forEach(node => {
      if (node.id !== '1') links.push({ source: 1, target: parseInt(node.id, 10) })
    })
    console.log(links)
  }
  const graphData = {
    nodes,
    links,
  }

  useEffect(() => {
    ReactDOM.render(
      <ForceGraph2D
        graphData={graphData}
        nodeColor={n => n.color || '#333333'}
        nodeLabel={n => n.label}
        nodeRelSize={4}
        onNodeClick={n => navigate(`/nodes/${n.id}`)}
      />,
      ref.current
    )
  })

  return (
    <Container>
      <div ref={ref} />
    </Container>
  )
}

export default NodesGraph
