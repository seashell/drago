import React, { useEffect } from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'

import { navigate, useLocation } from '@reach/router'
import { useMutation, useQuery } from 'react-apollo'

import { DELETE_HOST, GET_HOSTS } from '_graphql/actions'

import { Dragon as Spinner } from '_components/spinner'
import Button from '_components/button'
import toast from '_components/toast'
import Flex from '_components/flex'
import Text from '_components/text'
import Box from '_components/box'

import { illustrations } from '_assets/'

import HostsList from './hosts-list'

const Container = styled(Flex)`
  flex-direction: column;
`

const ErrorStateContainer = styled(Box).attrs({
  border: 'none',
  height: '300px',
})`
  svg {
    height: 300px;
    width: auto;
  }
  padding: 20px;
  flex-direction: column;
  align-items: center;
  justify-content: center;
`

const ErrorState = () => (
  <ErrorStateContainer>
    <illustrations.Error />
    <Text textStyle="description" mt={4}>
      Oops! It seems that an error has occurred.
    </Text>
  </ErrorStateContainer>
)

const EmptyStateContainer = styled(Box).attrs({
  border: 'none',
  height: '300px',
})`
  svg {
    height: 300px;
    width: auto;
  }
  padding: 20px;
  flex-direction: column;
  align-items: center;
  justify-content: center;
`

const EmptyState = () => (
  <EmptyStateContainer>
    <illustrations.Empty />
    <Text textStyle="description" mt={4}>
      Oops! It seems that there are no hosts registered in this network.
    </Text>
  </EmptyStateContainer>
)

export const StyledButton = styled(Button)``

const HostsView = ({ networkId }) => {
  const location = useLocation()

  const getHostsQuery = useQuery(GET_HOSTS, {
    variables: { networkId },
  })

  const [deleteHost, deleteHostMutation] = useMutation(DELETE_HOST, {
    variables: { id: undefined },
    onCompleted: handleHostDeleted,
    onError: handleHostDeleteError,
  })

  useEffect(() => {
    window.scrollTo(0, 0)
    getHostsQuery.refetch()
  }, [location])

  const handleHostDeleted = () => {
    toast.success('Host deleted')
    getHostsQuery.refetch()
  }

  const handleHostDeleteError = () => {
    toast.error('Error deleting host')
  }

  const handleHostSelect = id => {
    navigate(`/networks/${networkId}/hosts/${id}`)
  }

  const handleHostDelete = (e, id) => {
    e.preventDefault()
    e.stopPropagation()
    deleteHost({ variables: { networkId, id } })
  }

  const handleGraphViewButtonClick = () => {
    navigate(`/networks/${networkId}/topology`)
  }

  const handleCreateHostButtonClick = () => {
    navigate(`/networks/${networkId}/hosts/new`)
  }

  const isError = getHostsQuery.error
  const isLoading = getHostsQuery.loading || deleteHostMutation.loading
  const isEmpty = !isLoading && !isError && getHostsQuery.data.result.items.length === 0

  return (
    <Container>
      <Box mb={3}>
        <Text textStyle="title">Hosts</Text>
        <Button
          onClick={handleGraphViewButtonClick}
          variant="primaryInverted"
          borderRadius={3}
          width="100px"
          height="40px"
          ml="auto"
        >
          Graph view
        </Button>
        <Button
          onClick={handleCreateHostButtonClick}
          variant="primary"
          borderRadius={3}
          width="100px"
          height="40px"
          ml={2}
        >
          Create
        </Button>
      </Box>
      {isError ? (
        <ErrorState />
      ) : isLoading ? (
        <Spinner />
      ) : isEmpty ? (
        <EmptyState />
      ) : (
        <HostsList
          hosts={getHostsQuery.data.result.items}
          onHostSelect={handleHostSelect}
          onHostDelete={handleHostDelete}
        />
      )}
    </Container>
  )
}

HostsView.propTypes = {
  networkId: PropTypes.string.isRequired,
}

export default HostsView
