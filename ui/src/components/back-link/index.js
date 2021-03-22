import PropTypes from 'prop-types'
import React from 'react'
import styled from 'styled-components'
import { icons } from '_assets/'
import Box from '_components/box'
import Icon from '_components/icon'
import Link from '_components/link'
import Text from '_components/text'

const Container = styled(Box)`
  width: max-content;
  opacity: 0.3;
  :hover {
    opacity: 1;
  }
`

const BackLink = ({ text, to, ...props }) => (
  <Container {...props}>
    <Link to={to} width="max-content">
      <Box alignItems="center">
        <Icon icon={<icons.Back />} color="foreground2" size="14px" />
        <Text textStyle="body" ml={2}>
          {text}
        </Text>
      </Box>
    </Link>
  </Container>
)

BackLink.propTypes = {
  text: PropTypes.string.isRequired,
  to: PropTypes.string.isRequired,
}

export default BackLink
