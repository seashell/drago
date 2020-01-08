import React from 'react'
import ReactDOM from 'react-dom'

import Drift from 'react-driftjs'

import Keycloak from 'keycloak-js'
import { ApolloProvider } from 'react-apollo'
import { KeycloakProvider } from 'react-keycloak'
import { ThemeProvider } from 'styled-components'

import Router from './router'
import client from './graphql/client'
import * as serviceWorker from './serviceWorker'
import { themes, GlobalStyles } from './styles'

import { AUTH_PROVIDER_URL, AUTH_PROVIDER_REALM, AUTH_PROVIDER_CLIENT_ID } from './environment'

const theme = 'light'

const keycloak = Keycloak({
  url: AUTH_PROVIDER_URL,
  realm: AUTH_PROVIDER_REALM,
  clientId: AUTH_PROVIDER_CLIENT_ID,
})

ReactDOM.render(
  <ApolloProvider client={client}>
    <ThemeProvider theme={themes[theme]}>
      <KeycloakProvider keycloak={keycloak}>
        <GlobalStyles />
        <Router />
        <Drift appId="" />
      </KeycloakProvider>
    </ThemeProvider>
  </ApolloProvider>,
  document.getElementById('root')
)

serviceWorker.unregister()
