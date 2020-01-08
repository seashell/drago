import React from 'react'
import styled from 'styled-components'
import { grid, space, color } from 'styled-system'

import Link from '_components/link'

const Container = styled.div`
  height: 36px;

  background: white;
  border-top: 1px solid #f1f1f1;

  display: flex;
  justify-content: center;

  align-items: center;
  * + * {
    margin-left: 12px;
  }
  ${space}
  ${grid}
`

const StyledLink = styled(Link)`
  font-family: Lato;
  font-size: 14px;
  &:hover {
    color: ${props => props.theme.colors.primary};
  }
  ${color}
`

StyledLink.defaultProps = {
  color: 'neutralDarker',
}

const Footer = props => (
  <Container {...props}>
    <StyledLink to="/">Support</StyledLink>
    <StyledLink to="/">Docs</StyledLink>
  </Container>
)

Footer.defaultProps = {
  padding: 2,
}

export default Footer
