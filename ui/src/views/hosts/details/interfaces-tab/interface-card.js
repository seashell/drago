import React from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'
import * as pluralize from 'pluralize'

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
  height: 40px;
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
  align-items: center;
  justify-content: center;
`

const KeyContainer = styled(Box).attrs({
  display: 'flex',
})`
  align-items: center;
  width: 20px;
  height: 20px;
  padding: 4px;
  border-radius: 50%;
`

const KeyIconContainer = styled(Box).attrs({
  display: 'flex',
  height: '20px',
  width: '20px',
})`
  position: relative;
  align-items: center;
  justify-content: center;
  fill: ${props => props.theme.colors.success};
`

const Badge = styled.div`
  color: ${props => props.theme.colors.neutralDark};
  background: ${props => props.theme.colors.neutralLighter};
  padding: 4px 8px;
  border-radius: 2px;
`

const InterfaceCard = ({ name, ipAddress, listenPort, publicKey, numLinks, onClick, onDelete }) => {
  const handleDeleteButtonClick = e => {
    e.preventDefault()
    e.stopPropagation()
    onDelete()
  }

  return (
    <Container onClick={onClick}>
      <Box width="200px" alignItems="center">
        <IconContainer mr="12px">
          <IconButton ml="auto" icon={<icons.Interface />} />
        </IconContainer>
        <div>
          <Text textStyle="subtitle" fontSize="14px">
            {name}
          </Text>
          <Text textStyle="detail">
            {ipAddress}
            {listenPort ? `:${listenPort}` : null}
          </Text>
        </div>
      </Box>
      {publicKey && (
        <KeyContainer ml={2}>
          <KeyIconContainer>
            <icons.Key />
          </KeyIconContainer>
        </KeyContainer>
      )}
      <Badge>
        <Text textStyle="detail">
          {numLinks} {pluralize('link', numLinks)}
        </Text>
      </Badge>
      <IconButton icon={<icons.Times />} onClick={handleDeleteButtonClick} />
    </Container>
  )
}

InterfaceCard.propTypes = {
  name: PropTypes.string,
  ipAddress: PropTypes.string,
  listenPort: PropTypes.string,
  publicKey: PropTypes.string,
  numLinks: PropTypes.number.isRequired,
  onClick: PropTypes.func,
  onDelete: PropTypes.func,
}

InterfaceCard.defaultProps = {
  name: '',
  ipAddress: '',
  listenPort: '',
  publicKey: '',
  onClick: () => {},
  onDelete: () => {},
}

export default InterfaceCard
