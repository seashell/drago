import React from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'

import { icons } from '_assets/'
import Box from '_components/box'
import Text from '_components/text'
import IconButton from '_components/icon-button'

const Container = styled(Box).attrs({
  display: 'flex',
  border: 'discrete',
  m: 1,
  p: 3,
})`
  height: 40px;
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
  align-items: center;
  justify-content: center;
`

const AllowedIPBlock = styled.span`
  display: inline-block;
  background: white;
  margin: 0 4px;
`

const LinkCard = ({
  fromInterface,
  toInterface,
  allowedIps,
  persistentKeepalive,
  onClick,
  onDelete,
}) => {
  const handleDeleteButtonClick = e => {
    e.stopPropagation()
    e.preventDefault()
    onDelete()
  }

  return (
    <Container onClick={onClick}>
      <Box width="240px">
        <IconContainer mr="12px">
          <IconButton ml="auto" icon={<icons.Link />} />
        </IconContainer>
        <div>
          <Text textStyle="subtitle" fontSize="14px">
            From {fromInterface.name} on {fromInterface.host.name}
          </Text>
          <Text textStyle="subtitle" fontSize="14px">
            To {toInterface.name} on {toInterface.host.name}
          </Text>
        </div>
      </Box>
      <Box flexDirection="column" alignItems="center" width="300px">
        <Text textStyle="detail">Allowed IPs</Text>
        <Box width="250px" justifyContent="center">
          <Text textStyle="detail">
            {allowedIps.map(el => (
              <AllowedIPBlock>{el}</AllowedIPBlock>
            ))}
          </Text>
        </Box>
      </Box>
      <Box ml="auto" flexDirection="column" alignItems="center">
        <Text textStyle="detail">Keepalive</Text>
        <Text textStyle="detail">{persistentKeepalive}</Text>
      </Box>
      <IconButton ml="auto" icon={<icons.Times />} onClick={handleDeleteButtonClick} />
    </Container>
  )
}

LinkCard.propTypes = {
  fromInterface: PropTypes.string.isRequired,
  toInterface: PropTypes.string.isRequired,
  allowedIps: PropTypes.string,
  persistentKeepalive: PropTypes.number,
  onClick: PropTypes.func,
  onDelete: PropTypes.func,
}

LinkCard.defaultProps = {
  allowedIps: undefined,
  persistentKeepalive: undefined,
  onClick: () => {},
  onDelete: () => {},
}

export default LinkCard
