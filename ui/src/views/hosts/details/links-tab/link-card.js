import React, { useState } from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'

import { icons } from '_assets/'
import Box from '_components/box'
import Text from '_components/text'
import IconButton from '_components/icon-button'
import Button from '_components/button'

const Container = styled(Box).attrs({
  border: 'discrete',
  m: 1,
})`
  display: flex;
  flex-direction: column;
  cursor: pointer;
  :hover {
    box-shadow: ${props => props.theme.shadows.medium};
  }
`

const SummaryContainer = styled(Box).attrs({
  p: 3,
})`
  align-items: center;
  height: 40px;
  display: flex;
`

const DetailsContainer = styled(Box)`
  height: 100px;
  padding: 8px;
  position: relative;
  display: ${props => (props.isVisible ? 'flex' : 'none')};
  button {
    position: absolute;
    right: 0px;
    bottom: 0px;
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
  const [isExpanded, setExpanded] = useState(false)

  const handleDeleteButtonClick = e => {
    e.stopPropagation()
    e.preventDefault()
    onDelete()
  }

  const handleCardClick = () => {
    setExpanded(!isExpanded)
  }

  const handleEditButtonClick = e => {
    e.stopPropagation()
    e.preventDefault()
    onClick()
  }

  return (
    <Container onClick={handleEditButtonClick}>
      <SummaryContainer>
        <Box width="100%">
          <IconContainer mr="12px">
            <IconButton ml="auto" icon={<icons.Link />} />
          </IconContainer>
          <div>
            <Box alignItems="center">
              <Text textStyle="subtitle" fontSize="14px">
                From
              </Text>
              <Text ml={1} textStyle="detail" fontSize="12px">
                {fromInterface.name}
              </Text>
              <Text ml={1} textStyle="subtitle" fontSize="14px">
                on
              </Text>
              <Text ml={1} textStyle="detail" fontSize="12px">
                {fromInterface.host.name}
              </Text>
            </Box>
            <Box alignItems="center">
              <Text textStyle="subtitle" fontSize="14px">
                To
              </Text>
              <Text ml={1} textStyle="detail" fontSize="12px">
                {toInterface.name}
              </Text>
              <Text ml={1} textStyle="subtitle" fontSize="14px">
                on
              </Text>
              <Text ml={1} textStyle="detail" fontSize="12px">
                {toInterface.host.name}
              </Text>
            </Box>
          </div>
        </Box>
        <Box flexDirection="column" ml="auto" alignItems="center">
          <Text textStyle="detail">Allowed IPs</Text>
          <Box width="250px" justifyContent="center">
            <Text textStyle="detail">
              {allowedIps.map(el => (
                <AllowedIPBlock>{el}</AllowedIPBlock>
              ))}
            </Text>
          </Box>
        </Box>
        <Box width="150px" flexDirection="column" alignItems="center">
          <Text textStyle="detail">Keepalive</Text>
          <Text textStyle="detail">{persistentKeepalive}</Text>
        </Box>
        <Box>
          <IconButton ml="auto" icon={<icons.Times />} onClick={handleDeleteButtonClick} />
        </Box>
      </SummaryContainer>
      <DetailsContainer isVisible={isExpanded}>
        <Button variant="primaryInverted" width="80px" onClick={handleEditButtonClick}>
          Edit
        </Button>
      </DetailsContainer>
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
