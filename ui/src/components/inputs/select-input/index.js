import React from 'react'
import PropTypes from 'prop-types'

import { Container, StyledSelect } from './styled'

const SelectInput = props => {
  const { options, value, title, onChange } = props

  return (
    <Container {...props}>
      <StyledSelect
        value={value}
        placeholder="Select an option.."
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
  onChange: PropTypes.func,
}

SelectInput.defaultProps = {
  title: undefined,
  onChange: undefined,
}

export default SelectInput
