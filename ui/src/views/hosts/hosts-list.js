import React from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'
import moment from 'moment'

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
    <icons.EmptyStateCube />
    <Text textStyle="description" mt={4}>
      Oops! It seems that there are no hosts registered.
    </Text>
  </EmptyStateContainer>
)

const HostsList = ({ hosts, onHostSelect, onHostDelete }) =>
  hosts.length === 0 ? (
    <EmptyState />
  ) : (
    hosts.map(h => (
      <HostCard
        key={h.id}
        id={h.id}
        label={h.name}
        address={h.address}
        lastSeen={h.lastSeen}
        onClick={onHostSelect}
        onDelete={e => onHostDelete(e, h.id)}
      />
    ))
  )

HostsList.propTypes = {
  hosts: PropTypes.arrayOf(PropTypes.object).isRequired,
  onHostSelect: PropTypes.func.isRequired,
  onHostDelete: PropTypes.func.isRequired,
}

export default HostsList
