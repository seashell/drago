import React from 'react'
import styled from 'styled-components'
import { layout, space, border, color } from 'styled-system'

import { icons } from '_assets/'

const StyledIcon = styled(icons.Search).attrs({
  width: 22,
})`
  fill: ${({ theme: { colors } }) => colors.neutralDark};
  margin: auto;
  position: absolute;
  left: 8px;
  top: 50%;
  transform: translateY(-50%);
  z-index: 99;
`

const Container = styled.div`
  position: relative;
  height: 40px;
  ${layout}
  ${space}
  ${border}
  ${color}
`

const StyledInput = styled.input`
  box-sizing: border-box;
  padding-left: 32px;
  border: 1px solid ${props => props.theme.colors.neutralLighter};
  border-radius: 2px;
  font-size: 14px;
  width: 100%;
  height: 100%;
  :focus {
    border: 1px solid ${props => props.theme.colors.primary};
  }
  :disabled {
    background: inherit;
    opacity: 0.5;
    border: none;
  }
  :invalid {
    border: 1px solid ${props => props.theme.colors.danger};
  }
  ${space}
  ${layout}
`

// eslint-disable-next-line react/prop-types
const SearchInput = ({ placeholder, onChange, ...props }) => (
  <Container {...props}>
    <StyledIcon />
    <StyledInput onChange={onChange} placeholder={placeholder} />
  </Container>
)

export default SearchInput
