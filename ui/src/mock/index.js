import { createServer, RestSerializer } from 'miragejs'

import seeds from './seeds'
import models from './models'
import routes from './routes'
import factories from './factories'

import IdentityManager from './identity'

export default () => {
  createServer({
    identityManagers: {
      application: IdentityManager,
    },
    serializers: {
      application: RestSerializer,
    },
    models,
    seeds,
    routes() {
      routes(this)
    },
    factories,
  })
}
