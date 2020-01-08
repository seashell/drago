import React from 'react'
import styled from 'styled-components'
import { colorStyle } from 'styled-system'
import { Link } from '@reach/router'
import { ReactComponent as Logo } from '_assets/icons/logo.svg'

const StyledLink = styled(Link)`
  position: relative;
  padding-left: 14px;

  height: 76px;

  display: flex;
  align-items: center;

  svg {
    fill: ${props => props.theme.colors.neutralLightest};
    opacity: 0.3;
  }
  &:hover {
    svg {
      opacity: 0.8;
      transition: all 0.4s linear;
    }
  }
  ${colorStyle}
`

const Brand = props => (
  <StyledLink to="/">
    <Logo width={48} height={48} {...props} />
  </StyledLink>
)

export default Brand
