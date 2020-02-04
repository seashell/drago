import React from 'react'
import styled from 'styled-components'

import NotFound from '_views/not-found'
import NodesView from '_views/nodes'
import NodesGraph from '_views/nodes-graph'
import NewNode from '_views/new-node'

import Header from '_containers/header'
import Footer from '_containers/footer'

import { Router } from '@reach/router'
import EditNode from '_views/edit-node'

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
      <NodesView path="/ui/nodes" />
      <NewNode path="/ui/nodes/new" />
      <EditNode path="/ui/nodes/:nodeId" />
      <NodesGraph path="/ui/nodes/graph" />
      <NotFound default />
    </Content>
    <Footer gridArea="footer" />
  </Dashboard>
)

export default App
