import styled from 'styled-components'
import { grid, space, color, layout, flexbox, border, shadow, borderStyle } from 'styled-system'
import { containers } from '../../styles'

const Box = styled.div`
  ${grid}
  ${space}
  ${color}
  ${layout}
  ${flexbox}
  ${border}
  ${containers}
  ${borderStyle}
  ${shadow}
`

Box.defaultProps = {
  border: 'none',
  height: 'auto',
  display: 'flex',
  gridColumn: 'span 12',
}

export default Box
