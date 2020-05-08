import styled from 'styled-components'
import { grid, colorStyle, color, space, layout } from 'styled-system'

import Link from '_components/link'
import Collapse from '_components/collapse'
import Separator from '_components/separator'

export const Container = styled.div`
${grid}
${space}
${layout}
${colorStyle}

position: fixed;
width: 200px;

display: flex;
flex-direction: column;

top:0;
left:0;
bottom:0;
`

export const StyledSeparator = styled(Separator).attrs({
  bg: 'neutralDarkest',
  my: 3,
})``

export const NavLink = styled(Link).attrs({
  px: 3,
  py: 2,
  hoverStyle: { background: 'rgba(255,255,255,0.05)' },
  activeStyle: { background: 'rgba(255,255,255,0.1)' },
})`
  display: flex;
  ${color}
  ${space}
  ${layout}
`

export const ActionNavLink = styled(Link).attrs({
  px: 3,
  py: 2,
  color: 'secondary',
  hoverStyle: { background: 'none' },
  activeStyle: { background: 'none' },
})`
  display: flex;

  svg {
    pading: 16px;
    width: 20px;
    height: 20px;
    border-radius: 1px;
    margin-right: 8px;
    fill: ${props => props.theme.colors.secondary};
  }

  :hover {
    svg {
      background: ${props => props.theme.colors.secondary};
      fill: ${props => props.theme.colors.white};
      transition: 0.5s ease-in;
    }
  }

  ${color}
  ${space}
  ${layout}
`

export const CollapsibleSection = styled(Collapse).attrs({
  px: 3,
  py: 2,
})`
  opacity: 0.5;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.1em;

  .indicator {
    fill: ${props => props.theme.colors.neutralDarker};
  }
`
