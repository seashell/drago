import PropTypes from 'prop-types'
import React from 'react'
import styled from 'styled-components'
import { icons } from '_assets/'
import Box from '_components/box'
import Icon from '_components/icon'
import IconButton from '_components/icon-button'
import Text from '_components/text'

const Container = styled(Box).attrs({
  border: 'discrete',
})`
  flex-direction: column;
  :not(:last-child) {
    border-bottom: none;
  }
`

const HeaderContainer = styled(Box).attrs({
  px: 3,
})`
  display: grid;
  grid-template-columns: 320px auto 160px 120px auto;

  height: 72px;
  cursor: pointer;
  align-items: center;
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

const StyledLink = styled.a`
  text-decoration: none;
`

const PeerCard = ({
  id,
  nodeName,
  nodeStatus,
  nodeAdvertiseAddress,
  name,
  address,
  hasPublicKey,
  onClick,
  onDelete,
}) => {
  const handleCardClick = (e) => {
    e.preventDefault()
    e.stopPropagation()
    onClick()
  }

  const handleDeleteButtonClick = (e) => {
    e.preventDefault()
    e.stopPropagation()
    onDelete()
  }

  return (
    <Container onClick={handleCardClick}>
      <HeaderContainer>
        <Box>
          <IconContainer mr="12px">
            <StyledIcon icon={<icons.Host />} color="neutralDarker" />
            <StatusIndicator color={nodeStatus === 'ready' ? 'success' : 'danger'} />
          </IconContainer>
          <Box flexDirection="column" width="240px">
            <Text textStyle="subtitle" fontSize="14px">
              {id.split('-')[0]}
            </Text>
            <Text textStyle="detail" fontSize="12px">
              {nodeName} (via {name ? `${name}/${id.split('-')[0]}` : id.split('-')[0]})
            </Text>
          </Box>
        </Box>

        {hasPublicKey ? (
          <Badge alignItems="center" justifySelf="center">
            <Icon icon={<icons.Key />} color="neutralDarker" size="14px" mr={'4px'} />
          </Badge>
        ) : (
          <span />
        )}

        <Box width="200px" justifyContent="center">
          <Text textStyle="detail" fontSize="12px" title="Managed address">
            {address}
          </Text>
        </Box>
        <Box width="200px" justifyContent="center">
          <StyledLink href={`http://${nodeAdvertiseAddress}`}>
            <Text textStyle="detail" fontWeight="600" fontSize="12px" title="Public address">
              {nodeAdvertiseAddress}
            </Text>
          </StyledLink>
        </Box>
        <IconButton
          icon={<icons.Trash />}
          color="neutralDark"
          size="16px"
          ml="auto"
          onClick={handleDeleteButtonClick}
        />
      </HeaderContainer>
    </Container>
  )
}

PeerCard.propTypes = {
  id: PropTypes.string.isRequired,
  name: PropTypes.string,
  address: PropTypes.string,
  nodeName: PropTypes.string,
  nodeStatus: PropTypes.string,
  nodeAdvertiseAddress: PropTypes.string,
  hasPublicKey: PropTypes.bool,
  onClick: PropTypes.func,
  onDelete: PropTypes.func,
}

PeerCard.defaultProps = {
  name: '',
  address: undefined,
  nodeName: '',
  nodeStatus: 'down',
  nodeAdvertiseAddress: undefined,
  hasPublicKey: false,
  onClick: () => {},
  onDelete: () => {},
}

export default PeerCard
