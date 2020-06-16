import React, { useState, useEffect } from 'react'
import PropTypes from 'prop-types'

import { useNavigate, useLocation } from '@reach/router'

import { useQuery, useMutation } from 'react-apollo'
import { GET_LINKS, DELETE_LINK } from '_graphql/actions/links'

import Box from '_components/box'
import Button from '_components/button'
import SearchInput from '_components/inputs/search-input'
import { Dragon as Spinner } from '_components/spinner'
import EmptyState from '_components/empty-state'

import LinkCard from './link-card'

const HostLinksTab = ({ hostId }) => {
  const navigate = useNavigate()
  const location = useLocation()

  const [searchFilter, setSearchFilter] = useState('')

  const getHostLinksQuery = useQuery(GET_LINKS, {
    variables: { fromHostId: hostId },
  })

  const [deleteLink, deleteLinkLinksMutation] = useMutation(DELETE_LINK, {
    variables: { id: undefined },
  })

  useEffect(() => {
    getHostLinksQuery.refetch()
  }, [location, getHostLinksQuery])

  const handleCreateButtonClicked = () => {
    navigate(`/ui/links/new`, { state: { hostId } })
  }

  const handleLinkSelect = id => {
    navigate(`/ui/links/${id}`, { state: { hostId } })
  }

  const handleLinkDelete = id => {
    deleteLink({ variables: id })
  }

  const handleSearchInputChanged = e => {
    setSearchFilter(e.target.value)
  }

  const isLoading = getHostLinksQuery.loading || deleteLinkLinksMutation.loading

  if (isLoading) {
    return <Spinner />
  }

  const filteredLinks = isLoading
    ? []
    : getHostLinksQuery.data.result.items.filter(
        el =>
          el.fromInterface.host.name.includes(searchFilter) ||
          el.toInterface.host.name.includes(searchFilter)
      )

  return (
    <Box flexDirection="column">
      <Box my={3}>
        <SearchInput
          width="100%"
          placeholder="Search..."
          onChange={handleSearchInputChanged}
          mr={2}
        />
        <Button onClick={handleCreateButtonClicked}>Create</Button>
      </Box>
      {filteredLinks.length === 0 ? (
        <EmptyState />
      ) : (
        filteredLinks.map(l => (
          <LinkCard
            key={l.id}
            id={l.id}
            fromInterface={l.fromInterface}
            toInterface={l.toInterface}
            allowedIps={l.allowedIps}
            persistentKeepalive={l.persistentKeepalive}
            onClick={() => {
              handleLinkSelect(l.id)
            }}
            onDelete={() => {
              handleLinkDelete(l.id)
            }}
          />
        ))
      )}
    </Box>
  )
}

HostLinksTab.propTypes = {
  hostId: PropTypes.string.isRequired,
}

export default HostLinksTab
