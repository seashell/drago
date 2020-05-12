/* eslint-disable react/prop-types */
import React, { useEffect } from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'
import { border, shadow } from 'styled-system'
import Modal from 'styled-react-modal'
import { useFormState } from 'react-use-form-state'

import { useQuery } from 'react-apollo'
import * as _ from 'lodash'

import { icons } from '_assets/'

import Button from '_components/button'
import IconButton from '_components/icon-button'

import Box from '_components/box'
import Text from '_components/text'
import TextInput from '_components/inputs/text-input'
import SearchInput from '_components/inputs/search-input'

import { GET_HOSTS } from '_graphql/actions'

const StyledModal = Modal.styled`
background: white;
width: 480px;
height: 620px;
display: flex;
flex-direction: column;
border-radius: 4px;
position: relative;
padding: 32px;
${border}
${shadow}
`

const StyledIconButton = styled(IconButton)`
  position: absolute;
  right: 2px;
  top: 2px;
  svg {
    transform: scale(1);
  }
  :hover {
    background: ${props => props.theme.colors.neutralLighter};
  }
`

const StyledSearchInput = styled(SearchInput)`
  border: 1px solid ${props => props.theme.colors.neutralLighter};
  height: 48px;
  width: 100%;
`

const SearchResult = styled(Box)`
  cursor: pointer;
  height: 48px;
  padding: 8px;
  align-items: center;
  :hover {
    background: #eee;
  }
`

const IconContainer = styled(Box).attrs({
  display: 'flex',
  height: '40px',
  width: '40px',
  bg: 'neutralLighter',
  borderRadius: '4px',
})`
  position: relative;
  button {
    margin-right: auto;
  }
  align-items: center;
  justify-content: center;
`

const NewLinkModal = ({
  networkId,
  fromHost,
  onCreateLink,
  onBackgroundClick,
  onEscapeKeydown,
  isOpen,
}) => {
  const [formState, { raw }] = useFormState({
    from: fromHost.id,
    to: null,
    allowedIPs: [],
    persistentKeepalive: null,
  })

  const getHostsQuery = useQuery(GET_HOSTS, {
    variables: { networkId },
  })

  useEffect(() => {
    getHostsQuery.refetch()
  }, [getHostsQuery])

  const HostOption = ({ innerRef, innerProps, ...props }) => (
    <SearchResult innerRef={innerRef} {...innerProps} {...props}>
      <IconContainer mr="12px">
        <IconButton ml="auto" icon={<icons.Cube />} />
      </IconContainer>
      {props.data.name} ({props.data.ipAddress})
    </SearchResult>
  )

  const SelectedHostValue = ({ innerRef, innerProps, ...props }) => (
    <div innerRef={innerRef} {...innerProps} {...props}>
      {props.data.name} ({props.data.ipAddress})
    </div>
  )
  const handleTargetHostSelected = targetHost => {
    formState.setField('to', targetHost.id)
  }

  const handleCreateButtonClicked = e => {
    e.preventDefault()
    e.stopPropagation()
    onCreateLink(formState)
    onEscapeKeydown()
  }

  const filterHosts = (option, searchText) => {
    if (
      option.data.name.toLowerCase().includes(searchText.toLowerCase()) ||
      option.data.ipAddress.toLowerCase().includes(searchText.toLowerCase())
    ) {
      return true
    }
    return false
  }

  const options = !getHostsQuery.loading ? getHostsQuery.data.result.items : []
  return (
    <StyledModal
      isOpen={isOpen}
      border="discrete"
      boxShadow="medium"
      onEscapeKeydown={onEscapeKeydown}
      onBackgroundClick={onBackgroundClick}
    >
      <Text textStyle="subtitle" mt={3} mb={2}>
        New link
      </Text>

      <Text my={3}>To</Text>
      <StyledSearchInput
        options={options}
        optionComponent={HostOption}
        singleValueComponent={SelectedHostValue}
        placeholder="Search for host name, or address"
        filterOption={filterHosts}
        onChange={handleTargetHostSelected}
      />

      <Text my={3}>Allowed IPs</Text>
      <TextInput
        {...raw({
          name: 'allowedIPs',
          compare(initialValue, value) {
            return _.isEqual(value.sort(), initialValue.sort())
          },
          value: formState.allowedIPs,
          onChange: e => e.target.value.replace(/\s/g, '').split(','),
        })}
        placeholder="0.0.0.0/0"
        mb={2}
      />

      <Text my={3}>Persistent keepalive</Text>
      <TextInput
        type="number"
        {...raw({
          name: 'persistentKeepalive',
          compare(initialValue, value) {
            return value === initialValue
          },
          value: formState.persistentKeepalive,
          onChange: e => parseInt(e.target.value, 10),
        })}
        placeholder={20}
        mb={2}
      />

      <Button borderRadius={3} mt={3} mx="auto" onClick={handleCreateButtonClicked}>
        Create
      </Button>

      <StyledIconButton ml="auto" icon={<icons.Times />} onClick={onBackgroundClick} />
    </StyledModal>
  )
}

NewLinkModal.propTypes = {
  networkId: PropTypes.string.isRequired,
  fromHost: PropTypes.node.isRequired,
  onCreateLink: PropTypes.func.isRequired,
  onBackgroundClick: PropTypes.func.isRequired,
  onEscapeKeydown: PropTypes.func.isRequired,
  isOpen: PropTypes.bool.isRequired,
}

export default NewLinkModal
