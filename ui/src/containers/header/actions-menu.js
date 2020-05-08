import React from 'react'
import styled from 'styled-components'
import { color } from 'styled-system'

import Separator from '_components/separator'
import List from '_components/list'
import Link from '_components/link'

const MenuItem = styled.li`
  width: 120px;
  min-height: 36px;

  box-sizing: border-box;
  padding: 8px 16px;

  display: flex;

  align-items: center;
  justify-content: flex-start;

  cursor: pointer;
  font-size: 16px;
  color: ${props => props.theme.colors.neutralDarker};

  > div {
    display: flex;
    justify-items: flex-start;
    align-items: flex-start;
    flex-direction: column;
    margin-left: 8px;
    line-height: 14px;

    > :last-child {
      font-size: 14px;
      color: ${props => props.theme.colors.neutralDark};
    }
  }
  ${color}
`

const StyledSeparator = styled(Separator).attrs({
  bg: 'neutralLighter',
  my: 1,
})``

const ActionsMenu = () => (
  <List width="max-contents" padding="8px 0">
    <MenuItem>
      <Link to="/hosts/new">Add node</Link>
    </MenuItem>
    <StyledSeparator />
    <MenuItem>
      <Link to="/topology">Graph</Link>
    </MenuItem>
  </List>
)

export default ActionsMenu
