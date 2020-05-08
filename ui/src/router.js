import React from 'react'
import { Router } from '@reach/router'

import App from '_views/app'
import NotFound from '_views/not-found'

export default () => (
  <Router>
    <App path="/*" />
    <NotFound default />
  </Router>
)
