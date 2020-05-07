import React from 'react'
import PropTypes from 'prop-types'
import NetworkCard from './network-card'

const NetworksList = ({ networks, onNetworkSelect, onNetworkDelete }) =>
  networks.map(n => (
    <NetworkCard
      key={n.id}
      id={n.id}
      name={n.name}
      hostCount={n.hostCount}
      onClick={onNetworkSelect}
      onDelete={e => onNetworkDelete(e, n.id)}
    />
  ))

NetworksList.propTypes = {
  networks: PropTypes.arrayOf(PropTypes.object).isRequired,
  onNetworkSelect: PropTypes.func.isRequired,
  onNetworkDelete: PropTypes.func.isRequired,
}

export default NetworksList
