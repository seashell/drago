import React, { useEffect } from 'react'
import styled from 'styled-components'

import { navigate } from '@reach/router'

import { useQuery, useMutation } from 'react-apollo'

import { GET_HOSTS, DELETE_HOST } from '_graphql/actions'
import { Dragon as Spinner } from '_components/spinner'
import Text from '_components/text'
import Flex from '_components/flex'
import toast from '_components/toast'
import Button from '_components/button'
import Box from '_components/box'
import { icons } from '_assets/'

import HostsList from './hosts-list'

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

export const StyledButton = styled(Button)``

const HostsView = () => {
  const { loading, data, error, refetch } = useQuery(GET_HOSTS)

  useEffect(() => {
    refetch()
  })

  const handleHostDeleted = () => {
    toast.success('Host deleted')
    refetch()
  }

  const handleHostDeleteError = () => {
    toast.error('Error deleting host')
  }

  const [deleteHost, { loading: deleting }] = useMutation(DELETE_HOST, {
    variables: { id: undefined },
    onCompleted: handleHostDeleted,
    onError: handleHostDeleteError,
  })

  const handleHostSelect = id => {
    navigate(`/hosts/${id}`)
  }

  const handleHostDelete = (e, id) => {
    e.preventDefault()
    e.stopPropagation()
    deleteHost({ variables: { id } })
  }

  const handleCreateHostClick = () => {
    navigate('/hosts/new')
  }

  return (
    <Container>
      <Box mb={3}>
        <Text textStyle="title">Hosts</Text>
        <Button
          onClick={handleCreateHostClick}
          variant="primary"
          borderRadius={3}
          width="100px"
          height="40px"
          ml="auto"
        >
          Create
        </Button>
      </Box>
      {error ? (
        <ErrorState />
      ) : loading || deleting ? (
        <Spinner />
      ) : (
        <HostsList
          hosts={data.result.items}
          onHostSelect={handleHostSelect}
          onHostDelete={handleHostDelete}
        />
      )}
    </Container>
  )
}

export default HostsView
