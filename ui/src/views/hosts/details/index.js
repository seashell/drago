import React, { useEffect } from 'react'
import { useQuery } from 'react-apollo'
import { Router, useLocation } from '@reach/router'
import styled from 'styled-components'
import PropTypes from 'prop-types'

import Box from '_components/box'
import Flex from '_components/flex'
import Text from '_components/text'
import { Dragon as Spinner } from '_components/spinner'
import { Nav, HorizontalNavLink as NavLink } from '_components/nav'
import IconButton from '_components/icon-button'

import { icons } from '_assets/'

import { GET_HOST } from '_graphql/actions/hosts'
import HostAttributesTab from './attributes-tab'
import HostInterfacesTab from './interfaces-tab'
import HostLinksTab from './links-tab'

const Container = styled(Flex)`
  flex-direction: column;
`

const IconContainer = styled(Box).attrs({
  display: 'flex',
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

const HostDetails = ({ hostId }) => {
  const location = useLocation()

  const getHostQuery = useQuery(GET_HOST, {
    variables: { id: hostId },
  })

  useEffect(() => {
    window.scrollTo(0, 0)
  }, [location])

  const isLoading = getHostQuery.loading

  return (
    <Container>
      {isLoading ? (
        <Spinner />
      ) : (
        <Box flexDirection="column">
          <Box alignItems="center" mb={4}>
            <IconContainer mr="12px">
              <IconButton ml="auto" size={32} icon={<icons.Host />} />
            </IconContainer>
            <Text textStyle="title">{getHostQuery.data.result.name}</Text>
          </Box>
          <Nav>
            <NavLink to="">Overview</NavLink>
            <NavLink to="interfaces">Interfaces</NavLink>
            <NavLink to="links">Links</NavLink>
          </Nav>
          <Router>
            <HostAttributesTab path="/" hostId={hostId} />
            <HostInterfacesTab path="/interfaces" hostId={hostId} />
            <HostLinksTab path="/links" hostId={hostId} />
          </Router>
        </Box>
      )}
    </Container>
  )
}

HostDetails.propTypes = {
  hostId: PropTypes.string.isRequired,
}

export default HostDetails
