import { ApolloClient } from 'apollo-client'
import { ApolloLink } from 'apollo-link'
import { RestLink } from 'apollo-link-rest'
import { InMemoryCache } from 'apollo-cache-inmemory'
import { onError } from 'apollo-link-error'

import log from 'loglevel'

import { REST_API_URL } from '../environment'
import { defaults } from './local-state'

const composeUrl = (url, protocol) => `${protocol}://${url}`

const restLink = new RestLink({
  uri: composeUrl(REST_API_URL, 'https'),
})

const authLink = new ApolloLink((operation, forward) => {
  operation.setContext(({ headers }) => ({
    headers: {
      'X-Drago-Token': `${localStorage.getItem('drago.settings.acl.token')}`,
      ...headers,
    },
  }))
  return forward(operation)
})

const errorLink = onError(({ graphQLErrors, networkError }) => {
  if (graphQLErrors)
    graphQLErrors.map(({ message, locations, path }) =>
      log.error(`[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`)
    )
  if (networkError) log.error(`[Network error]: ${networkError}`)
})

const cache = new InMemoryCache()
cache.writeData(defaults)

const link = ApolloLink.from([authLink, errorLink, restLink])

export default new ApolloClient({
  link,
  cache,
  typeDefs: {},
  connectToDevTools: true,
})
