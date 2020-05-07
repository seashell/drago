import React from 'react'

import { useQuery } from 'react-apollo'
import { GET_PROJECTS } from '_graphql/actions'

import Avatar from '_components/avatar'

import { icons } from '_assets/'

import { Container, CollapsibleSection, NavLink, ActionNavLink, StyledSeparator } from './styled'

import Brand from './brand'

const SideNav = props => {
  const { loading, data } = useQuery(GET_PROJECTS)

  if (loading || !data) return <Container />

  const {
    viewer: {
      projects: { nodes },
    },
  } = data

  return (
    <Container {...props}>
      <Brand />
      <CollapsibleSection title="Projects" isOpen>
        {nodes.map(p => (
          <NavLink to={`/projects/${p.id}`}>
            <Avatar round="2px" mr={2} name={p.name} size={20} textSizeRatio={2} maxInitials={1} />
            {p.name}
          </NavLink>
        ))}
        <ActionNavLink to="/projects/new">
          <icons.Plus />
          New project
        </ActionNavLink>
      </CollapsibleSection>

      <StyledSeparator />

      <CollapsibleSection title="Manage" isOpen>
        <NavLink to="/foo">Foo</NavLink>
        <NavLink to="/bar">Bar</NavLink>
      </CollapsibleSection>

      <StyledSeparator />

      <CollapsibleSection title="Discover" isOpen>
        <NavLink to="/foobar">Foobar</NavLink>
      </CollapsibleSection>

      <StyledSeparator />

      <CollapsibleSection title="Account" isOpen>
        <NavLink to="/account/security">Security</NavLink>
        <NavLink to="/account/referrals">Referrals</NavLink>
      </CollapsibleSection>
    </Container>
  )
}

SideNav.defaultProps = {
  colors: 'darkHighContrast',
}

export default SideNav
