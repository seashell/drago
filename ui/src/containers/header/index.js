import React from 'react'
import styled from 'styled-components'
import { grid, space, color, shadow, border } from 'styled-system'

import Button from '_components/button'
import Brand from '_containers/side-nav/brand'

export const Container = styled.div`
  display: flex;

  align-items: center;
  justify-content: space-between;  
  
  border-bottom: 1px solid ${props => props.theme.colors.neutralLighter};
  position: fixed;

  top:0;
  right:0;
  left: 0;

  z-index: 99;

  ${border}
  ${shadow}
  ${color}
  ${space}
  ${grid}
`

Container.defaultProps = {
  backgroundColor: 'white',
  boxShadow: 'light',
  border: 'dark',
}

export const StyledButton = styled(Button).attrs({
  variant: 'primary',
  height: '40px',
  width: '100px',
  borderRadius: 3,
  mr: 3,
})`
  display: flex;
  align-items: center;
  justify-content: center;
  padding-left: 12px;
  line-height: 8px;
`

const Header = props => (
  <Container {...props}>
    <Brand />
  </Container>
)

Header.defaultProps = {
  padding: 3,
}

export default Header