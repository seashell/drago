import React from 'react'
import PropTypes from 'prop-types'

import { Container, StyledSelect } from './styled'

const SelectInput = props => {
  const { options, value, title, placeholder, onChange } = props

  return (
    <Container {...props}>
      <StyledSelect
        value={value}
        placeholder={placeholder}
        onChange={onChange}
        options={options}
        searchable
      />
    </Container>
  )
}

SelectInput.propTypes = {
  options: PropTypes.arrayOf(
    PropTypes.shape({
      value: PropTypes.string,
      label: PropTypes.string,
    })
  ).isRequired,
  title: PropTypes.string,
  value: PropTypes.string.isRequired,
  placeholder: PropTypes.string,
  onChange: PropTypes.func,
}

SelectInput.defaultProps = {
  title: undefined,
  placeholder: 'Select an option..',
  onChange: undefined,
}

export default SelectInput
