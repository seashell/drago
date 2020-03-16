import React from 'react'
import PropTypes from 'prop-types'

import styled from 'styled-components'
import Box from '_components/box'
import Text from '_components/text'

const Container = styled(Box)`
  background: #eee;
  height: 60px;
  width: 60px;
  border-radius: 50%;
  position: relative;
`

const RemoveButton = styled.button`
  position: absolute;
  height: 20px;
  width: 20px;
  border-radius: 50%;
  bottom: -4px;
  right: -4px;
  background: #444;
  border: none;
  color: white;
  opacity: 0.3;
  cursor: pointer;
  :hover {
    opacity: 1;
  }
`

const Peer = ({ id, name, onRemove, ...props }) => (
  <Container m={2}>
    <Text>{name}</Text>
    <RemoveButton onClick={() => onRemove(id)}>x</RemoveButton>
  </Container>
)

Peer.propTypes = {
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  onRemove: PropTypes.func.isRequired,
}

export default Peer
