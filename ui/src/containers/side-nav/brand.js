import React from 'react'
import styled from 'styled-components'
import { colorStyle } from 'styled-system'
import { Link } from '@reach/router'
import { ReactComponent as Logo } from '_assets/icons/logo.svg'

const StyledLink = styled(Link)`
  position: relative;

  height: auto;

  display: flex;
  align-items: center;

  svg {
    fill: ${props => props.theme.colors.primary};
    opacity: 0.3;
  }
  &:hover {
    svg {
      opacity: 1;
      transition: all 0.4s linear;
    }
  }
  ${colorStyle}
`

const Brand = props => (
  <StyledLink to="">
    <Logo width={56} height={56} {...props} />
  </StyledLink>
)

export default Brand
