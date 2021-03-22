import { Model, belongsTo, hasMany } from 'miragejs'

export default {
  networks: Model.extend({}),

  nodes: Model.extend({
    interfaces: hasMany('interface'),
  }),

  interfaces: Model.extend({
    node: belongsTo('node'),
    network: belongsTo('network'),
    connections: hasMany('connection'),
  }),

  connections: Model.extend({
    from: belongsTo('interface'),
    to: belongsTo('interface'),
  }),
}
