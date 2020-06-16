import React, { useEffect, useState } from 'react'
import styled from 'styled-components'

import { navigate, useLocation } from '@reach/router'
import { useMutation, useQuery } from 'react-apollo'

import { DELETE_HOST, GET_HOSTS } from '_graphql/actions/hosts'

import { Dragon as Spinner } from '_components/spinner'
import SearchInput from '_components/inputs/search-input'
import ErrorState from '_components/error-state'
import EmptyState from '_components/empty-state'
import Button from '_components/button'
import toast from '_components/toast'
import Text from '_components/text'
import Box from '_components/box'

import HostsList from './hosts-list'

const Container = styled(Box)`
  flex-direction: column;
`

const HostsListView = () => {
  const [searchFilter, setSearchFilter] = useState('')

  const location = useLocation()

  const getHostsQuery = useQuery(GET_HOSTS, {
    variables: {},
  })

  const [deleteHost, deleteHostMutation] = useMutation(DELETE_HOST, {
    variables: { id: undefined },
    onCompleted: () => {
      toast.success('Host deleted')
      getHostsQuery.refetch()
    },
    onError: () => {
      toast.error('Error deleting host')
    },
  })

  useEffect(() => {
    window.scrollTo(0, 0)
    getHostsQuery.refetch()
  }, [location])

  const handleHostSelect = id => {
    navigate(`hosts/${id}`)
  }

  const handleHostDelete = (e, id) => {
    e.preventDefault()
    e.stopPropagation()
    deleteHost({ variables: { id } })
  }

  const handleCreateHostButtonClick = () => {
    navigate(`hosts/new`)
  }

  const isError = getHostsQuery.error
  const isLoading = getHostsQuery.loading || deleteHostMutation.loading
  const isEmpty = !isLoading && !isError && getHostsQuery.data.result.items.length === 0

  const filteredHosts = isLoading
    ? []
    : getHostsQuery.data.result.items.filter(el => el.name.includes(searchFilter))

  return (
    <Container>
      <Box mb={3}>
        <Text textStyle="title">Hosts</Text>
        <Button onClick={handleCreateHostButtonClick} variant="primary" ml="auto">
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
        <ErrorState />
      ) : isLoading ? (
        <Spinner />
      ) : isEmpty ? (
        <EmptyState description="Oops! It seems that you don't have any hosts yet registered in this network." />
      ) : (
        <HostsList
          hosts={filteredHosts}
          onHostSelect={handleHostSelect}
          onHostDelete={handleHostDelete}
        />
      )}
    </Container>
  )
}

HostsListView.propTypes = {}

export default HostsListView
