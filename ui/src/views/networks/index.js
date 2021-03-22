import { Router } from '@reach/router'
import React from 'react'
import Details from './details'
import List from './list'
import New from './new'

const NetworksRouter = (props) => (
  <Router {...props}>
    <List path="/" />
    <New path="/new" />
    <Details path="/:networkId" />
  </Router>
)

export default NetworksRouter
