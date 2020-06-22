import React from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'

import Box from '_components/box'
import Text from '_components/text'
import IconButton from '_components/icon-button'
import { icons } from '_assets/'
import { space } from 'styled-system'

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
  ${space}
`

const HostCard = ({ id, name, labels, advertiseAddress, onClick, onDelete }) => (
  <Container onClick={() => onClick(id)}>
    <Box alignItems="center" width="240px">
      <IconContainer mr="12px">
        <IconButton ml="auto" icon={<icons.Host />} />
      </IconContainer>
      <div>
        <Text textStyle="subtitle" fontSize="14px">
          {name}
        </Text>
        <Text textStyle="detail">{advertiseAddress && ` @ ${advertiseAddress}`}</Text>
      </div>
    </Box>
    <Box>
      {labels.map(el => (
        <Badge ml={2}>{el}</Badge>
      ))}
    </Box>
    <IconButton icon={<icons.Times />} onClick={onDelete} />
  </Container>
)

HostCard.propTypes = {
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  labels: PropTypes.arrayOf(PropTypes.string),
  advertiseAddress: PropTypes.string,
  onClick: PropTypes.func,
  onDelete: PropTypes.func,
}

HostCard.defaultProps = {
  labels: [],
  advertiseAddress: undefined,
  onClick: () => {},
  onDelete: () => {},
}

export default HostCard
