import React from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'

import Box from '_components/box'
import Text from '_components/text'
import IconButton from '_components/icon-button'
import { icons } from '_assets/'
import TextInput from '_components/inputs/text-input'

const Container = styled(Box).attrs({
  display: 'flex',
  border: 'discrete',
  m: 2,
  p: 3,
})`
  flex-direction: column;
  width: 300px;
  height: 400px;
  box-shadow: ${props => props.theme.shadows.medium};
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
    transform: scale(1.2);
  }
`

const StatusBadge = styled.div`
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: 4px solid white;
  position: absolute;
  right: -2px;
  bottom: -2px;
  background: ${props => (props.status === 'online' ? 'green' : 'red')};
`

const HostCard = ({ id, name, address, advertiseAddress, listenPort, ...props }) => (
  <Container {...props}>
    <IconContainer mx="auto" mt={2} mb={3}>
      <IconButton ml="auto" icon={<icons.Cube />} />
      <StatusBadge status="online" />
    </IconContainer>
    <Box flexDirection="column">
      <Text textStyle="detail" my={1}>
        NAME
      </Text>
      <TextInput height={1} value={name} mb={1} />
      <Text textStyle="detail" my={1}>
        ADDRESS
      </Text>
      <TextInput height={1} value={address} mb={1} />
      <Text textStyle="detail" my={1}>
        ADVERTISE ADDRESS
      </Text>
      <TextInput height={1} value={advertiseAddress} placeholder="N/A" mb={1} />
      <Text textStyle="detail" my={1}>
        LISTEN PORT
      </Text>
      <TextInput height={1} value={listenPort} placeholder="N/A" mb={1} />
    </Box>
  </Container>
)

HostCard.propTypes = {
  id: PropTypes.number.isRequired,
  name: PropTypes.string.isRequired,
  address: PropTypes.string.isRequired,
  advertiseAddress: PropTypes.string.isRequired,
  listenPort: PropTypes.string.isRequired,
}

HostCard.defaultProps = {}

export default HostCard
