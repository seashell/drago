import React from 'react'
import { Router } from '@reach/router'

import Details from './details'
import New from './new'

const LinksRouter = () => (
  <Router>
    <New path="/new" />
    <Details path="/:linkId" />
  </Router>
)

export default LinksRouter
