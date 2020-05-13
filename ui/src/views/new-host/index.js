import React from 'react'
import PropTypes from 'prop-types'
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
import { CREATE_HOST } from '_graphql/actions'
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
  ipAddress: value => {
    if (validator.isIPRange(value)) {
      return true
    }
    return 'Invalid IP address. Must be in the format <IP address/CIDR suffix>.'
  },
  advertiseAddress: value => {
    if (validator.isEmpty(value) || validator.isIP(value) || validator.isURL(value)) {
      return true
    }
    return 'Invalid advertise address. Must be an IP address or host name.'
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

const NewHost = ({ networkId }) => {
  const [formState, { text }] = useFormState({
    name: null,
    ipAddress: null,
    advertiseAddress: null,
  })

  const [createHost, { loading }] = useMutation(CREATE_HOST, {
    variables: { networkId },
    onCompleted: () => {
      toast.success('Host created successfully')
      navigate(-1)
    },
    onError: () => {
      toast.error('Error creating host')
      navigate(-1)
    },
  })

  const onSave = () => {
    createHost({ variables: { networkId, ...formState.values } })
  }

  return (
    <Container>
      <Text textStyle="title" mb={4}>
        New host
      </Text>
      {loading ? (
        <Spinner />
      ) : (
        <Box flexDirection="column">
          <form>
            <Text my={3}>Name</Text>
            <FormInput type={text} name="name" placeholder="device-1" formState={formState} />
            <Text my={3}>Address</Text>
            <FormInput
              type={text}
              name="ipAddress"
              placeholder="10.0.0.1/24"
              formState={formState}
            />
            <Text my={3}>Advertise address</Text>
            <FormInput
              type={text}
              name="advertiseAddress"
              placeholder="wg.domain.com"
              formState={formState}
            />
            <Button width="100%" borderRadius={3} mt={3} mb={4} onClick={onSave}>
              Save
            </Button>
          </form>
        </Box>
      )}

      <Box justifyContent="center" gridColumn="4 / span 6">
        <Link color="primary" to={`/networks/${networkId}/hosts`}>
          Cancel
        </Link>
      </Box>
    </Container>
  )
}

NewHost.propTypes = {
  networkId: PropTypes.string.isRequired,
}

export default NewHost
