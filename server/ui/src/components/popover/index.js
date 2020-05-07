import React from 'react'
import styled from 'styled-components'
import Tippy from '@tippy.js/react'
import PropTypes from 'prop-types'

const StyledTippy = styled(Tippy).attrs({
  arrow: true,
  sticky: true,
  touchHold: true,
  interactive: true,
  trigger: 'click',
  boundary: 'viewport',
  duration: [500, 200],
  delay: [200, 10],
})`
  background: #fff;
  border: ${props => props.theme.borders.discrete};
  box-shadow: ${props => props.theme.shadows.medium};
  padding: 0;
  .tippy-arrow {
    :before {
      content: '';
      position: absolute;
      top: -6px;
      left: -8px;
      border-top: 8px solid transparent;
      border-left: 8px solid transparent;
      border-right: 8px solid transparent;
      border-bottom: 8px solid ${props => props.theme.colors.white};
    }
  }

  &[x-placement^='bottom'] {
    .tippy-arrow {
      border-bottom-color: ${props => props.theme.colors.neutralLighter};
    }
  }
`

const Popover = ({ children, content, ...props }) => (
  <StyledTippy content={content} {...props}>
    <div style={{ cursor: 'pointer' }}>{children}</div>
  </StyledTippy>
)

Popover.propTypes = {
  children: PropTypes.node.isRequired,
  content: PropTypes.node.isRequired,
}

export default Popover
