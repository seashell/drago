import React from 'react'
import PropTypes from 'prop-types'

import HostCard from './host-card'

const HostsList = ({ hosts, onHostSelect, onHostDelete }) =>
  hosts.map(n => (
    <HostCard
      key={n.id}
      id={n.id}
      name={n.name}
      address={n.ipAddress}
      advertiseAddress={n.advertiseAddress}
      onClick={onHostSelect}
      onDelete={e => onHostDelete(e, n.id)}
    />
  ))

HostsList.propTypes = {
  hosts: PropTypes.arrayOf(PropTypes.object).isRequired,
  onHostSelect: PropTypes.func.isRequired,
  onHostDelete: PropTypes.func.isRequired,
}

export default HostsList
