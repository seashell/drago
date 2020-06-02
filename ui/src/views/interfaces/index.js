import React from 'react'
import { Router } from '@reach/router'

import Details from './details'
import New from './new'

const InterfacesRouter = () => (
  <Router>
    <New path="/new" />
    <Details path="/:interfaceId" />
  </Router>
)

export default InterfacesRouter
