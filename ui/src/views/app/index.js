import React from 'react'
import styled from 'styled-components'

import NotFound from '_views/not-found'
import HostsView from '_views/hosts'
import Topology from '_views/topology'
import NewHost from '_views/new-host'

import Header from '_containers/header'
import Footer from '_containers/footer'

import { Router } from '@reach/router'
import EditHost from '_views/edit-host'

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

  min-height: 100vh;
  grid-area: body;

  width: 90%;
  max-width: 800px;
  justify-self: center;
`

const App = () => (
  <Dashboard>
    <Header />
    <Content>
      <HostsView path="/hosts" />
      <NewHost path="/hosts/new" />
      <EditHost path="/hosts/:hostId" />
      <Topology path="/topology" />
      <NotFound default />
    </Content>
    <Footer gridArea="footer" />
  </Dashboard>
)

export default App
