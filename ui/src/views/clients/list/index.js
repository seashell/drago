import { useQuery } from '@apollo/client'
import { navigate, useLocation } from '@reach/router'
import React, { useEffect, useState } from 'react'
import { FixedSizeList } from 'react-window'
import InfiniteLoader from 'react-window-infinite-loader'
import styled from 'styled-components'
import Box from '_components/box'
import EmptyState from '_components/empty-state'
import ErrorState from '_components/error-state'
import SearchInput from '_components/inputs/search-input'
import { Dragon as Spinner } from '_components/spinner'
import Text from '_components/text'
import { GET_NODES } from '_graphql/queries'
import ClientCard from './client-card'

const Container = styled(Box)`
  flex-direction: column;
`

const ClientsListView = () => {
  const location = useLocation()

  const [searchFilter, setSearchFilter] = useState('')
  const getNodesQuery = useQuery(GET_NODES, {
    variables: {},
  })

  useEffect(() => {
    window.scrollTo(0, 0)
    getNodesQuery.refetch()
  }, [location])

  const handleClientCardClick = (id) => {
    navigate(`/ui/clients/${id}`)
  }

  const nodes = getNodesQuery.data ? getNodesQuery.data.result : []

  const isError = getNodesQuery.error
  const isLoading = getNodesQuery.loading

  const isEmpty = false

  const filteredNodes = nodes.filter(
    (el) =>
      el.ID.includes(searchFilter) ||
      el.Name.includes(searchFilter) ||
      el.Status.includes(searchFilter)
  )

  if (isLoading) {
    return <Spinner />
  }

  return (
    <Container>
      <Box mb={3}>
        <Text textStyle="title">Clients</Text>
      </Box>
      <Box my={3}>
        <SearchInput
          width="100%"
          placeholder="Search..."
          onChange={(value) => setSearchFilter(value)}
          mr={2}
        />
      </Box>
      {isError ? (
        <ErrorState />
      ) : isEmpty ? (
        <EmptyState description="Oops! It seems that you don't have any clients yet registered in this network." />
      ) : (
        <InfiniteLoader
          itemCount={nodes.length}
          isItemLoaded={(index) => index + 1 <= nodes.length}
          loadMoreItems={() => {}}
        >
          {({ onItemsRendered, ref }) => (
            <FixedSizeList
              className="virtualized-list"
              height={72 * filteredNodes.length}
              itemCount={filteredNodes.length}
              itemSize={72}
              onItemsRendered={onItemsRendered}
              itemData={filteredNodes}
              ref={ref}
              width={'100%'}
            >
              {({ index, style }) => {
                const node = filteredNodes[index]
                return (
                  <ClientCard
                    key={node.ID}
                    id={node.ID}
                    name={node.Name}
                    status={node.Status}
                    address={node.Address}
                    connectionCount={node.ConnectionsCount}
                    interfaceCount={node.InterfacesCount}
                    onClick={() => handleClientCardClick(node.ID)}
                    updatedAt={node.UpdatedAt}
                    createdAt={node.CreatedAt}
                    style={style}
                  />
                )
              }}
            </FixedSizeList>
          )}
        </InfiniteLoader>
      )}
    </Container>
  )
}

ClientsListView.propTypes = {}

export default ClientsListView
