import React from 'react'
import PropTypes from 'prop-types'

import styled from 'styled-components'

import Box from '_components/box'
import Text from '_components/text'

import NodeCard from './node-card'

const EmptyStateContainer = styled(Box).attrs({
  border: 'discrete',
  height: '100px',
})``

const EmptyState = () => (
  <EmptyStateContainer>
    <Text>[TODO: Empty state] - No nodes found</Text>
  </EmptyStateContainer>
)

const NodesList = ({ nodes, onNodeSelect, onNodeDelete }) =>
  nodes.length === 0 ? (
    <EmptyState />
  ) : (
    nodes.map(n => (
      <NodeCard
        key={n.id}
        id={n.id}
        label={n.name}
        address={n.interface.address}
        onClick={onNodeSelect}
        onDelete={e => onNodeDelete(e, n.id)}
      />
    ))
  )

NodesList.propTypes = {
  nodes: PropTypes.arrayOf(PropTypes.object).isRequired,
  onNodeSelect: PropTypes.func.isRequired,
  onNodeDelete: PropTypes.func.isRequired,
}

export default NodesList
