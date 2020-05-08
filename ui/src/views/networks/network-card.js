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
  m: 1,
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
  height: '36px',
  width: '36px',
  bg: 'neutralLighter',
  borderRadius: '4px',
})`
  position: relative;
  button {
    margin-right: auto;
  }
  svg {
    transform: scale(1.2);
  }
  align-items: center;
  justify-content: center;
`

const NetworkCard = ({ id, name, hostCount, onClick, onDelete }) => (
  <Container onClick={() => onClick(id)}>
    <IconContainer mr="12px">
      <IconButton ml="auto" icon={<icons.Network />} />
    </IconContainer>
    <Box flexDirection="column">
      <Text textStyle="subtitle" fontSize="14px">
        {name}
      </Text>
      <Text textStyle="detail" fontSize="12px">
        {hostCount} hosts
      </Text>
    </Box>
    <IconButton ml="auto" icon={<icons.Times />} onClick={onDelete} />
  </Container>
)

NetworkCard.propTypes = {
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  hostCount: PropTypes.number.isRequired,
  onClick: PropTypes.func,
  onDelete: PropTypes.func,
}

NetworkCard.defaultProps = {
  onClick: () => {},
  onDelete: () => {},
}

export default NetworkCard
