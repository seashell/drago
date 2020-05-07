import styled, { css } from 'styled-components'
import { layout, space, color } from 'styled-system'
import PropTypes from 'prop-types'

const Separator = styled.span`
  ${props =>
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
  bg: 'neutral',
}

export default Separator
