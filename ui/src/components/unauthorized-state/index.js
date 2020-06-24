import React from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'

import Box from '_components/box'
import Text from '_components/text'

import { illustrations } from '_assets/'

const UnaurhorizedStateContainer = styled(Box).attrs({
  border: 'none',
  height: '300px',
})`
  svg {
    height: 300px;
    width: auto;
  }
  padding: 20px;
  flex-direction: column;
  align-items: center;
  justify-content: center;
`

const UnaurhorizedState = ({ title, description }) => (
  <UnaurhorizedStateContainer>
    <illustrations.Unaurhorized />
    <Text textStyle="subtitle" mt={4}>
      {title}
    </Text>
    <Text textStyle="description" mt={4}>
      {description}
    </Text>
  </UnaurhorizedStateContainer>
)

UnaurhorizedState.propTypes = {
  title: PropTypes.string,
  description: PropTypes.string,
}

UnaurhorizedState.defaultProps = {
  title: 'Forbidden',
  description: `Oops! It seems that you don't have enough permissions to access this resource`,
}

export default UnaurhorizedState
