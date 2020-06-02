import React, { useState, useEffect } from 'react'
import PropTypes from 'prop-types'
import { useLocation, navigate } from '@reach/router'

import { useQuery, useMutation } from 'react-apollo'
import { GET_INTERFACES, DELETE_INTERFACE } from '_graphql/actions/interfaces'

import Box from '_components/box'
import toast from '_components/toast'
import Button from '_components/button'
import SearchInput from '_components/inputs/search-input'
import { Dragon as Spinner } from '_components/spinner'
import EmptyState from '_components/empty-state'

import InterfaceCard from './interface-card'

const HostInterfacesTab = ({ hostId }) => {
  const [searchFilter, setSearchFilter] = useState('')
  const location = useLocation()

  const getHostInterfacesQuery = useQuery(GET_INTERFACES, {
    variables: { hostId },
  })

  const [deleteInterface, deleteInterfaceMutation] = useMutation(DELETE_INTERFACE, {
    variables: { id: null },
    onCompleted: () => {
      toast.success('Interface deleted')
      getHostInterfacesQuery.refetch()
    },
    onError: () => {
      toast.error('Error deleting interface')
    },
  })

  useEffect(() => {
    getHostInterfacesQuery.refetch()
  }, [location])

  const handleInterfaceSelect = ({ id }) => {
    navigate(`/interfaces/${id}`, { state: { hostId } })
  }

  const handleInterfaceDelete = ({ id }) => {
    deleteInterface({ variables: { id } })
  }

  const handleCreateButtonClicked = () => {
    navigate(`/interfaces/new`, { state: { hostId } })
  }

  const handleSearchInputChanged = e => {
    setSearchFilter(e.target.value)
  }

  const isLoading = getHostInterfacesQuery.loading || deleteInterfaceMutation.loading

  const filteredInterfaces = isLoading
    ? []
    : getHostInterfacesQuery.data.result.items.filter(
        el =>
          el.name.includes(searchFilter) ||
          (el.ipAddress != null && el.ipAddress.includes(searchFilter))
      )

  return isLoading ? (
    <Spinner />
  ) : (
    <Box flexDirection="column">
      <Box my={3}>
        <SearchInput
          width="100%"
          placeholder="Search..."
          onChange={handleSearchInputChanged}
          mr={2}
        />
        <Button onClick={handleCreateButtonClicked}>Create</Button>
      </Box>
      {filteredInterfaces.length === 0 ? (
        <EmptyState />
      ) : (
        filteredInterfaces.map(iface => (
          <InterfaceCard
            key={iface.id}
            id={iface.id}
            name={iface.name}
            ipAddress={iface.ipAddress}
            listenPort={iface.listenPort}
            numLinks={iface.links.count}
            onClick={() => handleInterfaceSelect(iface)}
            onDelete={() => handleInterfaceDelete(iface)}
          />
        ))
      )}
    </Box>
  )
}

HostInterfacesTab.propTypes = {
  hostId: PropTypes.string.isRequired,
}

export default HostInterfacesTab
