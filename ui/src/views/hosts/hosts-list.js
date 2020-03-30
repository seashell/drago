import React from 'react'
import PropTypes from 'prop-types'

import styled from 'styled-components'

import Box from '_components/box'
import Text from '_components/text'
import { icons } from '_assets/'

import HostCard from './host-card'

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
    <icons.GhostCube />
    <Text textStyle="description" mt={4}>
      Oops! It seems that there are no hosts registered.
    </Text>
  </EmptyStateContainer>
)

const HostsList = ({ hosts, onHostSelect, onHostDelete }) =>
  hosts.length === 0 ? (
    <EmptyState />
  ) : (
    hosts.map(n => (
      <HostCard
        key={n.id}
        id={n.id}
        label={n.name}
        address={n.address}
        onClick={onHostSelect}
        onDelete={e => onHostDelete(e, n.id)}
      />
    ))
  )

HostsList.propTypes = {
  hosts: PropTypes.arrayOf(PropTypes.object).isRequired,
  onHostSelect: PropTypes.func.isRequired,
  onHostDelete: PropTypes.func.isRequired,
}

export default HostsList
