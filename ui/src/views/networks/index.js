import React, { useEffect } from 'react'
import styled from 'styled-components'
import { navigate, useLocation } from '@reach/router'
import { useQuery, useMutation } from 'react-apollo'

import { GET_NETWORKS, DELETE_NETWORK } from '_graphql/actions'

import { Dragon as Spinner } from '_components/spinner'
import Button from '_components/button'
import toast from '_components/toast'
import Text from '_components/text'
import Flex from '_components/flex'
import Box from '_components/box'
import { illustrations } from '_assets/'

import NetworksList from './networks-list'

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
      Oops! It seems that there are no networks registered.
    </Text>
  </EmptyStateContainer>
)

export const StyledButton = styled(Button)``

const NetworksView = () => {
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
    navigate(`/networks/${id}/hosts`)
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
