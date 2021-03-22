import PropTypes from 'prop-types'
import React from 'react'
import styled from 'styled-components'
import { illustrations } from '_assets/'
import Box from '_components/box'
import Text from '_components/text'

const ErrorStateContainer = styled(Box)`
  svg {
    height: 160px;
    width: auto;
  }
  padding: 20px;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: ${(props) => props.theme.colors.background1};
`

const ErrorState = ({ title, description, extra }) => (
  <ErrorStateContainer>
    <illustrations.Error />
    <Text textStyle="subtitle" mt={4}>
      {title}
    </Text>
    <Text textStyle="body" mt={3} color="neutral">
      {description}
    </Text>
    {extra}
  </ErrorStateContainer>
)

ErrorState.propTypes = {
  title: PropTypes.string,
  description: PropTypes.string,
  extra: PropTypes.node,
}

ErrorState.defaultProps = {
  title: 'Something is not right',
  description: 'It seems that an error has occurred. Check your network and try again.',
  extra: null,
}

export default ErrorState
