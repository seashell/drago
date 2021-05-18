import PropTypes from 'prop-types'
import React from 'react'
import styled, { css } from 'styled-components'
import { icons } from '_assets/'
import Box from '_components/box'
import Icon from '_components/icon'
import Text from '_components/text'

const Container = styled(Box).attrs({
  display: 'flex',
  p: 2,
})`
  height: max-content;
  cursor: pointer;
  align-items: center;
  border-bottom: 1px solid ${(props) => props.theme.colors.neutralLighter};
  :last-child {
    border-bottom: none;
  }
`

const StyledIcon = styled(Icon)`
  width: 36px;
  height: 36px;
  padding: 4px;
  border-radius: 4px;
  background: ${(props) => props.theme.colors.neutralLighter};
  align-items: center;
  justify-content: center;
`

const SelectionIndicator = styled(Box).attrs({})`
  border: 1px solid ${(props) => props.theme.colors.neutralLight};
  width: 10px;
  height: 10px;
  border-radius: 50%;
  ${(props) =>
    props.isSelected &&
    css`
      background: ${props.theme.colors.neutralDarker};
    `};
`

const PeerCard = ({ id, name, address, isSelected, onClick }) => (
  <Container onClick={() => onClick(id)}>
    <Box>
      <SelectionIndicator my="auto" mr={2} isSelected={isSelected} />
      <StyledIcon ml="auto" mr="12px" icon={<icons.Host />} color="neutralDarker" />
      <Box flexDirection="column">
        <Text textStyle="subtitle" fontSize="14px">
          {name}
        </Text>
        <Text textStyle="detail" fontSize="12px">
          {address ? `${address}` : 'NA'}
        </Text>
      </Box>
    </Box>
  </Container>
)

PeerCard.propTypes = {
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  address: PropTypes.string.isRequired,
  isSelected: PropTypes.bool,
  onClick: PropTypes.func,
}

PeerCard.defaultProps = {
  onClick: () => {},
  isSelected: false,
}

export default PeerCard
