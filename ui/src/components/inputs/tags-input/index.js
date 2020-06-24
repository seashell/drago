/* eslint-disable react/prop-types */
import React from 'react'
import styled from 'styled-components'

import { components } from 'react-select'
import CreatableSelect from 'react-select/creatable'

const MultiValueLabel = ({ data, ...props }) => (
  <components.MultiValueLabel {...props}>{data.label}</components.MultiValueLabel>
)

const StyledSelect = styled(CreatableSelect).attrs({
  classNamePrefix: 'select',
  style: {
    control: base => ({
      ...base,
      border: 0,
      boxShadow: 'none',
    }),
  },
})`
  border: 1px solid ${props => props.theme.colors.neutralLighter};
  height: 48px;
  width: 100%;

  .select__menu {
    display: none;
  }

  .select__control {
    border: none;

    height: 100%;
    .select__value-container {
      padding-left: ${props => (props.hasIcon ? '32px' : '8px')};
      .select__placeholder {
      }
    }
    :hover {
      border: none;
    }
  }

  .select__menu {
    z-index: 2;
  }

  .select__options {
    display: hidden;
  }

  .select__control--is-focused {
    box-shadow: none;
    border: none;
  }

  .select__control--is-focused.select__control--menu-is-open {
    box-shadow: none;
    border: 1px solid ${props => props.theme.colors.primary};
    :hover {
      border: 1px solid ${props => props.theme.colors.primary};
    }
  }

  .select__multi-value__remove {
    cursor: pointer;
    color: ${props => props.theme.colors.primary};
    :hover {
      color: ${props => props.theme.colors.primary};
      background: ${props => props.theme.colors.neutral};
    }
  }
`

const TagsInput = ({ ...props }) => (
  <StyledSelect
    {...props}
    isMulti
    components={{
      MultiValueLabel,
      DropdownIndicator: null,
    }}
  />
)

TagsInput.defaultProps = {}

export default TagsInput
