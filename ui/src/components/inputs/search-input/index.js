/* eslint-disable react/prop-types */
import React from 'react'
import styled from 'styled-components'
import { layout, space, border, color } from 'styled-system'

import { icons } from '_assets/'

const StyledIcon = styled(icons.Search).attrs({
  width: 22,
})`
  padding: 8px;
  fill: ${({ theme: { colors } }) => colors.neutralDark};
`

const StyledInput = styled.input.attrs({
  type: 'text',
})`
  width: 100%;
  height: 100%;

  font-size: 15px;

  background: none;
  border: none;

  padding: 0;

  color: ${props => props.theme.colors.neutralDarker};

  ::placeholder {
    color: ${props => props.theme.colors.neutralDarker};
  }
`

const Container = styled.div`
  display: flex;
  flex-direction: row-reverse;
  align-items: center;
  height: 100%;
  width: 60%;
  ${layout}
  ${space}
  ${border}
  ${color}

  ${StyledInput}:focus + ${StyledIcon}{
    fill: ${({ theme: { colors } }) => colors.primary};
  }
`

const SearchInput = ({ ...props }) => (
  <Container {...props}>
    <StyledInput ref={props.innerRef} {...props} />
    <StyledIcon />
  </Container>
)

export default React.forwardRef((props, ref) => <SearchInput innerRef={ref} {...props} />)
