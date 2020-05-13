import React from 'react'
import styled from 'styled-components'

import { useFormState } from 'react-use-form-state'

import Flex from '_components/flex'
import Box from '_components/box'
import Link from '_components/link'
import Text from '_components/text'
import Button from '_components/button'
import TextInput from '_components/inputs/text-input'
import { Dragon as Spinner } from '_components/spinner'

import { useMutation } from 'react-apollo'
import { CREATE_NETWORK } from '_graphql/actions'
import { navigate } from '@reach/router'
import toast from '_components/toast'

const Container = styled(Flex)`
  flex-direction: column;
`

const NewNetwork = () => {
  const [formState, { text }] = useFormState({
    name: undefined,
    ipAddressRange: undefined,
  })

  const onNetworkCreated = () => {
    toast.success('Network created')
    navigate('/networks')
  }

  const onNetworkCreationError = () => {
    toast.error('Error creating network')
    navigate('/networks')
  }

  const [createNetwork, { loading }] = useMutation(CREATE_NETWORK, {
    variables: formState.values,
    onCompleted: onNetworkCreated,
    onError: onNetworkCreationError,
  })

  const onSave = () => {
    createNetwork()
  }

  return (
    <Container>
      <Text textStyle="title" mb={4}>
        New network
      </Text>
      {loading ? (
        <Spinner />
      ) : (
        <Box flexDirection="column">
          <Text my={3}>Name</Text>
          <TextInput required {...text('name')} placeholder="new-network" mb={2} />
          <Text my={3}>Address range</Text>
          <TextInput required {...text('ipAddressRange')} placeholder="10.0.8.0/24" mb={2} />
          <Button width="100%" borderRadius={3} mt={3} mb={4} onClick={onSave}>
            Save
          </Button>
        </Box>
      )}

      <Box justifyContent="center" gridColumn="4 / span 6">
        <Link color="primary" to="/networks">
          Cancel
        </Link>
      </Box>
    </Container>
  )
}

export default NewNetwork
