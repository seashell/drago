import { Router } from '@reach/router'
import React from 'react'
import ReactDOM from 'react-dom'
import styled, { ThemeProvider } from 'styled-components'
import { BaseModalBackground, ModalProvider } from 'styled-react-modal'
import ConfirmationDialogProvider from '_components/confirmation-dialog'
import ApolloProvider from '_graphql/apollo-provider'
import ToastProvider from '_utils/toast-provider'
import App from '_views/app'
import NotFound from '_views/not-found'
import * as serviceWorker from './serviceWorker'
import { GlobalStyles, themes } from './styles'

const theme = 'light'

const ModalBackground = styled(BaseModalBackground)`
  z-index: 999;
`

ReactDOM.render(
  <ThemeProvider theme={themes[theme]}>
    <ToastProvider>
      <ModalProvider backgroundComponent={ModalBackground}>
        <ConfirmationDialogProvider>
          <ApolloProvider>
            <GlobalStyles />
            <Router>
              <App path="/ui/*" />
              <NotFound default />
            </Router>
          </ApolloProvider>
        </ConfirmationDialogProvider>
      </ModalProvider>
    </ToastProvider>
  </ThemeProvider>,
  document.getElementById('root')
)

serviceWorker.unregister()
