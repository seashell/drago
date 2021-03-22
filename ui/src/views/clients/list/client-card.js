import PropTypes from 'prop-types'
import React from 'react'
import styled from 'styled-components'
import { icons } from '_assets/'
import Box from '_components/box'
import Icon from '_components/icon'
import Text from '_components/text'

const Container = styled(Box).attrs({
  display: 'flex',
  border: 'discrete',
  p: 3,
})`
  cursor: pointer;
  align-items: center;
  :not(:last-child) {
    border-bottom: none;
  }
`

const IconContainer = styled(Box)`
  position: relative;
`

const StyledIcon = styled(Icon)`
  width: 36px;
  height: 36px;
  padding: 4px;
  border-radius: 4px;
  background: ${(props) => props.theme.colors.neutralLighter};
  align-items: center;
  justify-content: center;
  position: relative;
  z-index: 2;
`
const StatusIndicator = styled(Box)`
  position: absolute;
  display: block;
  right: -4px;
  bottom: -4px;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  border: 3px solid #fff;
  background-color: ${(props) => props.theme.colors[props.color]};
  z-index: 4;
`
const Badge = styled(Box)`
  color: ${(props) => props.theme.colors.neutralDark};
  background: ${(props) => props.theme.colors.neutralLighter};
  padding: 4px 8px;
  border-radius: 2px;
  justify-content: center;
`

const ClientCard = ({ id, name, status, address, interfaceCount, connectionCount, onClick }) => (
  <Container onClick={() => onClick(id)}>
    <IconContainer mr="12px">
      <StyledIcon icon={<icons.Host />} color="neutralDarker" />
      <StatusIndicator color={status === 'ready' ? 'success' : 'danger'} />
    </IconContainer>
    <Box flexDirection="column" width="240px">
      <Text textStyle="subtitle" fontSize="14px">
        {id.split('-')[0]}
      </Text>
      <Text textStyle="detail">{name}</Text>
    </Box>

    <Box width="48px" justifyContent="center">
      <Badge alignItems="center">
        <Icon icon={<icons.Interface />} color="neutralDarker" size="14px" mr={'4px'} />
        <Text textStyle="detail" fontSize="10px" title="Managed address">
          {interfaceCount}
        </Text>
      </Badge>
    </Box>

    <Box width="48px" justifyContent="center">
      <Badge alignItems="center">
        <Icon icon={<icons.ConnectionSmall />} color="neutralDarker" size="14px" mr={'4px'} />
        <Text textStyle="detail" fontSize="10px" title="Managed address">
          {connectionCount}
        </Text>
      </Badge>
    </Box>

    <Box flexDirection="column" width="240px">
      <Text textStyle="detail">{address && ` @ ${address}`}</Text>
    </Box>
  </Container>
)

ClientCard.propTypes = {
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  status: PropTypes.string.isRequired,
  interfaceCount: PropTypes.number,
  connectionCount: PropTypes.number,
  address: PropTypes.string,
  onClick: PropTypes.func,
}

ClientCard.defaultProps = {
  address: undefined,
  interfaceCount: 0,
  connectionCount: 0,
  onClick: () => {},
}

export default ClientCard
