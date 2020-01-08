import React from 'react'
import styled from 'styled-components'

const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  grid-column: span 12;
`

const NotFound = () => (
  <Container>
    <p>404 Not found</p>
  </Container>
)

export default NotFound
