import React from 'react'
import PropTypes from 'prop-types'
import { Redirect } from '@reach/router'

const protectRoute = WrappedComponent => {
  const ProtectedRoute = ({ location, ...props }) => {
    // TODO: Find out if user is authenticated
    const isAuthenticated = true
    if (!isAuthenticated) {
      return <Redirect to="/" state={{ referrer: location.pathname }} noThrow />
    }
    return <WrappedComponent location={location} {...props} />
  }

  ProtectedRoute.propTypes = {
    location: PropTypes.shape({
      pathname: PropTypes.string.isRequired,
    }).isRequired,
  }

  return ProtectedRoute
}

export default protectRoute
