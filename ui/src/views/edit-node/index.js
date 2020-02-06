import React from 'react'
import PropTypes from 'prop-types'

import styled from 'styled-components'

import { useFormState } from 'react-use-form-state'

import Box from '_components/box'
import Flex from '_components/flex'
import Link from '_components/link'
import Text from '_components/text'
import Button from '_components/button'
import { Dragon } from '_components/spinner'
import TextInput from '_components/inputs/text-input'

import { useMutation, useQuery } from 'react-apollo'
import { GET_NODE, UPDATE_NODE } from '_graphql/actions'
import { navigate } from '@reach/router'
import toast from '_components/toast'
import Collapse from '_components/collapse'

const Container = styled(Flex)`
  flex-direction: column;
`

const EditNode = ({ nodeId }) => {
  const onNodeUpdated = () => {
    toast.success('Node updated')
    navigate('/nodes')
  }

  const onNodeUpdateError = () => {
    toast.error('Error updating node')
    navigate('/nodes')
  }

  const [formState, { text }] = useFormState()

  const query = useQuery(GET_NODE, {
    variables: { id: nodeId },
    onCompleted: data => {
      formState.setField('name', data.result.name)
      formState.setField('address', data.result.interface.address)
    },
  })

  const [updateNode, mutation] = useMutation(UPDATE_NODE, {
    variables: { id: nodeId, ...formState.values },
    onCompleted: onNodeUpdated,
    onError: onNodeUpdateError,
  })

  const onSave = () => {
    updateNode({ id: nodeId, ...formState.values })
  }

  return (
    <Container>
      <Text textStyle="title" mb={4}>
        Edit node
      </Text>
      {query.loading || mutation.loading ? (
        <Dragon />
      ) : (
        <Box flexDirection="column">
          <Text my={3}>Name</Text>
          <TextInput required {...text('name')} placeholder="new-node-1" mb={2} />

          <Collapse title={<Text textStyle="description">Control settings</Text>}>
            <Text my={3}>Public key</Text>
            <TextInput {...text('publicKey')} placeholder="N/A" mb={2} disabled />
            <Text my={3}>JWT</Text>
            <TextInput {...text('jwt')} placeholder="N/A" mb={2} disabled />
            <Text my={3}>Advertise address</Text>
            <TextInput {...text('advertiseAddr')} placeholder="wireguard.domain.io" mb={2} />
          </Collapse>

          <Collapse title={<Text textStyle="description">Interface settings</Text>}>
            <Text my={3}>Overlay address</Text>
            <TextInput required {...text('address')} placeholder="10.0.8.0/24" mb={2} />

            <Text my={3}>DNS</Text>
            <TextInput {...text('dns')} placeholder="8.8.8.8" mb={2} />

            <Text my={3}>MTU</Text>
            <TextInput {...text('mtu')} placeholder="1420" mb={2} />

            <Text my={3}>Pre up</Text>
            <TextInput {...text('preUp')} placeholder="" mb={2} />

            <Text my={3}>Post up</Text>
            <TextInput {...text('postUp')} placeholder="" mb={2} />

            <Text my={3}>Pre down</Text>
            <TextInput {...text('preDown')} placeholder="" mb={2} />

            <Text my={3}>Post down</Text>
            <TextInput {...text('postDown')} placeholder="" mb={2} />
          </Collapse>

          <Collapse title={<Text textStyle="description">Peers</Text>}>[TODO]</Collapse>

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

EditNode.propTypes = {
  nodeId: PropTypes.number,
}

EditNode.defaultProps = {
  nodeId: undefined,
}

export default EditNode
