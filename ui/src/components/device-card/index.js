import React from 'react'
import PropTypes from 'prop-types'

import styled from 'styled-components'

import Box from '_components/box'
import Text from '_components/text'
import Button from '_components/button'
import List from '_components/list'
import Moment from 'react-moment'

const Container = styled(Box)`
  display: flex;
  width: 480px;
  flex-direction: column;
  border: 1px solid ${props => props.theme.colors.neutralLighter};
`

export const StyledButton = styled(Button).attrs({
  height: '30px',
  width: '80px',
  borderRadius: 3,
  mr: 3,
})`
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 8px;
`

const DeviceCard = props => (
  <Container border="dark" boxShadow="light" p={3} m={2} {...props}>
    <Text textStyle="title">{props.name}</Text>
    <Text textStyle="detail" my={3}>
      Created at <Moment>{props.createdAt}</Moment>
    </Text>
    <Text textStyle="detail" mb={4}>
      Last handshake at <Moment>{props.handshakeAt}</Moment>
    </Text>
    <List display="flex" alignItems="center" justifyContent="flex-end">
      <StyledButton onClick={props.onConnect}>Connect</StyledButton>
      <StyledButton variant="warning" onClick={props.onDisable}>
        Disable
      </StyledButton>
      <StyledButton variant="danger" onClidk={props.onDelete}>
        Delete
      </StyledButton>
    </List>
  </Container>
)

DeviceCard.propTypes = {
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  createdAt: PropTypes.string.isRequired,
  handshakeAt: PropTypes.string,
  onConnect: PropTypes.func.isRequired,
  onDisable: PropTypes.func.isRequired,
  onDelete: PropTypes.func.isRequired,
}

DeviceCard.defaultProps = {
  handshakeAt: '',
}

export default DeviceCard
