/* eslint-disable react/prop-types */
import React from 'react'
import Select from 'react-select'
import styled from 'styled-components'
import { layout, space } from 'styled-system'

export const OptionContainer = styled.div`
  display: flex;
  cursor: pointer;
  height: 24px;
  padding: 8px;
  align-items: center;
  :hover {
    background-color: ${(props) => props.theme.colors.background2};
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
  styles: {
    control: (base) => ({
      ...base,
      border: 0,
      boxShadow: 'none',
    }),
  },
})`
  border: 1px solid ${(props) => props.theme.colors.neutralLighter};
  height: 48px;
  width: 100%;
  ${layout}
  ${space}

  .select__control {
    border: none;
    cursor: pointer;
    color: ${(props) => props.theme.colors.neutral};
    background-color: ${(props) => props.theme.colors.neutralLightest};

    height: 100%;
    .select__value-container {
      padding-left: 12px;
      height: 100%;
      .select__placeholder {
      }
    }
    :hover {
      border: none;
      background-color: ${(props) => props.theme.colors.neutralLightest};
    }
  }

  .select__menu {
    z-index: 2;
    color: ${(props) => props.theme.colors.neutral};
    background-color: ${(props) => props.theme.colors.neutralLightest};
  }

  .select__control--is-focused {
    box-shadow: none;
    border: none;
  }

  .select__control--is-disabled {
    opacity: 0.5;
    cursor: not-allowed;
    pointer-events: none;
  }

  .select__control--is-focused.select__control--menu-is-open {
    box-shadow: none;
    border-radius: 2px;
    border: 1px solid ${(props) => props.theme.colors.primary};
    :hover {
      border: 1px solid ${(props) => props.theme.colors.primary};
    }
  }
`

const SelectInput = ({ optionComponent, singleValueComponent, disabled, ...props }) => (
  <StyledSelect
    {...props}
    isDisabled={disabled}
    components={{ Option: optionComponent, SingleValue: singleValueComponent }}
  />
)

SelectInput.defaultProps = {
  disabled: false,
  optionComponent: DefaultOptionComponent,
  singleValueComponent: DefaultSingleValueComponent,
  options: [],
}

export default SelectInput
