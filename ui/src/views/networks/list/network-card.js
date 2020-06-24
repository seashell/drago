import React from 'react'
import PropTypes from 'prop-types'
import * as pluralize from 'pluralize'

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
  justify-content: space-between;
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

const Badge = styled.div`
  color: ${props => props.theme.colors.neutralDark};
  background: ${props => props.theme.colors.neutralLighter};
  padding: 4px 8px;
  border-radius: 2px;
`

const NetworkCard = ({ id, name, ipAddressRange, numHosts, onClick, onDelete }) => {
  const handleDeleteButtonClick = e => {
    e.preventDefault()
    e.stopPropagation()
    onDelete()
  }
  return (
    <Container onClick={() => onClick(id)}>
      <Box>
        <IconContainer mr="12px">
          <IconButton ml="auto" icon={<icons.Network />} />
        </IconContainer>
        <div>
          <Text textStyle="subtitle" fontSize="14px">
            {name}
          </Text>
          <Text textStyle="detail" fontSize="12px">
            {ipAddressRange}
          </Text>
        </div>
      </Box>
      <Badge>
        <Text textStyle="detail">
          {numHosts} {pluralize('host', numHosts)}
        </Text>
      </Badge>
      <IconButton icon={<icons.Times />} onClick={handleDeleteButtonClick} />
    </Container>
  )
}

NetworkCard.propTypes = {
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  ipAddressRange: PropTypes.string.isRequired,
  numHosts: PropTypes.number.isRequired,
  onClick: PropTypes.func,
  onDelete: PropTypes.func,
}

NetworkCard.defaultProps = {
  onClick: () => {},
  onDelete: () => {},
}

export default NetworkCard
