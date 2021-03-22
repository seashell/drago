import PropTypes from 'prop-types'
import styled, { css } from 'styled-components'
import { color, layout, space } from 'styled-system'

const Separator = styled.span`
  ${(props) =>
    props.vertical
      ? css`
          width: 1px;
          height: 100%;
        `
      : css`
          height: 1px;
          width: 100%;
        `};
  display: block;
  ${layout}
  ${space}
  ${color}
`

Separator.propTypes = {
  vertical: PropTypes.bool,
}

Separator.defaultProps = {
  vertical: false,
  bg: 'neutralLighter',
}

export default Separator
