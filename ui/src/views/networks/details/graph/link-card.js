import PropTypes from 'prop-types'
import React from 'react'
import styled from 'styled-components'
import { icons } from '_assets/'
import Box from '_components/box'
import IconButton from '_components/icon-button'
import TextInput from '_components/inputs/text-input'
import Text from '_components/text'

const Container = styled(Box).attrs({
  display: 'flex',
  border: 'discrete',
  m: 2,
  p: 3,
})`
  flex-direction: column;
  width: 300px;
  height: 400px;
  box-shadow: ${(props) => props.theme.shadows.medium};
`

const IconContainer = styled(Box).attrs({
  display: 'flex',
  height: '72px',
  width: '72px',
  bg: 'neutralLighter',
  borderRadius: '4px',
})`
  position: relative;
  button {
    margin-right: auto;
  }
  align-items: center;
  justify-content: center;
  svg {
    transform: scale(1.6);
  }
`

const LinkCard = ({
  id,
  sourceName,
  sourceAddress,
  targetName,
  targetAddress,
  allowedIPs,
  persistentKeepalive,
  ...props
}) => (
  <Container {...props}>
    <IconContainer mx="auto" mt={2} mb={3}>
      <IconButton ml="auto" icon={<icons.Connection />} />
    </IconContainer>
    <Text textStyle="detail" my={1}>
      FROM
    </Text>
    <TextInput height={1} value={`${sourceName} (${sourceAddress})`} mb={1} />
    <Text textStyle="detail" my={1}>
      TO
    </Text>
    <TextInput height={1} value={`${targetName} (${targetAddress})`} mb={1} />
    <Text textStyle="detail" my={1}>
      ALLOWED IPS
    </Text>
    <TextInput height={1} value={allowedIPs} mb={1} />
    <Text textStyle="detail" my={1}>
      PERSISTENT KEEPALIVE
    </Text>
    <TextInput height={1} value={persistentKeepalive} mb={1} />
  </Container>
)

LinkCard.propTypes = {
  id: PropTypes.number.isRequired,
  sourceName: PropTypes.string.isRequired,
  sourceAddress: PropTypes.string.isRequired,
  targetName: PropTypes.string.isRequired,
  targetAddress: PropTypes.string.isRequired,
  allowedIPs: PropTypes.string,
  persistentKeepalive: PropTypes.string,
}

LinkCard.defaultProps = {
  allowedIPs: undefined,
  persistentKeepalive: undefined,
}

export default LinkCard
