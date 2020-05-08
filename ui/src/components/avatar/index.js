import React from 'react'
import styled from 'styled-components'
import { layout, space } from 'styled-system'

import ReactAvatar, { ConfigProvider } from 'react-avatar'

const StyledReactAvatar = styled(ReactAvatar)`
  ${layout}
  ${space}
`

const Avatar = props => (
  <ConfigProvider>
    <StyledReactAvatar size="40" maxInitials={2} round {...props} />
  </ConfigProvider>
)

export default Avatar
