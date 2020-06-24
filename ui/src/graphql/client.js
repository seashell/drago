import { ApolloClient } from 'apollo-client'
import { ApolloLink } from 'apollo-link'
import { InMemoryCache } from 'apollo-cache-inmemory'
import { RestLink } from 'apollo-link-rest'

import { REST_API_URL } from '../environment'

const composeUrl = (url, protocol) => `${protocol}://${url}`

const restLink = new RestLink({
  uri: composeUrl(REST_API_URL || `${window.location.host}/api/`, 'http'),
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

const link = ApolloLink.from([authLink, restLink])

const inMemoryCache = new InMemoryCache()

export default new ApolloClient({
  link,
  cache: inMemoryCache,
  connectToDevTools: true,
  typeDefs: {},
})
