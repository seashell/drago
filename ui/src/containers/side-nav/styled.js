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
width: 240px;

display: flex;
flex-direction: column;
padding-top:110px;

top:0;
left:0;
bottom:0;

background: transparent;

z-index: 99;
`

export const StyledSeparator = styled(Separator).attrs({
  bg: 'neutralDarkest',
  my: 3,
})``

export const NavLink = styled(Link).attrs(props => ({
  px: 3,
  py: 2,
  getProps: ({ isPartiallyCurrent }) => ({
    style: { color: isPartiallyCurrent && props.theme.colors.primary },
  }),
}))`
  display: flex;
  align-items: center;
  :hover {
    color: ${({ theme }) => theme.colors.neutralDark};
  };
  ${color}
  ${space}
  ${layout}
`

export const CollapsibleSection = styled(Collapse).attrs({
  px: 3,
  py: 2,
})`
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.1em;
`
