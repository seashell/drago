import PropTypes from 'prop-types'
import React from 'react'
import styled from 'styled-components'
import { illustrations } from '_assets/'
import Box from '_components/box'
import Text from '_components/text'

const EmptyStateContainer = styled(Box).attrs({
  border: 'thin',
})`
  > svg {
    height: 160px;
    width: auto;
  }
  padding: 48px;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: ${(props) => props.theme.colors.white};
  border-color: ${(props) => props.theme.colors.neutralLighter};
`

const EmptyState = ({ image, title, description, extra, ...props }) => (
  <EmptyStateContainer {...props}>
    {image}
    <Text textStyle="subtitle" mt={3}>
      {title}
    </Text>
    <Text textStyle="body" textAlign="center" mt={3} color="neutral">
      {description}
    </Text>
    {extra}
  </EmptyStateContainer>
)

EmptyState.propTypes = {
  image: PropTypes.node,
  title: PropTypes.string,
  description: PropTypes.string,
  extra: PropTypes.node,
}

EmptyState.defaultProps = {
  image: <illustrations.Empty />,
  title: 'No results found',
  description: `Oops! It seems that this query has not returned any resources.`,
  extra: null,
}

export default EmptyState
