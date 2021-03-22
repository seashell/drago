import { useQuery } from '@apollo/client'
import PropTypes from 'prop-types'
import React, { useEffect, useState } from 'react'
import styled from 'styled-components'
import Modal from 'styled-react-modal'
import { icons } from '_assets/'
import Box from '_components/box'
import Button from '_components/button'
import IconButton from '_components/icon-button'
import SearchInput from '_components/inputs/search-input'
import List from '_components/list'
import Text from '_components/text'
import { GET_NETWORKS } from '_graphql/queries'
import NetworkCard from './network-card'

const Container = styled(Box)`
  height: 600px;
  width: 400px;
  background: ${(props) => props.theme.colors.white};
  border: 1px solid;
  border-color: ${(props) => props.theme.colors.neutralLighter};
  border-radius: 4px;
  padding: 32px;
  flex-direction: column;
  padding-bottom: 32px;
  position: relative;
`

const CloseButton = styled(IconButton).attrs({
  icon: <icons.Times />,
  size: '32px',
  color: 'neutralLight',
})`
  position: absolute;
  top: 0;
  right: 0;
`

const JoinNetworkModal = ({ isOpen, onJoin, onClose }) => {
  const [searchString, setSearchString] = useState('')
  const [selectedNetwork, setSelectedNetwork] = useState()
  const getNetworksQuery = useQuery(GET_NETWORKS, {})

  const handleNetworkCardClick = (id) => {
    setSelectedNetwork(selectedNetwork !== id ? id : undefined)
  }

  const handleJoinButtonClick = () => {
    onJoin(selectedNetwork)
    onClose()
  }

  useEffect(() => {
    window.scrollTo(0, 0)
  }, [])

  useEffect(() => {
    setSearchString('')
    setSelectedNetwork(undefined)
    getNetworksQuery.refetch()
  }, [isOpen])

  const networks = getNetworksQuery.data ? getNetworksQuery.data.result : []
  const filteredNetworks = networks.filter((el) => el.Name.includes(searchString))

  return (
    <Modal isOpen={isOpen} onBackgroundClick={onClose} onEscapeKeydown={onClose}>
      <Container>
        <CloseButton onClick={onClose} />
        <Text textStyle="title" mb={3}>
          Join Network
        </Text>
        <SearchInput
          onChange={(s) => setSearchString(s)}
          mb={3}
          placeholder="Search for networks..."
        />
        <List pb={3}>
          {filteredNetworks.map((el) => (
            <NetworkCard
              key={el.ID}
              id={el.ID}
              name={el.Name}
              ipAddressRange={el.AddressRange}
              numHosts={el.InterfacesCount}
              isSelected={el.ID === selectedNetwork}
              onClick={() => handleNetworkCardClick(el.ID)}
            />
          ))}
        </List>
        <Button
          disabled={selectedNetwork === undefined}
          variant="primary"
          width="auto"
          height="48px"
          mt="auto"
          onClick={handleJoinButtonClick}
        >
          Join
        </Button>
      </Container>
    </Modal>
  )
}

JoinNetworkModal.propTypes = {
  isOpen: PropTypes.bool,
  onJoin: PropTypes.func,
  onClose: PropTypes.func,
}

JoinNetworkModal.defaultProps = {
  isOpen: false,
  onJoin: () => {},
  onClose: () => {},
}

export default JoinNetworkModal
