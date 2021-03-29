import { Router } from '@reach/router'
import React from 'react'
import styled from 'styled-components'
import Footer from '_containers/footer'
import Header from '_containers/header'
import SideNav from '_containers/side-nav'
import ClientsRouter from '_views/clients'
import HomeView from '_views/home'
import NetworksRouter from '_views/networks'
import NotFound from '_views/not-found'
import SettingsRouter from '_views/settings'

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

  // Laptops and above
  @media (min-width: 1280px) {
    padding-left: 200px;
  }
`

const App = () => (
  <Dashboard>
    <Header />
    <SideNav />
    <Content>
      <HomeView path="/" />
      <ClientsRouter path="clients/*" />
      {/* <InterfacesRouter path="interfaces/*" /> */}
      {/* <LinksRouter path="links/*" /> */}
      <NetworksRouter path="networks/*" />
      <SettingsRouter path="settings/*" />
      <NotFound default />
    </Content>
    <Footer gridArea="footer" />
  </Dashboard>
)
export default App
