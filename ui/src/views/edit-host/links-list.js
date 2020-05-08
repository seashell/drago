import React from 'react'
import PropTypes from 'prop-types'

import styled from 'styled-components'

import Box from '_components/box'
import Text from '_components/text'
import Button from '_components/button'

import LinkCard from './link-card'

const Container = styled(Box)`
  display: grid;
  grid-template-columns: repeat(4, 1fr);
`

const EmptyStateContainer = styled(Box).attrs({
  border: 'discrete',
  height: '100px',
})``

const AddLinkCard = styled(Box).attrs({
  display: 'flex',
  border: 'discrete',
  m: 2,
  p: 3,
})`
  border-style: dashed;
  align-items: center;
  justify-content: center;
  height: 200px;
  grid-column: span 1;
  flex-direction: column;
  position: relative;
  cursor: pointer;
`
const EmptyState = () => (
  <EmptyStateContainer>
    <Text>[TODO: Empty state] - No links found</Text>
  </EmptyStateContainer>
)

const LinksList = ({ links, onLinkAdd, onLinkUpdate, onLinkDelete }) => {
  const handleAddLinkButtonClick = () => {
    onLinkAdd()
  }

  return (
    <Container>
      <AddLinkCard>
        <Button variant="primaryInverted" onClick={handleAddLinkButtonClick}>
          Add link
        </Button>
      </AddLinkCard>
      {links.map(l => (
        <LinkCard
          key={l.id}
          id={l.id}
          toName={l.to.name}
          toAddress={l.to.address}
          allowedIPs={l.allowedIPs}
          persistentKeepalive={l.persistentKeepalive}
          onChange={onLinkUpdate}
          onDelete={e => onLinkDelete(e, l.id)}
        />
      ))}
    </Container>
  )
}
LinksList.propTypes = {
  links: PropTypes.arrayOf(PropTypes.object).isRequired,
  onLinkAdd: PropTypes.func.isRequired,
  onLinkUpdate: PropTypes.func.isRequired,
  onLinkDelete: PropTypes.func.isRequired,
}

export default LinksList
