import React from 'react'
import ReactDOM from 'react-dom'

import { ApolloProvider } from 'react-apollo'
import { ModalProvider, BaseModalBackground } from 'styled-react-modal'
import styled, { ThemeProvider } from 'styled-components'

import { ToastContainer } from '_components/toast'

import Router from './router'
import client from './graphql/client'

import { themes, GlobalStyles } from './styles'
import * as serviceWorker from './serviceWorker'

import { makeMirageServer } from './mirage/server'

const theme = 'light'

const ModalBackground = styled(BaseModalBackground)`
  z-index: 999;
`

if (process.env.NODE_ENV === 'development') {
  makeMirageServer()
}

ReactDOM.render(
  <ApolloProvider client={client}>
    <ThemeProvider theme={themes[theme]}>
      <ModalProvider backgroundComponent={ModalBackground}>
        <GlobalStyles />
        <Router />
        <ToastContainer />
      </ModalProvider>
    </ThemeProvider>
  </ApolloProvider>,
  document.getElementById('root')
)

serviceWorker.unregister()
