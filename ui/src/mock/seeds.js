export default (server) => {
  const create = (...args) => server.create(...args, 'withTime')
  const networks = {
    n1: create('network', { name: 'net1' }),
    n2: create('network', { name: 'net2' }),
  }
  const nodes = {
    n1: create('node', { name: 'node1' }),
    n2: create('node', { name: 'node2' }),
  }
  const interfaces = {
    i1: create('interface', { name: 'wg0', node: nodes.n1, network: networks.n1 }),
    i2: create('interface', { name: 'wg0', node: nodes.n2, network: networks.n2 }),
  }
  // eslint-disable-next-line no-unused-vars
  const connections = {
    c1: create('interface', { from: interfaces.i1, to: interfaces.i2 }),
  }
}
