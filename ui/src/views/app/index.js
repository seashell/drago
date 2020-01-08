import React from 'react'
import styled from 'styled-components'

import NotFound from '_views/not-found'

import Header from '_containers/header'
import Footer from '_containers/footer'

import { Router } from '@reach/router'
import DevicesView from '_views/devices'
import NewDevice from '_views/new-device'
import ConnectDevice from '_views/connect-device'

const Dashboard = styled.div`
  position: relative;
  display: grid;
  grid-template: 72px auto 40px / auto;
  grid-template-areas:
    'header'
    'body'
    'footer';
`

const Content = styled(Router).attrs({ primary: false })`
  padding-top: 44px;
  padding-bottom: 32px;
  margin: auto;
  min-height: 100vh;
  grid-area: body;
`

const App = () => (
  <Dashboard>
    <Header />
    <Content>
      <DevicesView path="/" />
      <NewDevice path="/devices/new" />
      <ConnectDevice path="/devices/:deviceId/connect" />
      <NotFound default />
    </Content>
    <Footer gridArea="footer" />
  </Dashboard>
)

export default App
