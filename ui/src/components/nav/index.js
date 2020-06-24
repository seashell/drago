import styled from 'styled-components'

import Box from '_components/box'
import Link from '_components/link'

export const Nav = styled(Box).attrs({
  width: '100%',
  borderBottom: 'medium',
  borderColor: 'neutralLighter',
})``

export const HorizontalNavLink = styled(Link).attrs(({ theme }) => ({
  py: 3,
  px: 2,
  color: 'neutral',
  fontSize: '14px',
  fontWeight: '500',
  marginBottom: '-2px',
  activeStyle: {
    color: theme.colors.primaryDark,
    'border-bottom': `2px solid ${theme.colors.primary}`,
  },
  hoverStyle: { 'border-bottom': `2px solid ${theme.colors.neutralDark}` },
}))`
  :first-child {
    margin-left: 0;
  }
`

export const VerticalNavLink = styled(Link).attrs(props => ({
  py: 2,
  color: 'neutral',
  activeStyle: { color: props.theme.colors.primary },
}))``
