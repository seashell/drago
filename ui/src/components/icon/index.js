import React from 'react'
import PropTypes from 'prop-types'
import styled, { css } from 'styled-components'
import { color, border } from 'styled-system'

import Box from '_components/box'

export const Container = styled(Box)`
  display: flex;
  align-items: center;
  justify-content: center;

  border: none;
  box-shadow: none;
  background: none;
  box-sizing: border-box;
  outline: none;

  flex-shrink: 0;
  width: ${(props) => props.size};
  height: ${(props) => props.size};

  svg {
    z-index: 1;
    width: 100%;
    height: auto;
    ${(props) =>
      props.color &&
      css`
        path {
          fill: ${props.theme.colors[props.color]};
        }
      `}
  }
  ${border}
  ${color}
`

const Icon = ({ icon, ...props }) => <Container {...props}>{icon}</Container>

Icon.propTypes = {
  icon: PropTypes.node,
  size: PropTypes.string,
  color: PropTypes.string,
}

Icon.defaultProps = {
  icon: undefined,
  color: undefined,
  size: '24px',
}

export default Icon
