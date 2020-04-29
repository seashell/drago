import React, { useState, useRef, useEffect } from 'react'
import PropTypes from 'prop-types'
import styled from 'styled-components'
import { border, shadow } from 'styled-system'
import Modal from 'styled-react-modal'

import { useFormState } from 'react-use-form-state'

import { icons } from '_assets/'

import Button from '_components/button'
import IconButton from '_components/icon-button'

import Box from '_components/box'
import Text from '_components/text'
import TextInput from '_components/inputs/text-input'
import SearchInput from '_components/inputs/search-input'
import toast from '_components/toast'

import { GET_HOSTS, CREATE_LINK } from '_graphql/actions'
import { useQuery, useMutation } from 'react-apollo'

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
  padding-left: 12px;
  height: 48px;
  width: 100%;
`
const SearchContainer = styled(Box)`
  position: relative;
`

const SearchResults = styled(Box)`
  position: absolute;
  left: 0;
  right: 0;
  top: 48px;

  flex-direction: column;
  background: white;
  border: 1px solid #aaa;

  :empty {
    display: none;
  }
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

const NewLinkModal = ({ fromHost, onBackgroundClick, onEscapeKeydown, isOpen }) => {
  const searchInputRef = useRef(null)

  const [searchResults, setSearchResults] = useState([])

  const [formState, { text, number }] = useFormState({
    from: fromHost.id,
    to: null,
    allowedIPs: null,
    persistentKeepalive: null,
  })

  const getHostsQuery = useQuery(GET_HOSTS)

  const [createLink] = useMutation(CREATE_LINK, {
    variables: formState.values,
    onCompleted: () => {
      toast.success('Link created')
      onEscapeKeydown()
    },
    onError: () => {
      toast.error('Error creating link')
      onEscapeKeydown()
    },
  })

  useEffect(() => {
    getHostsQuery.refetch()
  }, [getHostsQuery])

  const handleSearchInputChanged = e => {
    if (!getHostsQuery.loading) {
      if (e.target.value === '') {
        setSearchResults([])
      } else {
        setSearchResults(
          getHostsQuery.data.result.items.filter(
            r =>
              r.name.includes(e.target.value) &&
              r.id !== fromHost.id &&
              fromHost.links.items.find(el => r.id === el.to.id) === undefined
          )
        )
      }
    }
  }

  const handleSearchResultSelected = selectedResult => {
    formState.setField('to', selectedResult.id)
    searchInputRef.current.value = `${selectedResult.name} (${selectedResult.address})`
    setSearchResults([])
  }

  const handleSearchInputFocused = () => {
    setSearchResults([])
  }

  const handleSearchInputBlurred = () => {
    setTimeout(() => {
      setSearchResults([])
    }, 500)
  }

  const handleCreateButtonClicked = () => {
    createLink()
  }

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
      <SearchContainer>
        <StyledSearchInput
          ref={searchInputRef}
          onChange={handleSearchInputChanged}
          onFocus={handleSearchInputFocused}
          onBlur={handleSearchInputBlurred}
          placeholder="Search for host name, or address"
        />
        {!getHostsQuery.loading && (
          <SearchResults>
            {searchResults.map(r => (
              <SearchResult key={r.id} onClick={() => handleSearchResultSelected(r)}>
                <IconContainer mr="12px">
                  <IconButton ml="auto" icon={<icons.Cube />} />
                </IconContainer>
                {r.name} ({r.address})
              </SearchResult>
            ))}
          </SearchResults>
        )}
      </SearchContainer>

      <Text my={3}>Allowed IPs</Text>
      <TextInput {...text('allowedIPs')} placeholder="0.0.0.0/0" mb={2} />

      <Text my={3}>Persistent keepalive</Text>
      <TextInput {...number('persistentKeepalive')} placeholder={20} mb={2} />

      <Button borderRadius={3} mt={3} mx="auto" onClick={handleCreateButtonClicked}>
        Create
      </Button>

      <StyledIconButton ml="auto" icon={<icons.Times />} onClick={onBackgroundClick} />
    </StyledModal>
  )
}

NewLinkModal.propTypes = {
  fromHost: PropTypes.node.isRequired,
  onBackgroundClick: PropTypes.func.isRequired,
  onEscapeKeydown: PropTypes.func.isRequired,
  isOpen: PropTypes.bool.isRequired,
}

export default NewLinkModal
