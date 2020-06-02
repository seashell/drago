import React from 'react'

import { Container, CollapsibleSection, NavLink } from './styled'

const SideNav = props => (
  <Container {...props}>
    <CollapsibleSection title="Manage" isOpen>
      <NavLink to="/networks">Networks</NavLink>
      <NavLink to="/hosts">Hosts</NavLink>
    </CollapsibleSection>
  </Container>
)

SideNav.defaultProps = {
  colors: 'darkHighContrast',
}

export default SideNav
