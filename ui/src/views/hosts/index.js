import React from 'react'
import styled from 'styled-components'

import { navigate } from '@reach/router'

import { useQuery, useMutation } from 'react-apollo'

import { GET_NODES, DELETE_NODE } from '_graphql/actions'
import { Dragon } from '_components/spinner'
import Text from '_components/text'
import Flex from '_components/flex'
import toast from '_components/toast'
import Button from '_components/button'
import Box from '_components/box'

import NodesList from './nodes-list'

const Container = styled(Flex)`
  flex-direction: column;
`

export const StyledButton = styled(Button)``

const NodesView = () => {
  const { loading, data, refetch } = useQuery(GET_NODES)

  const handleNodeDeleted = () => {
    toast.success('Node deleted')
    refetch()
  }

  const handleNodeDeleteError = () => {
    toast.error('Error deleting node')
  }

  const [deleteNode, { loading: deleting }] = useMutation(DELETE_NODE, {
    variables: { id: undefined },
    onCompleted: handleNodeDeleted,
    onError: handleNodeDeleteError,
  })

  const handleNodeSelect = id => {
    navigate(`/nodes/${id}`)
  }

  const handleNodeDelete = (e, id) => {
    e.preventDefault()
    e.stopPropagation()
    deleteNode({ variables: { id } })
  }

  const handleCreateNodeClick = () => {
    navigate('/nodes/new')
  }

  return (
    <Container>
      <Box mb={3}>
        <Text textStyle="title">Nodes</Text>
        <Button
          onClick={handleCreateNodeClick}
          variant="primary"
          borderRadius={3}
          width="100px"
          height="40px"
          ml="auto"
        >
          Create
        </Button>
      </Box>
      {loading || deleting ? (
        <Dragon />
      ) : (
        <NodesList
          nodes={data.result.items}
          onNodeSelect={handleNodeSelect}
          onNodeDelete={handleNodeDelete}
        />
      )}
    </Container>
  )
}

export default NodesView
