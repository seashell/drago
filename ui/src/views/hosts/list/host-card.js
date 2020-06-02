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
  button {
    margin-right: auto;
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

const HostCard = ({ id, label, address, advertiseAddress, onClick, onDelete }) => (
  <Container onClick={() => onClick(id)}>
    <Box alignItems="center" width="240px">
      <IconContainer mr="12px">
        <IconButton ml="auto" icon={<icons.Host />} />
      </IconContainer>
      <Text textStyle="subtitle" fontSize="14px">
        {label}
      </Text>
      <Text textStyle="detail" fontSize="12px">
        {address}
      </Text>
    </Box>
    {advertiseAddress && <Badge>{advertiseAddress}</Badge>}
    <IconButton icon={<icons.Times />} onClick={onDelete} />
  </Container>
)

HostCard.propTypes = {
  id: PropTypes.string.isRequired,
  label: PropTypes.string.isRequired,
  address: PropTypes.string.isRequired,
  advertiseAddress: PropTypes.string,
  onClick: PropTypes.func,
  onDelete: PropTypes.func,
}

HostCard.defaultProps = {
  advertiseAddress: undefined,
  onClick: () => {},
  onDelete: () => {},
}

export default HostCard
