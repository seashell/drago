import React from 'react'
import styled, { css } from 'styled-components'
import { color } from 'styled-system'

import { useQuery } from '@apollo/react-hooks'
import { GET_CURRENT_USER } from '_graphql/actions'

import Separator from '_components/separator'
import Avatar from '_components/avatar'
import List from '_components/list'
import Link from '_components/link'

const MenuItem = styled.li`
  width: 100%;
  min-height: 36px;

  box-sizing: border-box;
  padding: 8px 16px;

  display: flex;

  align-items: center;
  justify-content: flex-start;

  cursor: pointer;
  font-size: 16px;
  color: ${props => props.theme.colors.neutralDarker};

  ${props =>
    props.selected &&
    css`
      background: ${props.theme.colors.primaryLightest};
    `}

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

const UserMenu = () => {
  const { loading, data } = useQuery(GET_CURRENT_USER)

  if (loading) return null
  const { viewer } = data
  const { firstName, lastName, email, gravatar, organizations } = viewer

  return (
    <List width="max-contents" padding="8px 0">
      <MenuItem selected>
        <Avatar name={`${firstName} ${lastName}`} src={gravatar} size="36" style={{ margin: 0 }} />
        <div>
          <p>
            {firstName} {lastName}
          </p>
          <p>{email}</p>
        </div>
      </MenuItem>
      {organizations.nodes.map(org => (
        <MenuItem key={org.id}>
          <Avatar name={org.name} size="36" style={{ margin: 0 }} />
          <div>
            <p>{org.name}</p>
            <p />
          </div>
        </MenuItem>
      ))}
      <Separator bg="neutralLighter" my={1} />
      <MenuItem color="primary">
        <Link to="/organizations/new/">New organization</Link>
      </MenuItem>
      <MenuItem>
        <Link to="/account">Account</Link>
      </MenuItem>
    </List>
  )
}

export default UserMenu
