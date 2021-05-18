import { useLazyQuery, useQuery } from '@apollo/client'
import { useParams } from '@reach/router'
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
import NetworkSelectInput from '_components/network-select-input'
import Text from '_components/text'
import { GET_NETWORKS, GET_NODE_WITH_INTERFACES, GET_PEERS } from '_graphql/queries'
import PeerCard from './peer-card'

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

const ConnectPeerModal = ({ isOpen, onJoin, onClose }) => {
  const { nodeId } = useParams()

  const [networks, setNetworks] = useState([])
  const [selectedNetwork, setSelectedNetwork] = useState()

  const [peers, setPeers] = useState([])
  const [selectedPeer, setSelectedPeer] = useState()

  const [searchString, setSearchString] = useState('')

  const getNodeQuery = useQuery(GET_NODE_WITH_INTERFACES, {
    variables: { id: nodeId },
  })

  const handleGetNetworksQueryData = (data) => {
    setNetworks(data.result)
  }

  const handleGetPeersQueryData = (data) => {
    setPeers(data.result)
  }

  const getNetworksQuery = useQuery(GET_NETWORKS, {
    onCompleted: handleGetNetworksQueryData,
  })

  const [getPeers, getPeersQuery] = useLazyQuery(GET_PEERS, {
    onCompleted: handleGetPeersQueryData,
  })

  useEffect(() => {
    window.scrollTo(0, 0)
  }, [])

  useEffect(() => {
    if (isOpen) {
      setPeers([])
      setNetworks([])
      setSelectedPeer(undefined)
      setSelectedNetwork(undefined)
      if (getPeersQuery.called) {
        getNetworksQuery.refetch({}).then((res) => handleGetNetworksQueryData(res.data))
      }
      getPeersQuery.refetch()
      setSearchString('')
    }
  }, [isOpen])

  useEffect(() => {
    if (getPeersQuery.called) {
      getPeersQuery.refetch({ nodeId: '', networkId: selectedNetwork })
    } else {
      getPeers({
        variables: { nodeId: '', networkId: selectedNetwork },
      })
    }
  }, [selectedNetwork])

  const handleSelectedNetworkChanged = (id) => {
    setSelectedNetwork(id)
    setSelectedPeer(undefined)
  }

  const handlePeerCardClick = (id) => {
    setSelectedPeer(selectedPeer !== id ? id : undefined)
  }

  const handleJoinButtonClick = () => {
    onJoin(selectedNetwork, selectedPeer)
    onClose()
  }

  const node = getNodeQuery.data ? getNodeQuery.data.result : { Interfaces: [] }

  // Filter network options to include only those containing the current node
  const nodeNetworkIDs = node.Interfaces.map((el) => el.NetworkID)
  const nodeNetworks = networks.filter(
    (el) => nodeNetworkIDs.find((id) => id === el.ID) !== undefined
  )

  // Find node interface within the selected network
  const sourceInterface = node.Interfaces.find((el) => el.NetworkID === selectedNetwork)
  
  const filteredPeers = peers
    .filter((el) => el.ID !== sourceInterface.ID) // Do not show source interface as an option
    .filter((el) => (el.Name !== null ? el.Name.includes(searchString) : true)) // Filter peers based on search query

  return (
    <Modal isOpen={isOpen} onBackgroundClick={onClose} onEscapeKeydown={onClose}>
      <Container>
        <CloseButton onClick={onClose} />
        <Text textStyle="title" mb={3}>
          Connect to Peer
        </Text>
        <NetworkSelectInput
          networks={nodeNetworks}
          mb={3}
          onChange={handleSelectedNetworkChanged}
        />
        <SearchInput
          onChange={(s) => setSearchString(s)}
          mb={3}
          placeholder="Search for peers..."
        />
        <List pb={3}>
          {filteredPeers.map((el) => (
            <PeerCard
              key={el.ID}
              id={el.ID}
              name={el.Node.Name}
              address={el.Address}
              isSelected={el.ID === selectedPeer}
              onClick={() => handlePeerCardClick(el.ID)}
            />
          ))}
        </List>
        <Button
          disabled={selectedPeer === undefined}
          variant="primary"
          width="auto"
          height="48px"
          mt="auto"
          onClick={handleJoinButtonClick}
        >
          Connect
        </Button>
      </Container>
    </Modal>
  )
}

ConnectPeerModal.propTypes = {
  isOpen: PropTypes.bool,
  onJoin: PropTypes.func,
  onClose: PropTypes.func,
}

ConnectPeerModal.defaultProps = {
  isOpen: false,
  onJoin: () => {},
  onClose: () => {},
}

export default ConnectPeerModal
