import React, { useState } from 'react'
import styled from 'styled-components'
import * as jwtDecode from 'jwt-decode'

import Text from '_components/text'
import Flex from '_components/flex'
import Button from '_components/button'
import Box from '_components/box'
import TextInput from '_components/inputs/text-input'

import { useFormState } from 'react-use-form-state'
import { useLocalStorage } from 'react-use'

const Container = styled(Flex)`
  flex-direction: column;
`

const StyledButton = styled(Button).attrs({
  height: '40px',
  borderRadius: 3,
  mr: 3,
})`
  display: block;
  white-space: nowrap;
`

const ImportantBox = styled(Box)`
  flex-direction: row;
  align-items: center;
`

const TokensView = () => {
  const [formState, { text }] = useFormState({
    token: undefined,
  })
  const [token, setToken] = useLocalStorage('drago.settings.acl.token')

  const handleSetTokenButtonClick = () => {
    setToken(formState.values.token)
  }

  const handleClearTokenButtonClick = () => {
    setToken(undefined)
    formState.clear()
  }

  let decodedToken
  let hasError
  if (token !== undefined) {
    try {
      decodedToken = jwtDecode(token)
      hasError = false
    } catch (error) {
      hasError = true
    }
  }

  return (
    <Container>
      <Box mb={3} mr="auto">
        <Text textStyle="title">Access control tokens</Text>
      </Box>
      <Text textStyle="bodyText">
        By providing a token Secret ID, each future request will be authenticated, potentially
        authorizing read access to additional information. By providing a token Accessor ID, the
        policies and rules for the token will be listed.
      </Text>
      <ImportantBox my={3} p={3} border="discrete">
        <div>
          <Text textStyle="subtitle" mb={2} mr="auto" fontSize="18px">
            Token storage
          </Text>
          <Text textStyle="bodyText">
            Tokens are stored client-side in local storage. This will persist your token across
            sessions. You can manually clear your token here.
          </Text>
        </div>
        <StyledButton variant="primaryInverted" mx={4} onClick={handleClearTokenButtonClick}>
          Clear token
        </StyledButton>
      </ImportantBox>

      {decodedToken === undefined ? (
        <>
          <Text my={3}>Secret ID</Text>
          <TextInput
            required
            {...text('token')}
            placeholder="XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
            mb={2}
          />
          <Text textStyle="detail">Sent with every request to determine authorization</Text>
          <StyledButton my={3} mb={4} px={2} onClick={handleSetTokenButtonClick}>
            Set token
          </StyledButton>
        </>
      ) : (
        <ImportantBox my={3} p={3} bg="green">
          <div>
            <Text textStyle="subtitle" mb={2} mr="auto" fontSize="18px">
              Token authenticated!
            </Text>
            <Text textStyle="bodyText">
              Your token is valid and authorized for the following policies.
            </Text>
          </div>
        </ImportantBox>
      )}

      {hasError && (
        <ImportantBox my={3} p={3} bg="danger">
          <div>
            <Text textStyle="subtitle" mb={2} mr="auto" fontSize="18px" color="white">
              Token invalid!
            </Text>
            <Text textStyle="bodyText" color="white">
              The secret you have provided is either invalid or does not match an existing token.
            </Text>
          </div>
        </ImportantBox>
      )}

      {decodedToken !== undefined && (
        <>
          <Text textStyle="subtitle" fontSize="18px" mb={2}>
            Token: {decodedToken.sub}
          </Text>
          <Text textStyle="bodyText">Accessor: {decodedToken.sub}</Text>
          <Text textStyle="bodyText">Secret: {decodedToken.sub}</Text>
          <Text textStyle="subtitle" fontSize="18px" mt={3}>
            Policies
          </Text>
          <ImportantBox my={3} p={3} border="discrete" />
        </>
      )}
    </Container>
  )
}

export default TokensView
