import { useMutation, useQuery } from '@apollo/client'
import { navigate, useLocation } from '@reach/router'
import React, { useEffect, useState } from 'react'
import { FixedSizeList } from 'react-window'
import InfiniteLoader from 'react-window-infinite-loader'
import styled from 'styled-components'
import Box from '_components/box'
import Button from '_components/button'
import { useConfirmationDialog } from '_components/confirmation-dialog'
import EmptyState from '_components/empty-state'
import ErrorState from '_components/error-state'
import Flex from '_components/flex'
import SearchInput from '_components/inputs/search-input'
import { Dragon as Spinner } from '_components/spinner'
import Text from '_components/text'
import { DELETE_NETWORK } from '_graphql/mutations'
import { GET_NETWORKS } from '_graphql/queries'
import { useToast } from '_utils/toast-provider'
import NetworkCard from './network-card'

const Container = styled(Flex)`
  flex-direction: column;
`

const NetworksView = () => {
  const [searchFilter, setSearchFilter] = useState('')
  const location = useLocation()

  const { success } = useToast()
  const { confirm } = useConfirmationDialog()

  const getNetworksQuery = useQuery(GET_NETWORKS, {})

  const [deleteNetwork, deleteNetworkMutation] = useMutation(DELETE_NETWORK, {
    variables: { id: undefined },
    onCompleted: () => {
      success('Network deleted')
      getNetworksQuery.refetch()
    },
  })

  useEffect(() => {
    window.scrollTo(0, 0)
    getNetworksQuery.refetch()
  }, [location])

  const handleNetworkSelect = (id) => {
    navigate(`networks/${id}`)
  }

  const handleNetworkDelete = (id) => {
    confirm({
      title: 'Are you sure?',
      details: `This will destroy all resources, including interfaces and connections, associated with this network.`,
      isDestructive: true,
      onConfirm: () => {
        deleteNetwork({ variables: { id } })
          .then(() => {
            getNetworksQuery.refetch()
          })
          .catch(() => {})
      },
    })
  }

  const handleCreateNetworkClick = () => {
    navigate('networks/new')
  }

  const networks = getNetworksQuery.data ? getNetworksQuery.data.result : []

  const isError = getNetworksQuery.error
  const isLoading = getNetworksQuery.loading || deleteNetworkMutation.loading

  const isEmpty = networks.length === 0

  const filteredNetworks = networks.filter(
    (el) =>
      el.ID.includes(searchFilter) ||
      el.Name.includes(searchFilter) ||
      el.AddressRange.includes(searchFilter)
  )

  if (isLoading) {
    return <Spinner />
  }

  return (
    <Container>
      <Box mb={3}>
        <Text textStyle="title">Networks</Text>
        <Button onClick={handleCreateNetworkClick} variant="primary" ml="auto">
          Create
        </Button>
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
        <ErrorState message="" />
      ) : isEmpty ? (
        <EmptyState message="Oops! It seems that you don't have any networks yet registered." />
      ) : (
        <InfiniteLoader
          itemCount={networks.length}
          isItemLoaded={(index) => index + 1 <= networks.length}
          loadMoreItems={() => {}}
        >
          {({ onItemsRendered, ref }) => (
            <FixedSizeList
              className="virtualized-list"
              height={74 * filteredNetworks.length}
              itemCount={filteredNetworks.length}
              itemSize={74}
              onItemsRendered={onItemsRendered}
              itemData={filteredNetworks}
              ref={ref}
              width={'100%'}
            >
              {({ index, style }) => {
                const network = filteredNetworks[index]
                return (
                  <NetworkCard
                    key={network.ID}
                    id={network.ID}
                    name={network.Name}
                    addressRange={network.AddressRange}
                    interfacesCount={network.InterfacesCount}
                    connectionsCount={network.ConnectionsCount}
                    onClick={handleNetworkSelect}
                    onDelete={() => handleNetworkDelete(network.ID)}
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

export default NetworksView
