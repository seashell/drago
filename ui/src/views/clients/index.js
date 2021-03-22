import { Router } from '@reach/router'
import React from 'react'
import Details from './details'
import List from './list'

const ClientsRouter = () => (
  <Router>
    <List path="/" />
    <Details path="/:nodeId/*" />
  </Router>
)

export default ClientsRouter
