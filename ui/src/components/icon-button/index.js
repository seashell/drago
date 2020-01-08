/* eslint-disable react/default-props-match-prop-types */
import React from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'

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
  width: ${props => props.size};
  height: ${props => props.size};

  svg {
    width: 100%;
    height: auto;
    fill: ${props => props.theme.colors.neutralDark};
  }

  &:hover {
    background: ${props => props.theme.colors.neutralLighter};
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }
`

const IconButton = ({ icon, ...props }) => <StyledButton {...props}>{icon}</StyledButton>

IconButton.propTypes = {
  icon: PropTypes.node,
  size: PropTypes.string,
  color: PropTypes.string,
}

IconButton.defaultProps = {
  icon: undefined,
  color: 'primary',
  size: '24px',
}

export default IconButton
