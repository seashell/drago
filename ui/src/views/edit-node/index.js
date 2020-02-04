import React from 'react'
import PropTypes from 'prop-types'
import { useFormState } from 'react-use-form-state'

import Box from '_components/box'
import Link from '_components/link'
import Text from '_components/text'
import Button from '_components/button'
import TextInput from '_components/inputs/text-input'
import SelectInput from '_components/inputs/select-input'
import { useQuery } from 'react-apollo'
import { GET_NODE, GET_NODES } from '_graphql/actions'

const EditNode = ({ nodeId }) => {
  const allNodes = useQuery(GET_NODES)

  const { loading, error, data } = useQuery(GET_NODE, {
    variables: { nodeId },
  })

  const [formState, { text }] = useFormState({
    nodeLabel: 'data.result.label',
    nodeAddress: 'data.result.address',
  })

  const peers = []
  const handleSelectInputChange = idx => {
    peers.push(allNodes.data.result.items[idx])
  }

  return (
    <>
      <Box>
        <Text textStyle="title">Edit node</Text>
      </Box>

      <Box padding={5} gridColumn="4 / span 6" flexDirection="column">
        <TextInput {...text('nodeLabel')} placeholder="Node label" mb={2} />
        <TextInput {...text('nodeAddress')} placeholder="10.0.8.0/24" mb={2} />
        <Text textStyle="subtitle" fontSize="18px" my={3}>
          Peers
        </Text>
        <SelectInput
          onChange={handleSelectInputChange}
          options={
            loading
              ? []
              : allNodes.data.result.items.map(n => ({
                  label: n.label,
                  value: n.id,
                }))
          }
          mb={5}
        />
        <Button width="100%" borderRadius={3}>
          Save
        </Button>
      </Box>

      <Box justifyContent="center" gridColumn="4 / span 6">
        <Link color="primary" to="/nodes">
          Cancel
        </Link>
      </Box>
    </>
  )
}

EditNode.propTypes = {
  nodeId: PropTypes.string.isRequired,
}

export default EditNode
