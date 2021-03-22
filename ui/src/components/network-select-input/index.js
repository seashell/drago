/* eslint-disable react/prop-types */
import PropTypes from 'prop-types'
import React from 'react'
import styled from 'styled-components'
import { icons } from '_assets/'
import Box from '_components/box'
import Icon from '_components/icon'
import SelectInput, { OptionContainer } from '_components/inputs/select-input'
import Text from '_components/text'

const StyledText = styled(Text).attrs({
  textStyle: 'body',
})`
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
`

const NetworkOptionComponent = ({ innerRef, innerProps, ...props }) => (
  <OptionContainer innerRef={innerRef} {...innerProps} {...props} color="neutral">
    <Icon mr={2} icon={<icons.Network />} size="32px" color="neutral" />
    <Box flexDirection="column">
      <StyledText textStyle="body">{props.data.name}</StyledText>
      <StyledText textStyle="detail" fontSize="10px">
        {props.data.addressRange}
      </StyledText>
    </Box>
  </OptionContainer>
)

const NetworkSingleValueComponent = ({ innerRef, innerProps, ...props }) => (
  <Box alignItems="center" height="100%" innerRef={innerRef} {...innerProps} {...props}>
    <Icon mr={2} icon={<icons.Network />} size="32px" color="neutral" />
    <Box flexDirection="column">
      <StyledText textStyle="body">{props.data.name}</StyledText>
      <StyledText textStyle="detail" fontSize="10px">
        {props.data.addressRange}
      </StyledText>
    </Box>
  </Box>
)

const NetworkSelectInput = ({ selectedId, networks, onChange, placeholder, ...props }) => {
  const handleChange = (option) => {
    onChange(option.id)
  }

  const network = networks.find((el) => el.ID === selectedId)

  const options = networks.map((el) => ({
    id: el.ID,
    name: el.Name,
    addressRange: el.AddressRange,
  }))

  const value = network
    ? { id: network.ID, name: network.Name, addressRange: network.AddressRange }
    : undefined

  return (
    <SelectInput
      options={options}
      value={value}
      optionComponent={NetworkOptionComponent}
      singleValueComponent={NetworkSingleValueComponent}
      onChange={handleChange}
      placeholder={placeholder}
      {...props}
    />
  )
}

NetworkSelectInput.propTypes = {
  networks: PropTypes.arrayOf(
    PropTypes.shape({
      ID: PropTypes.string,
      Name: PropTypes.string,
      AddressRange: PropTypes.string,
    })
  ),
  selectedId: PropTypes.string,
  onChange: PropTypes.func,
  placeholder: PropTypes.string,
}

NetworkSelectInput.defaultProps = {
  networks: [],
  selectedId: undefined,
  onChange: () => {},
  placeholder: 'Select network...',
}

export default NetworkSelectInput
