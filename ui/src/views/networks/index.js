import React from 'react'
import { Router } from '@reach/router'

import List from './list'
import New from './new'
import Details from './details'

const NetworksRouter = props => (
  <Router {...props}>
    <List path="/" />
    <New path="/new" />
    <Details path="/:networkId" />
  </Router>
)

export default NetworksRouter
