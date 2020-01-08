import { ApolloClient } from 'apollo-client'
import { HttpLink } from 'apollo-link-http'
import { WebSocketLink } from 'apollo-link-ws'
import { ApolloLink, split } from 'apollo-link'
import { InMemoryCache } from 'apollo-cache-inmemory'
import { getMainDefinition } from 'apollo-utilities'
import { withClientState } from 'apollo-link-state'
import { onError } from 'apollo-link-error'

import { defaults, resolvers } from './local-state'
import { GRAPHQL_API_URL, USE_WS_LINK } from '../environment'

const composeUrl = (url, protocol) => `${protocol}://${url}`

const httpLink = new HttpLink({
  uri: composeUrl(GRAPHQL_API_URL, 'http'),
})

const wsLink = new WebSocketLink({
  uri: composeUrl(GRAPHQL_API_URL, 'wss'),
  options: {
    reconnect: true,
  },
})

const authLink = new ApolloLink((operation, forward) => {
  operation.setContext(({ headers }) => ({
    headers: {
      Authorization: `Bearer ${localStorage.getItem('kc_jwt')}`,
      ...headers,
    },
  }))
  return forward(operation)
})

const errorLink = onError(({ graphQLErrors, networkError }) => {
  if (graphQLErrors)
    graphQLErrors.map(({ message, locations, path }) =>
      console.log(`[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`)
    )
  if (networkError) console.log(`[Network error]: ${networkError}`)
})

const terminatingLink = split(
  ({ query }) => {
    const { kind, operation } = getMainDefinition(query)
    return kind === 'OperationDefinition' && operation === 'subscription'
  },
  wsLink, // Receives the operation in case the expression above evaluates to true
  httpLink // Receives the operation otherwise
)

const cache = new InMemoryCache()

const stateLink = withClientState({
  cache,
  resolvers,
  defaults,
})

const link = ApolloLink.from(
  USE_WS_LINK
    ? [authLink, stateLink, errorLink, terminatingLink]
    : [authLink, stateLink, errorLink, httpLink]
)

export default new ApolloClient({
  link,
  cache,
})
