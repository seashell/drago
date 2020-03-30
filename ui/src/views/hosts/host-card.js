import React from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'
import moment from 'moment'

import Box from '_components/box'
import Text from '_components/text'
import IconButton from '_components/icon-button'
import { icons } from '_assets/'

const Container = styled(Box).attrs({
  display: 'flex',
  border: 'discrete',
  m: 2,
  p: 3,
})`
  align-items: center;
  cursor: pointer;
  :hover {
    box-shadow: ${props => props.theme.shadows.medium};
  }
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

const StatusBadge = styled.div`
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: 4px solid white;
  position: absolute;
  right: -2px;
  bottom: -2px;
  background: ${props => (props.status === 'online' ? 'green' : props.theme.colors.neutralLight)};
`

const HostCard = ({ id, label, address, lastSeen, onClick, onDelete }) => {
  const isOnline = Math.abs(moment(lastSeen).diff(moment.now(), 'minutes')) < 5

  return (
    <Container onClick={() => onClick(id)}>
      <IconContainer mr="12px">
        <IconButton ml="auto" icon={<icons.Cube />} />
        <StatusBadge status={isOnline ? 'online' : 'offline'} />
      </IconContainer>
      <Box flexDirection="column">
        <Text textStyle="subtitle" fontSize="14px">
          {label}
        </Text>
        <Text textStyle="detail" fontSize="12px">
          {address}
        </Text>
      </Box>
      <IconButton ml="auto" icon={<icons.Times />} onClick={onDelete} />
    </Container>
  )
}

HostCard.propTypes = {
  id: PropTypes.number.isRequired,
  label: PropTypes.string.isRequired,
  address: PropTypes.string.isRequired,
  lastSeen: PropTypes.string.isRequired,
  onClick: PropTypes.func,
  onDelete: PropTypes.func,
}

HostCard.defaultProps = {
  onClick: () => {},
  onDelete: () => {},
}

export default HostCard
