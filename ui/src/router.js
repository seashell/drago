import React from 'react'
import { Router } from '@reach/router'

import App from '_views/app'

export default () => (
  <Router>
    <App path="/ui/*" />
  </Router>
)
