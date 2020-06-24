import React from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'

import Box from '_components/box'
import Text from '_components/text'

import { illustrations } from '_assets/'

const ErrorStateContainer = styled(Box).attrs({
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

const ErrorState = ({ title, description }) => (
  <ErrorStateContainer>
    <illustrations.Error />
    <Text textStyle="subtitle" mt={4}>
      {title}
    </Text>
    <Text textStyle="description" mt={4}>
      {description}
    </Text>
  </ErrorStateContainer>
)

ErrorState.propTypes = {
  title: PropTypes.string,
  description: PropTypes.string,
}

ErrorState.defaultProps = {
  title: 'Something is not right',
  description: 'It seems that an error has occurred. Check your network and try again.',
}

export default ErrorState
