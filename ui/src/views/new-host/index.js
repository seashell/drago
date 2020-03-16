import React from 'react'
import styled from 'styled-components'

import { useFormState } from 'react-use-form-state'

import Flex from '_components/flex'
import Box from '_components/box'
import Link from '_components/link'
import Text from '_components/text'
import Button from '_components/button'
import TextInput from '_components/inputs/text-input'
import { Dragon } from '_components/spinner'

import { useMutation } from 'react-apollo'
import { CREATE_NODE } from '_graphql/actions'
import { navigate } from '@reach/router'
import toast from '_components/toast'

const Container = styled(Flex)`
  flex-direction: column;
`

const NewNode = () => {
  const [formState, { text }] = useFormState()

  const onNodeCreated = () => {
    toast.success('Node created')
    navigate('/nodes')
  }

  const onNodeCreationError = () => {
    toast.error('Error creating node')
    navigate('/nodes')
  }

  const [createNode, { loading }] = useMutation(CREATE_NODE, {
    variables: formState.values,
    onCompleted: onNodeCreated,
    onError: onNodeCreationError,
  })

  const onSave = () => {
    createNode()
  }

  return (
    <Container>
      <Text textStyle="title" mb={4}>
        New node
      </Text>
      {loading ? (
        <Dragon />
      ) : (
        <Box flexDirection="column">
          <Text my={3}>Name</Text>
          <TextInput required {...text('name')} placeholder="new-node-1" mb={2} />
          <Text my={3}>Address</Text>
          <TextInput required {...text('address')} placeholder="10.0.8.0/24" mb={2} />
          <Button width="100%" borderRadius={3} mt={3} mb={4} onClick={onSave}>
            Save
          </Button>
        </Box>
      )}

      <Box justifyContent="center" gridColumn="4 / span 6">
        <Link color="primary" to="/nodes">
          Cancel
        </Link>
      </Box>
    </Container>
  )
}

export default NewNode
