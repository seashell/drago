import faker from 'faker'
import { Factory, trait } from 'miragejs'

export default {
  token: Factory.extend({
    withTime: trait({
      createdAt: () => faker.date.recent(),
      updatedAt: () => faker.date.recent(),
    }),
  }),
  network: Factory.extend({
    withTime: trait({
      createdAt: () => faker.date.recent(),
      updatedAt: () => faker.date.recent(),
    }),
  }),
  node: Factory.extend({
    withTime: trait({
      createdAt: () => faker.date.recent(),
      updatedAt: () => faker.date.recent(),
    }),
  }),
  interface: Factory.extend({
    withTime: trait({
      createdAt: () => faker.date.recent(),
      updatedAt: () => faker.date.recent(),
    }),
  }),
  connection: Factory.extend({
    withTime: trait({
      createdAt: () => faker.date.recent(),
      updatedAt: () => faker.date.recent(),
    }),
  }),
}
