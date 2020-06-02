/* eslint-disable react/prop-types */
import React from 'react'
import styled from 'styled-components'

import Select from 'react-select'
import { layout, space } from 'styled-system'

const OptionContainer = styled.div`
  display: flex;
  cursor: pointer;
  height: 40px;
  padding: 8px;
  align-items: center;
  :hover {
    background: #eee;
  }
`

const DefaultOptionComponent = ({ innerRef, innerProps, ...props }) => (
  <OptionContainer innerRef={innerRef} {...innerProps} {...props}>
    {props.data.label}
  </OptionContainer>
)

const DefaultSingleValueComponent = ({ innerRef, innerProps, ...props }) => (
  <div innerRef={innerRef} {...innerProps} {...props}>
    {props.data.label}
  </div>
)

const StyledSelect = styled(Select).attrs({
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
  ${layout}
  ${space}

  .select__control {
    border: none;

    height: 100%;
    .select__value-container {
      padding-left: 12px;
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

  .select__control--is-focused {
    box-shadow: none;
    border: none;
  }

  .select__control--is-focused.select__control--menu-is-open {
    box-shadow: none;
    border-radius: 2px;
    border: 1px solid ${props => props.theme.colors.primary};
    :hover {
      border: 1px solid ${props => props.theme.colors.primary};
    }
  }
`

const SelectInput = ({ optionComponent, singleValueComponent, ...props }) => (
  <StyledSelect
    {...props}
    components={{ Option: optionComponent, SingleValue: singleValueComponent }}
  />
)

SelectInput.defaultProps = {
  optionComponent: DefaultOptionComponent,
  singleValueComponent: DefaultSingleValueComponent,
  options: [],
}

export default SelectInput
