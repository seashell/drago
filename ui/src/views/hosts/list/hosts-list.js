import React from 'react'
import PropTypes from 'prop-types'

import HostCard from './host-card'

const HostsList = ({ hosts, onHostSelect, onHostDelete }) =>
  hosts.map(h => (
    <HostCard
      key={h.id}
      id={h.id}
      name={h.name}
      labels={h.labels}
      advertiseAddress={h.advertiseAddress}
      onClick={onHostSelect}
      onDelete={e => onHostDelete(e, h.id)}
    />
  ))

HostsList.propTypes = {
  hosts: PropTypes.arrayOf(PropTypes.object).isRequired,
  onHostSelect: PropTypes.func.isRequired,
  onHostDelete: PropTypes.func.isRequired,
}

export default HostsList
