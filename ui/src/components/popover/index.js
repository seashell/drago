import React from 'react'
import styled from 'styled-components'
import Tippy from '@tippy.js/react'
import PropTypes from 'prop-types'

const StyledTippy = styled(Tippy).attrs({
  arrow: true,
  sticky: true,
  touchHold: true,
  interactive: true,
  boundary: 'viewport',
})`
  background: #fff;
  border: ${props => props.theme.borders.discrete};
  box-shadow: ${props => props.theme.shadows.medium};
  padding: 0;
  .tippy-arrow {
    :before {
      content: '';
      position: absolute;
      top: -8px;
      left: -8px;
      border-bottom: 8px solid transparent;
      border-left: 8px solid transparent;
      border-right: 8px solid transparent;
      border-top: 8px solid ${props => props.theme.colors.white};
    }
  }

  &[x-placement^='bottom'] {
    .tippy-arrow {
      border-bottom-color: ${props => props.theme.colors.neutralLighter};
    }
  }
`

const Popover = ({ trigger, duration, delay, children, content, ...props }) => (
  <StyledTippy content={content} {...props}>
    <div style={{ cursor: 'pointer' }}>{children}</div>
  </StyledTippy>
)

Popover.propTypes = {
  trigger: PropTypes.string,
  duration: PropTypes.arrayOf(PropTypes.number),
  delay: PropTypes.arrayOf(PropTypes.number),
  children: PropTypes.node.isRequired,
  content: PropTypes.node.isRequired,
}

Popover.defaultProps = {
  trigger: 'hover',
  duration: [500, 200],
  delay: [100, 10],
}

export default Popover
