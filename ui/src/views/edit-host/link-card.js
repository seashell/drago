import React from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'

import Box from '_components/box'
import Text from '_components/text'
import IconButton from '_components/icon-button'
import TextInput from '_components/inputs/text-input'
import { icons } from '_assets/'

const Container = styled(Box).attrs({
  display: 'flex',
  border: 'discrete',
  m: 2,
  p: 3,
})`
  height: 200px;
  grid-column: span 1;
  flex-direction: column;
  position: relative;
  cursor: pointer;
  :hover {
    box-shadow: ${props => props.theme.shadows.medium};
    button {
      visibility: visible;
    }
  }
  button {
    visibility: hidden;
  }
`

const StyledIconButton = styled(IconButton)`
  height: 8px;
  width: 8px;
  position: absolute;
  right: 2px;
  top: 2px;
  :hover {
    background: ${props => props.theme.colors.neutralLighter};
  }
  svg {
    transform: scale(1.5);
  }
`

const LinkCard = ({
  id,
  toName,
  toAddress,
  allowedIPs,
  persistentKeepalive,
  onChange,
  onDelete,
}) => (
  <Container>
    <div>
      <Text textStyle="subtitle" fontSize="14px">
        {toName}
      </Text>
      <Text textStyle="detail" fontSize="14px">
        {toAddress}
      </Text>
    </div>
    <Box flexDirection="column" mt="auto">
      <Text textStyle="detail" my={1}>
        ALLOWED IPS
      </Text>
      <TextInput height={1} value={allowedIPs} placeholder="0.0.0.0/0" mb={1} />
      <Text textStyle="detail" my={1}>
        KEEPALIVE
      </Text>
      <TextInput height={1} value={persistentKeepalive} placeholder="20" mb={1} />
    </Box>
    <StyledIconButton ml="auto" icon={<icons.Times />} onClick={onDelete} />
  </Container>
)
LinkCard.propTypes = {
  id: PropTypes.number.isRequired,
  toName: PropTypes.string,
  toAddress: PropTypes.string,
  allowedIPs: PropTypes.string,
  persistentKeepalive: PropTypes.number,
  onChange: PropTypes.func,
  onDelete: PropTypes.func,
}

LinkCard.defaultProps = {
  toName: undefined,
  toAddress: undefined,
  allowedIPs: undefined,
  persistentKeepalive: undefined,
  onChange: () => {},
  onDelete: () => {},
}

export default LinkCard
