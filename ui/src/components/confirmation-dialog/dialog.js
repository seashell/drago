import React from 'react'
import { Slide } from 'react-awesome-reveal'
import styled from 'styled-components'
import Box from '_components/box'
import Button from '_components/button'
import Text from '_components/text'

const Container = styled(Box).attrs({
  bg: 'white',
  border: '1px solid neutralLighter',
})`
  flex-direction: column;

  border-radius: 4px;
  overflow: hidden;

  width: 100%;
  max-width: 80%;
  height: 200px;

  // Tablets (portrait)
  @media (min-width: 768px) and (max-width: 1024px) {
    width: auto;
    max-width: 600px;

    min-width: 420px;
    min-height: 240px;
  }

  // Laptops and above
  @media (min-width: 1280px) {
    width: auto;
    max-width: 48px;

    min-width: 420px;
    min-height: 240px;
  }
`

const ButtonsArea = styled(Box)`
  justify-content: stretch;
  height: 60px;
  * + * {
    border-left: 1px solid ${(props) => props.theme.colors.neutralLighter} !important;
  }
`

const StyledButton = styled(Button)`
  border-radius: 0;
  background: transparent;
  border-top: 1px solid ${(props) => props.theme.colors.neutralLighter} !important;
  :hover {
    background: ${(props) => props.theme.colors.white};
  }
`

// eslint-disable-next-line react/prop-types
const Dialog = ({ title, details, isDestructive, onConfirm, onCancel }) => (
  <Slide direction="up" duration={500}>
    <Container>
      <Box flexDirection="column" alignItems="center" justifyContent="center" flex="1">
        <Text mb={3} px={4} textAlign="center" textStyle="title" fontSize="24px">
          {title}
        </Text>
        <Text px={'48px'} lineHeight="18px" textStyle="body" textAlign="center" color="neutralDark">
          {details}
        </Text>
      </Box>
      <ButtonsArea>
        <StyledButton
          width="100%"
          height="100%"
          color={isDestructive ? 'primary' : 'neutralDark'}
          onClick={onCancel}
        >
          Cancel
        </StyledButton>
        <StyledButton
          width="100%"
          height="100%"
          color={isDestructive ? 'danger' : 'primary'}
          onClick={onConfirm}
          isDestructive
        >
          Confirm
        </StyledButton>
      </ButtonsArea>
    </Container>
  </Slide>
)

export default Dialog
