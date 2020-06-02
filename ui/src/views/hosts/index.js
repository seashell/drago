import React from 'react'
import { Router } from '@reach/router'

import List from './list'
import Details from './details'
import New from './new'

const HostsRouter = () => (
  <Router>
    <List path="/" />
    <New path="/new" />
    <Details path="/:hostId/*" />
  </Router>
)

export default HostsRouter
