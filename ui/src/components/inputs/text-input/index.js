import styled from 'styled-components'
import { space, layout } from 'styled-system'

const TextInput = styled.input`
  box-sizing: border-box;
  padding-bottom: 16px;
  padding-left: 13.6px;
  padding-right: 24.8px;
  padding-top: 16px;
  border: 1px solid ${props => props.theme.colors.neutralLighter};
  border-radius: 2px;
  font-size: 14px;
  width: 100%;
  :focus {
    border: 1px solid ${props => props.theme.colors.primary};
  }
  :disabled {
    background: inherit;
    opacity: 0.5;
    border: none;
  }
  ${space}
  ${layout}
`

export default TextInput
