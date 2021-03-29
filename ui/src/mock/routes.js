import { REST_API_URL } from '../environment'
import C from './commands'
import Q from './queries'

const routes = (server) => {
  const get = (path, handler) => server.get(`${REST_API_URL}${path}`, handler, { timing: 1000 })
  const post = (path, handler) => server.post(`${REST_API_URL}${path}`, handler, { timing: 2000 })

  // Setup query routes
  get('/api/self/token', Q.getSelf)

  // Setup command routes
  post('/api/networks/', C.createNetwork)

  // Setup passthrough routes (to which requests will not be intercepted by Mirage)
  // server.passthrough(<url>)
}

export default routes
