import React from 'react'
import { Router } from '@reach/router'

import Tokens from './tokens'

const SettingsRouter = (props) => (
  <Router {...props}>
    <Tokens path="/tokens" />
  </Router>
)

export default SettingsRouter
