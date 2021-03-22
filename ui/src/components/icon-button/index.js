import PropTypes from 'prop-types'
import React from 'react'
import styled, { css } from 'styled-components'
import { border, borderStyle, layout, space } from 'styled-system'

export const StyledButton = styled.button`
  display: flex;
  align-items: center;
  justify-content: center;

  border: none;
  box-shadow: none;
  background: none;
  box-sizing: content-box;
  padding: 8px;
  border-radius: 50%;
  outline: none;

  flex-shrink: 0;
  width: ${(props) => props.size};
  height: ${(props) => props.size};

  :disabled {
    border-color: ${(props) => props.theme.colors.border1};
    color: ${(props) => props.theme.colors.border1};
    box-shadow: none;
    cursor: not-allowed;
    svg {
      fill: ${(props) => props.theme.colors.border1};
    }
    :hover {
      background: none;
    }
  }

  svg {
    width: 100%;
    height: auto;
    fill: ${(props) =>
      props.color ? props.theme.colors[props.color] : props.theme.colors.foreground3};
  }

  ${(props) =>
    props.hoverEffect &&
    css`
      :hover {
        background: ${props.theme.colors.background2};
      }
    `}

  ${layout}
  ${space}
  ${border}
  ${borderStyle}
`

const IconButton = ({ icon, ...props }) => <StyledButton {...props}>{icon}</StyledButton>

IconButton.propTypes = {
  icon: PropTypes.node,
  size: PropTypes.string,
  color: PropTypes.string,
}

IconButton.defaultProps = {
  icon: undefined,
  color: 'foreground2',
  size: '24px',
}

export default IconButton
