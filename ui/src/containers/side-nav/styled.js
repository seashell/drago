import styled from 'styled-components'
import { color, colorStyle, grid, layout, space } from 'styled-system'
import Box from '_components/box'
import Collapse from '_components/collapse'
import Link from '_components/link'
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
  padding-top: 100px;

  top: 0;
  left: 0;
  bottom: 0;

  border-right: 1px solid ${(props) => props.theme.colors.neutralLighter};
  background: transparent;

  z-index: 9;

  display: none;

  // Laptops and above
  @media (min-width: 1280px) {
    display: flex;
  }
`

export const StyledSeparator = styled(Separator).attrs({
  bg: 'border1',
  my: 3,
})``

export const NavButton = styled(Box)`
  * {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  display: flex;
  align-items: center;
  cursor: pointer;

  font-size: 16px;
  color: ${({ theme }) => theme.colors.primary};
  :hover {
    border-right: 3px solid !important;
    border-color: ${({ theme }) => theme.colors.primary}!important;
  }

  padding-top: 8px;
  padding-bottom: 8px;

  border-right: ${(props) => (props.isCurrent ? '3px solid' : 'none')};
  border-color: ${(props) => (props.isCurrent ? props.theme.colors.primary : '')};

  ${color}
  ${space}
  ${layout}
`

export const NavLink = styled(Link).attrs((props) => ({
  px: 3,

  getProps: ({ isPartiallyCurrent }) => ({
    style: {
      borderRight: isPartiallyCurrent ? '3px solid' : 'none',
      borderColor: isPartiallyCurrent ? props.theme.colors.primary : '',
    },
  }),
}))`
  * {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  display: flex;
  align-items: center;

  font-size: 16px;
  color: ${({ theme }) => theme.colors.primary};
  :hover {
    border-right: 3px solid !important;
    border-color: ${({ theme }) => theme.colors.primary}!important;
  }

  padding-top: 8px;
  padding-bottom: 8px;

  ${color}
  ${space}
  ${layout}
`

export const CollapsibleSection = styled(Collapse).attrs({
  px: 3,
  py: 2,
})`
  color: ${(props) => props.theme.colors.neutralLight};
  font-weight: 500;
  font-size: 0.76rem;
  letter-spacing: 0.06rem;
  text-transform: uppercase;
`
