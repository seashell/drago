import React from 'react'
import ReactDOM from 'react-dom'

import { ApolloProvider } from 'react-apollo'
import { ThemeProvider } from 'styled-components'
import { ToastContainer } from '_components/toast'

import Router from './router'
import client from './graphql/client'
import * as serviceWorker from './serviceWorker'
import { themes, GlobalStyles } from './styles'

const theme = 'light'

ReactDOM.render(
  <ApolloProvider client={client}>
    <ThemeProvider theme={themes[theme]}>
      <GlobalStyles />
      <Router />
      <ToastContainer />
    </ThemeProvider>
  </ApolloProvider>,
  document.getElementById('root')
)

serviceWorker.unregister()
