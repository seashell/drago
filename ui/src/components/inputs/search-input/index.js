/* eslint-disable react/prop-types */
import React from 'react'
import styled from 'styled-components'
import { layout, space, border, color } from 'styled-system'

import { icons } from '_assets/'

import Select from 'react-select'

const StyledIcon = styled(icons.Search).attrs({
  width: 22,
})`
  fill: ${({ theme: { colors } }) => colors.neutralDark};
  margin: auto;
  position: absolute;
  left: 8px;
  top: 50%;
  transform: translateY(-50%);
  z-index: 99;
`

const Container = styled.div`
  position: relative;
  ${layout}
  ${space}
  ${border}
  ${color}
`

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
  height: 100%;
  border: 0px;

  .select__control {
    border: none;

    height: 100%;
    .select__value-container {
      padding-left: 32px;
      .select__placeholder {
      }
    }
    :hover {
      border: none;
    }
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
`

const SearchInput = ({
  options,
  placeholder,
  optionComponent,
  singleValueComponent,
  filterOption,
  onChange,
  ...props
}) => (
  <Container {...props}>
    <StyledIcon />
    <StyledSelect
      components={{ Option: optionComponent, SingleValue: singleValueComponent }}
      options={options}
      placeholder={placeholder}
      onChange={onChange}
      filterOption={filterOption}
      isSearchable
    />
  </Container>
)

SearchInput.defaultProps = {
  placeholder: 'Select...',
  options: [],
  optionRenderer: undefined,
}

export default SearchInput
