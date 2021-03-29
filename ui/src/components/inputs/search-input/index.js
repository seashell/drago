import React from 'react'
import styled from 'styled-components'
import { border, color, layout, space } from 'styled-system'
import { icons } from '_assets/'

const StyledIcon = styled(icons.Search).attrs({
  width: 22,
})`
  fill: ${({ theme: { colors } }) => colors.neutral};
  margin: auto;
  position: absolute;
  left: 8px;
  top: 50%;
  transform: translateY(-50%);
  z-index: 1;
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
  border: 1px solid ${(props) => props.theme.colors.neutralLighter};
  border-radius: 24px;
  font-size: 14px;
  width: 100%;
  height: 100%;
  background: ${(props) => props.theme.colors.background1};
  color: ${(props) => props.theme.colors.foreground1};

  :focus {
    border: 1px solid ${(props) => props.theme.colors.primary};
  }
  :disabled {
    background: inherit;
    opacity: 0.5;
    border: none;
  }
  :invalid {
    border: 1px solid ${(props) => props.theme.colors.danger};
  }
  ${space}
  ${layout}
`

// eslint-disable-next-line react/prop-types
const SearchInput = ({ placeholder, onChange, ...props }) => {
  const handleChange = (e) => {
    onChange(e.target.value)
  }

  return (
    <Container {...props}>
      <StyledIcon />
      <StyledInput onChange={handleChange} placeholder={placeholder} />
    </Container>
  )
}

export default SearchInput
