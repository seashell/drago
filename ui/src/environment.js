export const GRAPHQL_API_URL = process.env.REACT_APP_GRAPHQL_API_URL || 'localhost:9002/graphql'

export const REST_API_URL =
  process.env.REACT_APP_GRAPHQL_API_URL ||
  'localhost:8080/api/v1/' ||
  '5e243331c5fc8f001465cef2.mockapi.io/api/v1/'

export const AUTH_PROVIDER_URL =
  process.env.REACT_APP_AUTH_PROVIDER_URL || 'http://localhost:8080/auth'

export const AUTH_PROVIDER_REALM = process.env.REACT_APP_AUTH_PROVIDER_REALM || 'master'

export const GOOGLE_MAPS_API_KEY =
  process.env.REACT_APP_GOOGLE_MAPS_API_KEY || 'enter-your-google-maps-api-key-here'

export const AUTH_PROVIDER_CLIENT_ID =
  process.env.REACT_APP_AUTH_PROVIDER_CLIENT_ID || 'keycloak-client-id'

export const USE_WS_LINK = process.env.REACT_APP_USE_WS_LINK || false
