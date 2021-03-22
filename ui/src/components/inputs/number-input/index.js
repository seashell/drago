import styled from 'styled-components'
import { layout, space } from 'styled-system'

const NumberInput = styled.input.attrs({
  type: 'number',
})`
  box-sizing: border-box;

  border: none;
  border-bottom: 1px solid ${(props) => props.theme.colors.neutralLighter};

  padding-bottom: 4px;

  font-family: Lato;
  font-size: 16px;

  width: 100%;

  :focus {
    border-bottom: 1px solid ${(props) => props.theme.colors.primary};
  }
  :disabled {
    background: inherit;
    opacity: 0.5;
    border: none;
  }
  :invalid {
    border-bottom: 1px solid ${(props) => props.theme.colors.danger};
  }
  ::placeholder {
    color: ${(props) => props.theme.colors.neutralLight};
  }

  ::-webkit-outer-spin-button,
  ::-webkit-inner-spin-button {
    -webkit-appearance: none;
    margin: 0;
  }

  input[type='number'] {
    -moz-appearance: textfield; /* Firefox */
  }

  ${space}
  ${layout}
`

export default NumberInput
