import React, { useEffect } from 'react'
import styled from 'styled-components'

import { navigate } from '@reach/router'

import Text from '_components/text'
import Flex from '_components/flex'
import Button from '_components/button'
import Box from '_components/box'

const Container = styled(Flex)`
  flex-direction: column;
  align-items: center;
  text-align: center;
`

export const StyledButton = styled(Button).attrs({
  variant: 'primary',
  height: '40px',
  width: '100px',
  borderRadius: 3,
  mr: 3,
})`
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 8px;
`

const IconContainer = styled(Box).attrs({
  display: 'flex',
  height: '48px',
  width: '48px',
  bg: 'neutralLighter',
  borderRadius: '4px',
})`
  position: relative;
  button {
    margin-right: auto;
  }
  align-items: center;
  justify-content: center;
`

const HomeView = () => {
  useEffect(() => {})

  const handleHostsButtonClicked = () => {
    navigate(`/hosts`)
  }

  return (
    <Container>
      <Box mb={3}>
        <Text textStyle="title">Welcome to Drago!</Text>
      </Box>

      <Text textStyle="bodyText">
        Drago is a flexible configuration manager for WireGuard networks. It is designed to make it
        simple to define and manage secure network overlays spanning heterogeneous hosts distributed
        across different clouds and edge locations.
      </Text>

      <Text my={3} textStyle="bodyText">
        If you have questions, feature requests, or issues to report, please reach out through{' '}
        <a href="https://github.com/seashell/drago/issues">our Github repository</a>.
      </Text>

      <StyledButton mt={3} onClick={handleHostsButtonClicked}>
        Get started
      </StyledButton>
    </Container>
  )
}

export default HomeView
