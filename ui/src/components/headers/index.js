import React from 'react'
import styled from 'styled-components'

import Avatar from '_components/avatar'
import Text from '_components/text'
import Box from '_components/box'

const StyledBox = styled(Box).attrs({
  height: '120px',
  display: 'flex',
  alignItems: 'center',
  padding: 0,
  mb: 3,
})`
  > :first-child {
    margin-right: 16px;
  }
`

const ProjectHeader = props => (
  <StyledBox {...props}>
    <Avatar name="some-project" size={64} round="4px" textSizeRatio={2} />
    <div>
      <Text textStyle="title">some-project</Text>
      <Text textStyle="subtitle">This projects does something really well</Text>
    </div>
  </StyledBox>
)

const NodeHeader = props => (
  <StyledBox {...props}>
    <Text textStyle="title">xyZ761kgVK</Text>
    <Text textStyle="subtitle">This projects does something really well</Text>
  </StyledBox>
)

export { ProjectHeader, NodeHeader }
