import React, { useEffect } from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'

import { navigate, useLocation } from '@reach/router'
import { useMutation, useQuery } from 'react-apollo'

import { DELETE_HOST, GET_HOSTS } from '_graphql/actions'

import { Dragon as Spinner } from '_components/spinner'
import ErrorState from '_components/error-state'
import EmptyState from '_components/empty-state'
import Button from '_components/button'
import toast from '_components/toast'
import Flex from '_components/flex'
import Text from '_components/text'
import Box from '_components/box'

import HostsList from './hosts-list'

const Container = styled(Flex)`
  flex-direction: column;
`

export const StyledButton = styled(Button)``

const HostsView = ({ networkId }) => {
  const location = useLocation()

  const getHostsQuery = useQuery(GET_HOSTS, {
    variables: { networkId },
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
        <EmptyState description="Oops! It seems that you don't have any hosts yet registered in this network." />
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
