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

const NodeOptionComponent = ({ innerRef, innerProps, ...props }) => (
  <OptionContainer innerRef={innerRef} {...innerProps} {...props} color="neutral">
    <Icon mr={2} icon={<icons.Host />} size="32px" color="neutral" />
    <Box flexDirection="column">
      <StyledText textStyle="body">{props.data.name}</StyledText>
      <StyledText textStyle="detail" fontSize="10px">
        {props.data.id}
      </StyledText>
    </Box>
  </OptionContainer>
)

const NodeSingleValueComponent = ({ innerRef, innerProps, ...props }) => (
  <Box alignItems="center" height="100%" innerRef={innerRef} {...innerProps} {...props}>
    <Icon mr={2} icon={<icons.Host />} size="32px" color="neutral" />
    <Box flexDirection="column">
      <StyledText textStyle="body">{props.data.name}</StyledText>
      <StyledText textStyle="detail" fontSize="10px">
        {props.data.id}
      </StyledText>
    </Box>
  </Box>
)

const NodeSelectInput = ({ selectedId, nodes, onChange, placeholder, ...props }) => {
  const handleChange = (option) => {
    onChange(option.id)
  }

  const node = nodes.find((el) => el.ID === selectedId)

  const options = nodes.map((el) => ({
    id: el.ID,
    name: el.Name,
    addressRange: el.AddressRange,
  }))

  const value = node ? { id: node.ID, name: node.Name } : undefined

  return (
    <SelectInput
      options={options}
      value={value}
      optionComponent={NodeOptionComponent}
      singleValueComponent={NodeSingleValueComponent}
      onChange={handleChange}
      placeholder={placeholder}
      {...props}
    />
  )
}

NodeSelectInput.propTypes = {
  nodes: PropTypes.arrayOf(
    PropTypes.shape({
      ID: PropTypes.string,
      Name: PropTypes.string,
    })
  ),
  placeholder: PropTypes.string,
  selectedId: PropTypes.string,
  onChange: PropTypes.func,
}

NodeSelectInput.defaultProps = {
  nodes: [],
  selectedId: undefined,
  placeholder: 'Select node...',
  onChange: () => {},
}

export default NodeSelectInput
