const path = require('path')

module.exports = {
  alias: {
    _assets: path.resolve(__dirname, '..', '..', 'src/assets/'),
    _components: path.resolve(__dirname, '..', '..', 'src/components/'),
    _containers: path.resolve(__dirname, '..', '..', 'src/containers/'),
    _modals: path.resolve(__dirname, '..', '..', 'src/modals/'),
    _graphql: path.resolve(__dirname, '..', '..', 'src/graphql/'),
    _utils: path.resolve(__dirname, '..', '..', 'src/utils/'),
    _views: path.resolve(__dirname, '..', '..', 'src/views/'),
  },
}
