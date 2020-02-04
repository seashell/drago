import React from 'react'
import { useFormState } from 'react-use-form-state'

import Box from '_components/box'
import Link from '_components/link'
import Text from '_components/text'
import Button from '_components/button'
import TextInput from '_components/inputs/text-input'
import SelectInput from '_components/inputs/select-input'
import { useQuery } from 'react-apollo'
import { GET_NODES } from '_graphql/actions'

const NewNode = () => {
  const [formState, { text }] = useFormState()

  const { loading, error, data } = useQuery(GET_NODES)

  return (
    <>
      <Box>
        <Text textStyle="title">New node</Text>
      </Box>

      <Box padding={5} gridColumn="4 / span 6" flexDirection="column">
        <TextInput {...text('projectName')} placeholder="Node name" mb={2} />
        <TextInput {...text('allowedIPs')} placeholder="10.0.8.0/24" mb={2} />
        <Text textStyle="subtitle" fontSize="18px" my={3}>
          Peers
        </Text>
        <SelectInput
          options={
            loading
              ? []
              : data.result.items.map(n => ({
                  label: n.label,
                  value: n.id,
                }))
          }
          mb={5}
        />
        <Button width="100%" borderRadius={3}>
          Continue
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

export default NewNode
