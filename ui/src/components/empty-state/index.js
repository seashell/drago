import React from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'

import Box from '_components/box'
import Text from '_components/text'

import { illustrations } from '_assets/'

const EmptyStateContainer = styled(Box).attrs({
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

const EmptyState = ({ title, description }) => (
  <EmptyStateContainer>
    <illustrations.Empty />
    <Text textStyle="subtitle" mt={4}>
      {title}
    </Text>
    <Text textStyle="description" mt={4}>
      {description}
    </Text>
  </EmptyStateContainer>
)

EmptyState.propTypes = {
  title: PropTypes.string,
  description: PropTypes.string,
}

EmptyState.defaultProps = {
  title: 'No results found',
  description: `Oops! It seems that you don't have any resources yet registered.`,
}

export default EmptyState
