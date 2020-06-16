import React, { useEffect, useState } from 'react'
import styled from 'styled-components'
import { navigate, useLocation } from '@reach/router'
import { useQuery, useMutation } from 'react-apollo'

import { GET_NETWORKS, DELETE_NETWORK } from '_graphql/actions/networks'

import { Dragon as Spinner } from '_components/spinner'
import SearchInput from '_components/inputs/search-input'
import ErrorState from '_components/error-state'
import EmptyState from '_components/empty-state'
import Button from '_components/button'
import toast from '_components/toast'
import Text from '_components/text'
import Flex from '_components/flex'
import Box from '_components/box'

import NetworkCard from './network-card'

const Container = styled(Flex)`
  flex-direction: column;
`

const NetworksView = () => {
  const [searchFilter, setSearchFilter] = useState('')

  const location = useLocation()

  const getNetworksQuery = useQuery(GET_NETWORKS)

  const [deleteNetwork, deleteNetworkMutation] = useMutation(DELETE_NETWORK, {
    variables: { id: undefined },
    onCompleted: () => {
      toast.success('Network deleted')
      getNetworksQuery.refetch()
    },
    onError: () => {
      toast.error('Error deleting host')
    },
  })

  useEffect(() => {
    window.scrollTo(0, 0)
    getNetworksQuery.refetch()
  }, [location])

  const handleNetworkSelect = id => {
    navigate(`networks/${id}`)
  }

  const handleNetworkDelete = id => {
    deleteNetwork({ variables: { id } })
  }

  const handleCreateNetworkClick = () => {
    navigate('networks/new')
  }

  const isError = getNetworksQuery.error
  const isLoading = getNetworksQuery.loading || deleteNetworkMutation.loading
  const isEmpty = !isError && !isLoading && getNetworksQuery.data.result.items.length === 0

  const filteredNetworks =
    isError || isLoading
      ? []
      : getNetworksQuery.data.result.items.filter(
          el => el.name.includes(searchFilter) || el.ipAddressRange.includes(searchFilter)
        )

  if (isLoading) {
    return <Spinner />
  }

  return (
    <Container>
      <Box mb={3}>
        <Text textStyle="title">Networks</Text>
        <Button onClick={handleCreateNetworkClick} variant="primary" ml="auto">
          Create
        </Button>
      </Box>
      <Box my={3}>
        <SearchInput
          width="100%"
          placeholder="Search..."
          onChange={e => setSearchFilter(e.target.value)}
          mr={2}
        />
      </Box>
      {isError ? (
        <ErrorState message="" />
      ) : isEmpty ? (
        <EmptyState message="Oops! It seems that you don't have any networks yet registered." />
      ) : (
        filteredNetworks.map(n => (
          <NetworkCard
            key={n.id}
            id={n.id}
            name={n.name}
            ipAddressRange={n.ipAddressRange}
            numHosts={n.hosts.count}
            onClick={handleNetworkSelect}
            onDelete={() => handleNetworkDelete(n.id)}
          />
        ))
      )}
    </Container>
  )
}

export default NetworksView
