import React from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'

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
  display: 'block',
  height: '48px',
  width: '48px',
  bg: 'neutralLighter',
  borderRadius: '4px',
})`
  position: relative;
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

const NodeCard = ({ id, label, address, onClick, onDelete }) => (
  <Container onClick={() => onClick(id)}>
    <IconContainer mr="12px">
      <StatusBadge status="online" />
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

NodeCard.propTypes = {
  id: PropTypes.number.isRequired,
  label: PropTypes.string.isRequired,
  address: PropTypes.string.isRequired,
  onClick: PropTypes.func,
  onDelete: PropTypes.func,
}

NodeCard.defaultProps = {
  onClick: () => {},
  onDelete: () => {},
}

export default NodeCard
