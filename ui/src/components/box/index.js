import styled from 'styled-components'
import { border, borderStyle, color, flexbox, grid, layout, shadow, space } from 'styled-system'
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
}

export default Box
