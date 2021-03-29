import React from 'react'
import styled from 'styled-components'
import { Portal } from 'react-portal'
import Box from '_components/box'
import { DEBUG } from '../environment'

const Container = styled(Box).attrs({
  border: 'discrete',
})`
  padding: 12px;
  position: absolute;
  bottom: 40px;
  left: 20px;
  height: auto;
`

const FormikState = (props) =>
  DEBUG ? (
    <Portal>
      <Container>
        <div>
          <pre
            style={{
              background: '#f6f8fa',
              fontSize: '0.8rem',
              padding: '.5rem',
              lineHeight: '1rem',
              lineBreak: 'auto',
              whiteSpace: 'pre-wrap',
            }}
          >
            {JSON.stringify(props, null, 4)}
          </pre>
        </div>
      </Container>
    </Portal>
  ) : null

export default FormikState
