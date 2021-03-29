import * as pluralize from 'pluralize'
import PropTypes from 'prop-types'
import React, { useState } from 'react'
import styled, { css } from 'styled-components'
import { icons } from '_assets/'
import Box from '_components/box'
import Button from '_components/button'
import Icon from '_components/icon'
import IconButton from '_components/icon-button'
import NumberInput from '_components/inputs/number-input'
import TextInput from '_components/inputs/text-input'
import Separator from '_components/separator'
import Text from '_components/text'

const Container = styled(Box).attrs({
  border: 'discrete',
})`
  flex-direction: column;
`

const HeaderContainer = styled(Box).attrs({
  px: 3,
})`
  height: 72px;
  cursor: pointer;
  align-items: center;
`

const HiddenContentContainer = styled(Box).attrs({
  p: 3,
})`
  flex-direction: column;
`

const StyledIcon = styled(Icon)`
  width: 36px;
  height: 36px;
  padding: 4px;
  border-radius: 4px;
  background: ${(props) => props.theme.colors.neutralLighter};
  align-items: center;
  justify-content: center;
`

const Badge = styled(Box)`
  color: ${(props) => props.theme.colors.neutralDark};
  background: ${(props) => props.theme.colors.neutralLighter};
  padding: 4px 8px;
  border-radius: 2px;
  width: 48px;
  justify-content: center;
`
const SelectionIndicator = styled(Box).attrs({})`
  border: 1px solid ${(props) => props.theme.colors.neutralLight};
  width: 10px;
  height: 10px;
  border-radius: 50%;
  ${(props) =>
    props.isSelected &&
    css`
      background: ${props.theme.colors.neutralDarker};
    `};
`

const NetworkCard = ({ id, name, addressRange, numHosts, isSelected, onClick }) => {
  const [isExpanded, setExpanded] = useState(false)

  const handleClick = () => {
    setExpanded(!isExpanded)
  }

  const handleLeaveButtonClick = (e) => {
    e.preventDefault()
    e.stopPropagation()
  }

  return (
    <Container>
      <HeaderContainer onClick={handleClick}>
        <StyledIcon mr="12px" icon={<icons.Network />} color="neutralDarker" />
        <Box flexDirection="column" width="240px">
          <Text textStyle="subtitle" fontSize="14px">
            {name}
          </Text>
          <Text textStyle="detail" fontSize="12px">
            {id}
          </Text>
        </Box>
        <Box width="200px" justifyContent="center">
          <Text textStyle="detail" fontSize="12px" title="Managed address">
            {addressRange}
          </Text>
        </Box>

        <Box width="200px" justifyContent="center">
          <Badge>
            <Text textStyle="detail">
              {numHosts} {pluralize('host', numHosts)}
            </Text>
          </Badge>
        </Box>
        <IconButton
          icon={<icons.Leave />}
          color="neutralDark"
          size="16px"
          ml="auto"
          onClick={handleLeaveButtonClick}
        />
      </HeaderContainer>
      {isExpanded && (
        <>
          <Separator mx={'10px'} width="auto" />
          <HiddenContentContainer>
            <Box alignItems="center" mb={3}>
              <Text textStyle="subtitle" fontSize="16px">
                Interface Configuration
              </Text>
              <Button variant="primary" ml="auto" height="32px">
                Save
              </Button>
            </Box>

            <Box alignItems="center" mb={2}>
              <Text textStyle="body" mr={2}>
                Address
              </Text>
              <TextInput placeholder="192.168.0.1/24" height="40px" width="180px" />
            </Box>

            <Box alignItems="center">
              <Text textStyle="body" mr={2}>
                Listen Port
              </Text>
              <NumberInput placeholder="51820" height="40px" width="180px" />
            </Box>
          </HiddenContentContainer>
        </>
      )}
    </Container>
  )
}

NetworkCard.propTypes = {
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  addressRange: PropTypes.string.isRequired,
  isSelected: PropTypes.bool,
  numHosts: PropTypes.number.isRequired,
  onClick: PropTypes.func,
}

NetworkCard.defaultProps = {
  onClick: () => {},
  isSelected: false,
}

export default NetworkCard
