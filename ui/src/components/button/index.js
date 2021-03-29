import styled from 'styled-components'

import { buttonStyle, shadow, layout, typography, border, space, color } from 'styled-system'

const Button = styled.button`
  font-family: 'Lato';
  font-weight: bold;
  letter-spacing: 0.04em;
  border: none;

  :disabled {
    opacity: 0.4;
    background: #ddd;
    border-color: #ddd;
    color: #666;
    box-shadow: none;
    cursor: not-allowed;
  }

  &:hover {
    filter: brightness(95%);
    transition: all 0.7s ease;
  }

  border-radius: 1px;

  ${buttonStyle}
  ${typography}
  ${shadow}
  ${layout}
  ${space}
  ${border}
  ${color}
`

Button.defaultProps = {
  variant: 'primary',
  width: 'max-content',
  px: '24px',
  height: '50px',
  type: 'button',
}

export default Button
