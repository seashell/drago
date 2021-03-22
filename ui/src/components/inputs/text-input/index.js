import styled from 'styled-components'
import { border, layout, space } from 'styled-system'

const TextInput = styled.input.attrs({
  type: 'text',
  autocomplete: 'off',
})`
  box-sizing: border-box;

  border: none;
  border-bottom: 1px solid ${(props) => props.theme.colors.neutralLighter};

  padding-bottom: 4px;

  background-color: ${(props) => props.theme.colors.white};
  color: ${(props) => props.theme.colors.neutralDarker};

  font-family: Lato;
  font-size: 16px;

  width: 100%;
  ::placeholder {
    color: ${(props) => props.theme.colors.neutralLight};
  }
  :focus {
    border-bottom: 1px solid ${(props) => props.theme.colors.primary};
  }
  :disabled {
    opacity: 0.5;
    cursor: not-allowed;
    border: none;
  }
  :invalid {
    border-bottom: 1px solid ${(props) => props.theme.colors.danger};
  }

  ${border}
  ${space}
  ${layout}
`

export default TextInput
