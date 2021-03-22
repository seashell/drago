import { handler } from './util'

export default {
  getSelf: handler((schema, request, self) => self.user.attrs),
}
