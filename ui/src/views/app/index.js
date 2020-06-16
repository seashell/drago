import React from 'react'
import styled from 'styled-components'

import HomeView from '_views/home'
import NotFound from '_views/not-found'

import HostsRouter from '_views/hosts'
import LinksRouter from '_views/links'
import NetworksRouter from '_views/networks'
import InterfacesRouter from '_views/interfaces'
import SettingsRouter from '_views/settings'

import Header from '_containers/header'
import SideNav from '_containers/side-nav'
import Footer from '_containers/footer'

import { Router } from '@reach/router'

const Dashboard = styled.div`
  position: relative;
  display: grid;
  height: 100vh;
  grid-template: 72px auto 40px / auto;
  grid-template-areas:
    'header'
    'body'
    'footer';
`

const Content = styled(Router).attrs({ primary: false })`
  padding-top: 84px;
  padding-bottom: 32px;

  grid-area: body;

  width: 90%;
  max-width: 800px;
  justify-self: center;
`

const App = () => (
  <Dashboard>
    <Header />
    <SideNav />
    <Content>
      <HomeView path="/" />
      <HostsRouter path="hosts/*" />
      <InterfacesRouter path="interfaces/*" />
      <LinksRouter path="links/*" />
      <NetworksRouter path="networks/*" />
      <SettingsRouter path="settings/*" />
      <NotFound default />
    </Content>
    <Footer gridArea="footer" />
  </Dashboard>
)

export default App
