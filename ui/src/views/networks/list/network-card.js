import PropTypes from 'prop-types'
import React from 'react'
import styled from 'styled-components'
import { icons } from '_assets/'
import Box from '_components/box'
import Icon from '_components/icon'
import IconButton from '_components/icon-button'
import Text from '_components/text'

const Container = styled(Box).attrs({
  px: 3,
  border: 'discrete',
})`
  height: 72px;
  cursor: pointer;
  align-items: center;
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
  justify-content: center;
`

const NetworkCard = ({
  id,
  name,
  addressRange,
  interfacesCount,
  connectionsCount,
  onClick,
  onDelete,
}) => {
  const handleDeleteButtonClick = (e) => {
    e.preventDefault()
    e.stopPropagation()
    onDelete()
  }
  return (
    <Container onClick={() => onClick(id)}>
      <Box>
        <StyledIcon ml="auto" mr="12px" icon={<icons.Network />} color="neutralDarker" />
        <Box flexDirection="column" width="200px">
          <Text textStyle="subtitle" fontSize="14px">
            {id.split('-')[0]}
          </Text>
          <Text textStyle="detail" fontSize="12px">
            {name}
          </Text>
        </Box>
      </Box>

      <Box width="48px" justifyContent="center">
        <Badge alignItems="center">
          <Icon icon={<icons.Host />} color="neutralDarker" size="14px" mr={'4px'} />
          <Text textStyle="detail" fontSize="10px" title="Managed address">
            {interfacesCount}
          </Text>
        </Badge>
      </Box>

      <Box width="48px" justifyContent="center">
        <Badge alignItems="center">
          <Icon icon={<icons.ConnectionSmall />} color="neutralDarker" size="14px" mr={'4px'} />
          <Text textStyle="detail" fontSize="10px" title="Managed address">
            {connectionsCount}
          </Text>
        </Badge>
      </Box>

      <Box width="200px" justifyContent="center">
        <Text textStyle="detail">{addressRange}</Text>
      </Box>

      <IconButton
        icon={<icons.Trash />}
        color="neutralDark"
        size="16px"
        ml="auto"
        onClick={handleDeleteButtonClick}
      />
    </Container>
  )
}

NetworkCard.propTypes = {
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  addressRange: PropTypes.string.isRequired,
  interfacesCount: PropTypes.number,
  connectionsCount: PropTypes.number,
  onClick: PropTypes.func,
  onDelete: PropTypes.func,
}

NetworkCard.defaultProps = {
  interfacesCount: 0,
  connectionsCount: 0,
  onClick: () => {},
  onDelete: () => {},
}

export default NetworkCard
