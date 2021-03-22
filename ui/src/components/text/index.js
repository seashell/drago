import styled from 'styled-components'
import { color, fontSize, fontWeight, layout, space, textAlign, textStyle } from 'styled-system'

const Text = styled.div`
  strong {
    font-weight: bold;
  }
  ${textStyle}
  ${textAlign}
  ${space}
  ${layout}
  ${fontSize}
  ${fontWeight}
  ${color}
`

export default Text
