import React from 'react'
import styled from 'styled-components'
import validator from 'validator'
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
import { icons } from '_assets/'

const Container = styled(Flex)`
  flex-direction: column;
`

const validators = {
  name: value => {
    if (validator.isAscii(value)) {
      return true
    }
    return 'Invalid name. Must be non-empty and contain only ASCII characters.'
  },
  ipAddressRange: value => {
    if (validator.isIPRange(value)) {
      return true
    }
    return 'Invalid IP address range. Must be non-empty and in the format <IP address/CIDR suffix>.'
  },
}

// eslint-disable-next-line react/prop-types
const FormInput = ({ type, name, placeholder, formState }) => {
  const validityIndicator = () => {
    const isValid = formState.validity[name]
    const error = formState.errors[name]
    const IconContainer = styled.div.attrs({
      title: error,
    })`
      position: absolute;
      right: 12px;
      top: 50%;
      transform: translateY(-50%);
      cursor: help;
    `
    if (error === undefined) return null
    return <IconContainer>{isValid ? <icons.Success /> : <icons.Error />}</IconContainer>
  }
  return (
    <Box style={{ position: 'relative' }} mb={2}>
      <TextInput
        {...type({
          name,
          validate: validators[name],
          validateOnBlur: true,
        })}
        placeholder={placeholder}
      />
      {validityIndicator(name)}
    </Box>
  )
}

const NewNetwork = () => {
  const [formState, { text }] = useFormState({
    name: null,
    ipAddressRange: null,
  })

  const onNetworkCreated = () => {
    toast.success('Network created successfully')
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
          <form>
            <Text my={3}>Name</Text>
            <FormInput type={text} name="name" placeholder="new-network" formState={formState} />
            <Text my={3}>Address range</Text>
            <FormInput
              type={text}
              name="ipAddressRange"
              placeholder="10.0.8.0/24"
              formState={formState}
            />
            <Button width="100%" borderRadius={3} mt={3} mb={4} onClick={onSave}>
              Save
            </Button>
          </form>
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
