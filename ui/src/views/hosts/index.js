import React, { useEffect } from 'react'
import styled from 'styled-components'

import { navigate } from '@reach/router'

import { useQuery, useMutation } from 'react-apollo'

import { GET_HOSTS, DELETE_HOST } from '_graphql/actions'
import { Dragon } from '_components/spinner'
import Text from '_components/text'
import Flex from '_components/flex'
import toast from '_components/toast'
import Button from '_components/button'
import Box from '_components/box'

import HostsList from './hosts-list'

const Container = styled(Flex)`
  flex-direction: column;
`

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
        <div>Error retrieving hosts</div>
      ) : loading || deleting ? (
        <Dragon />
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
