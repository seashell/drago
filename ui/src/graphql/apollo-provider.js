import { ApolloProvider, throwServerError } from '@apollo/client'
import { ApolloClient } from 'apollo-boost'
import { InMemoryCache } from 'apollo-cache-inmemory'
import { ApolloLink } from 'apollo-link'
import { setContext } from 'apollo-link-context'
import { onError } from 'apollo-link-error'
import { RestLink } from 'apollo-link-rest'
import log from 'loglevel'
import PropTypes from 'prop-types'
import React from 'react'
import { createNetworkStatusNotifier } from 'react-apollo-network-status'
import { useToast } from '_utils/toast-provider'
import { REST_API_URL } from '../environment'
import { defaults } from './local-state'

const {
  link: networkStatusLink,
  useApolloNetworkStatus: useNetworkStatus,
} = createNetworkStatusNotifier()

export async function customFetch(requestInfo, init) {
  const response = await fetch(requestInfo, init)
  const res = response.clone()

  if (!res.ok) {
    const body = await res.json()
    if (response.status === 404) {
      throwServerError(res, body, 'Not found error')
    }
    if (response.status === 500) {
      throwServerError(res, body, `Internal error: ${body.Message}`)
    }
  }
  return response
}

export const CustomApolloProvider = ({ children }) => {
  const { error } = useToast()
  const errorLink = onError(({ graphQLErrors, networkError }) => {
    if (graphQLErrors) {
      graphQLErrors.forEach(({ message, locations, path }) =>
        log.error(`[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`)
      )
    }
    if (networkError) {
      log.error('[Apollo Provider] Network error: ', networkError)
      if (networkError.statusCode > 300 && networkError.statusCode !== 404) {
        if (networkError.result && networkError.result.Message) {
          error(`${networkError.result.Message}`)
        }
      }
      // navigate('/ui/')
    }
  })

  const withToken = setContext(() => ({ token: localStorage.getItem('drago.settings.acl.token') }))

  const authLink = new ApolloLink((operation, forward) => {
    const { token } = operation.getContext()
    if (token && token !== null) {
      operation.setContext(({ headers }) => ({
        headers: {
          ...headers,
          'X-Drago-Token': `${token}`,
        },
      }))
    }
    return forward(operation)
  })

  const restLink = new RestLink({
    uri: REST_API_URL,
    customFetch,
  })

  const cache = new InMemoryCache({
    dataIdFromObject: (object) => {
      // eslint-disable-next-line no-underscore-dangle
      switch (object.__typename) {
        default:
          return object.ID
      }
    },
  })
  cache.writeData(defaults)

  const client = new ApolloClient({
    link: networkStatusLink.concat(ApolloLink.from([withToken, errorLink, authLink, restLink])),
    cache,
    typeDefs: {},
    connectToDevTools: true,
    queryDeduplication: true,
  })

  return <ApolloProvider client={client}>{children}</ApolloProvider>
}

CustomApolloProvider.propTypes = {
  children: PropTypes.node.isRequired,
}

export default CustomApolloProvider
