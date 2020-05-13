import React from 'react'
import styled from 'styled-components'
import { grid, space, color } from 'styled-system'

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

const StyledLink = styled.a`
  font-family: Lato;
  font-size: 14px;
  text-decoration: none;
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
    <StyledLink href="https://www.github.com/seashell/drago/issues">Support</StyledLink>
    <StyledLink href="https://www.github.com/seashell/drago">Docs</StyledLink>
  </Container>
)

Footer.defaultProps = {
  padding: 2,
}

export default Footer
