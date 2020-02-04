import React from 'react'
import styled from 'styled-components'
import { navigate } from '@reach/router'

import { useQuery } from 'react-apollo'

import { GET_NODES } from '_graphql/actions'
import Box from '_components/box'
import Text from '_components/text'

const Container = styled(Box).attrs({
  display: 'flex',
  border: 'discrete',
  m: 3,
  p: 3,
})`
  align-items: center;
  cursor: pointer;
  :hover {
    box-shadow: ${props => props.theme.shadows.medium};
  }
`

const IconContainer = styled(Box).attrs({
  display: 'block',
  height: '48px',
  width: '48px',
  bg: 'neutralLighter',
  borderRadius: '4px',
})`
  position: relative;
`

const StatusBadge = styled.div`
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: 4px solid white;
  position: absolute;
  right: -2px;
  bottom: -2px;
  background: ${props => (props.status === 'online' ? 'green' : 'red')};
`

// eslint-disable-next-line react/prop-types
const NodeCard = ({ id, label, onClick, address }) => (
  <Container onClick={() => onClick(id)}>
    <IconContainer mr="12px">
      <StatusBadge status="online" />
    </IconContainer>
    <Box flexDirection="column">
      <Text textStyle="subtitle" fontSize="14px">
        {label}
      </Text>
      <Text textStyle="detail" fontSize="12px">
        {address}
      </Text>
    </Box>
  </Container>
)

const onNodeCardClick = id => {
  navigate(`/nodes/${id}`)
}

const NodesView = () => {
  const { loading, error, data } = useQuery(GET_NODES)
  return (
    <>
      <Text textStyle="title">Nodes</Text>
      {loading
        ? 'loading'
        : data.result.items.map(n => (
            <NodeCard id={n.id} label={n.label} address={n.address} onClick={onNodeCardClick} />
          ))}
    </>
  )
}

export default NodesView
