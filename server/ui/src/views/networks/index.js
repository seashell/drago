import React, { useEffect } from 'react'
import styled from 'styled-components'
import { navigate } from '@reach/router'
import { useQuery, useMutation } from 'react-apollo'
import { GET_NETWORKS, DELETE_NETWORK } from '_graphql/actions'

import { Dragon as Spinner } from '_components/spinner'
import Text from '_components/text'
import Flex from '_components/flex'
import toast from '_components/toast'
import Button from '_components/button'
import Box from '_components/box'

import { icons } from '_assets/'

import NetworksList from './networks-list'

const Container = styled(Flex)`
  flex-direction: column;
`

const ErrorStateContainer = styled(Box).attrs({
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

const ErrorState = () => (
  <ErrorStateContainer>
    <icons.ErrorStateCube />
    <Text textStyle="description" mt={4}>
      Oops! It seems that an error has occurred.
    </Text>
  </ErrorStateContainer>
)

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
      Oops! It seems that there are no networks registered.
    </Text>
  </EmptyStateContainer>
)

export const StyledButton = styled(Button)``

const NetworksView = () => {
  const getNetworksQuery = useQuery(GET_NETWORKS)

  useEffect(() => {
    getNetworksQuery.refetch()
  })

  const handleNetworkDeleted = () => {
    toast.success('Network deleted')
    getNetworksQuery.refetch()
  }

  const handleNetworkDeleteError = () => {
    toast.error('Error deleting host')
  }

  const [deleteNetwork, deleteNetworkMutation] = useMutation(DELETE_NETWORK, {
    variables: { id: undefined },
    onCompleted: handleNetworkDeleted,
    onError: handleNetworkDeleteError,
  })

  const handleNetworkSelect = id => {
    navigate(`/networks/${id}`)
  }

  const handleNetworkDelete = (e, id) => {
    e.preventDefault()
    e.stopPropagation()
    deleteNetwork({ variables: { id } })
  }

  const handleCreateNetworkClick = () => {
    navigate('/networks/new')
  }

  const isError = getNetworksQuery.error
  const isLoading = getNetworksQuery.loading || deleteNetworkMutation.loading
  const isEmpty = !isError && !isLoading && getNetworksQuery.data.result.items.length === 0

  return (
    <Container>
      <Box mb={3}>
        <Text textStyle="title">Networks</Text>
        <Button
          onClick={handleCreateNetworkClick}
          variant="primary"
          borderRadius={3}
          width="100px"
          height="40px"
          ml="auto"
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
        <NetworksList
          networks={getNetworksQuery.data.result.items}
          onNetworkSelect={handleNetworkSelect}
          onNetworkDelete={handleNetworkDelete}
        />
      )}
    </Container>
  )
}

export default NetworksView
