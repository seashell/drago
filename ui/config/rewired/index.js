const path = require('path')

const { override, useBabelRc, useEslintRc, addWebpackResolve } = require('customize-cra')

const rewireStyledComponents = require('react-app-rewire-styled-components')
const resolve = require('../webpack/resolve')

module.exports = {
  webpack: override(
    addWebpackResolve(resolve),
    // eslint-disable-next-line react-hooks/rules-of-hooks
    // useBabelRc(path.resolve(__dirname, '..', '..', '.babelrc')),
    // eslint-disable-next-line react-hooks/rules-of-hooks
    // useEslintRc(path.resolve(__dirname, '..', '..', '.eslintrc')),
    (config, env) => rewireStyledComponents(config, env, {})
  ),
}
