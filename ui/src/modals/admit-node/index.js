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
import { GET_NODES } from '_graphql/queries'
import NodeCard from './node-card'

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

const AdmitNodeModal = ({ isOpen, onAdmit, onClose }) => {
  const [searchString, setSearchString] = useState('')
  const [selectedNode, setSelectedNode] = useState()
  const getNodesQuery = useQuery(GET_NODES, {})

  const handleNodeCardClick = (id) => {
    setSelectedNode(selectedNode !== id ? id : undefined)
  }

  const handleAdmitButtonClick = () => {
    onAdmit(selectedNode)
    onClose()
  }

  useEffect(() => {
    window.scrollTo(0, 0)
  }, [])

  useEffect(() => {
    setSearchString('')
    setSelectedNode(undefined)
    getNodesQuery.refetch()
  }, [isOpen])

  const networks = getNodesQuery.data ? getNodesQuery.data.result : []
  const filteredNetworks = networks.filter((el) => el.Name.includes(searchString))

  return (
    <Modal isOpen={isOpen} onBackgroundClick={onClose} onEscapeKeydown={onClose}>
      <Container>
        <CloseButton onClick={onClose} />
        <Text textStyle="title" mb={3}>
          Admit Node
        </Text>
        <SearchInput
          onChange={(s) => setSearchString(s)}
          mb={3}
          placeholder="Search for nodes..."
        />
        <List pb={3}>
          {filteredNetworks.map((el) => (
            <NodeCard
              key={el.ID}
              id={el.ID}
              name={el.Name}
              ipAddressRange={el.AddressRange}
              numHosts={el.InterfacesCount}
              isSelected={el.ID === selectedNode}
              onClick={() => handleNodeCardClick(el.ID)}
            />
          ))}
        </List>
        <Button
          disabled={selectedNode === undefined}
          variant="primary"
          width="auto"
          height="48px"
          mt="auto"
          onClick={handleAdmitButtonClick}
        >
          Admit
        </Button>
      </Container>
    </Modal>
  )
}

AdmitNodeModal.propTypes = {
  isOpen: PropTypes.bool,
  onClose: PropTypes.func,
  onAdmit: PropTypes.func,
}

AdmitNodeModal.defaultProps = {
  isOpen: false,
  onAdmit: () => {},
  onClose: () => {},
}

export default AdmitNodeModal
