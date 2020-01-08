import React, { useRef } from 'react'
import { navigate } from '@reach/router'
import { useHotkeys } from 'react-hotkeys-hook'
import styled from 'styled-components'
import { grid, space, color, shadow, border } from 'styled-system'

import { useQuery } from 'react-apollo'
import { GET_CURRENT_USER } from '_graphql/actions'

import Popover from '_components/popover'
import Button from '_components/button'

import { icons } from '_assets/'

import ActionsMenu from './actions-menu'

export const Container = styled.div`
  display: flex;

  align-items: center;
  justify-content: flex-end;  
  
  border-bottom: 1px solid ${props => props.theme.colors.neutralLighter};
  position: fixed;

  top:0;
  right:0;
  left: 0;

  z-index: 99;

  ${border}
  ${shadow}
  ${color}
  ${space}
  ${grid}
`

Container.defaultProps = {
  backgroundColor: 'white',
  boxShadow: 'light',
  border: 'dark',
}

export const StyledButton = styled(Button).attrs({
  variant: 'secondary',
  height: '40px',
  width: '100px',
  borderRadius: 3,
  mr: 3,
})`
  display: flex;
  align-items: center;
  justify-content: center;
  padding-left: 12px;
  line-height: 8px;
`

const handleOnSupportIconButtonClick = e => {
  e.preventDefault()
  e.stopPropagation()
  window.open('https://www.google.com', '_blank')
}

const handleOnNotificationsIconButtonClick = e => {
  e.preventDefault()
  e.stopPropagation()
  navigate('/notifications')
}

const Header = props => {
  const { loading, data } = useQuery(GET_CURRENT_USER)

  if (loading) return <Container />

  return (
    <Container {...props}>
      <Popover content={<ActionsMenu />}>
        <StyledButton>
          Create
          <icons.ArrowDown fill="white" width={16} height={16} />
        </StyledButton>
      </Popover>
    </Container>
  )
}

Header.defaultProps = {
  padding: 3,
}

export default Header
