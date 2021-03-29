import React from 'react'
import { CollapsibleSection, Container, NavLink } from './styled'

const SideNav = (props) => (
  <Container {...props}>
    <CollapsibleSection title="Manage" isOpen>
      <NavLink to="networks">Networks</NavLink>
      <NavLink to="clients">Clients</NavLink>
    </CollapsibleSection>
  </Container>
)

SideNav.defaultProps = {
  colors: 'darkHighContrast',
}

export default SideNav
