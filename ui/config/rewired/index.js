const path = require('path')

const {
  override,
  watchAll,
  useBabelRc,
  useEslintRc,
  addWebpackResolve,
  overrideDevServer,
} = require('customize-cra')

const resolve = require('../webpack/resolve')
const rewireStyledComponents = require('react-app-rewire-styled-components')

const devServerConfig = () => config => ({
  ...config,
  host: '0.0.0.0',
  port: '9999',
})

module.exports = {
  webpack: override(
    addWebpackResolve(resolve),
    useBabelRc(path.resolve(__dirname, '..', '..', '.babelrc')),
    useEslintRc(path.resolve(__dirname, '..', '..', '.eslintrc')),
    (config, env) => rewireStyledComponents(config, env, {})
  ),
  devServer: overrideDevServer(devServerConfig(), watchAll()),
}
