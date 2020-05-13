import styled from 'styled-components'

import { buttonStyle, shadow, layout, typography, border, space } from 'styled-system'

const Button = styled.button`
  height: 48px;
  border-radius: 24px;
  
  font-family: 'Lato';
  font-weight: bold;
  letter-spacing: 0.08em;
  border: none;

  :disabled {
    background: ${props => props.theme.colors.neutralLighter};
    box-shadow: none;
    cursor: default;
    :hover {
      filter: none;
    }
  }

  &:hover {
    filter: brightness(90%);
    transition: all 0.7s ease;
  }
  ${buttonStyle}
  ${typography}
  ${shadow}
  ${layout}
  ${space}
  ${border}
`

Button.defaultProps = {
  variant: 'primary',
  boxShadow: 'light',
  width: '120px',
  height: '48px',
}

export default Button
