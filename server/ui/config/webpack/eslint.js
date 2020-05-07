const path = require('path')

const resolve = require('./resolve')

module.exports = {
  entry: ['../src/index'],
  output: {
    path: path.join(__dirname, 'dist'),
    filename: 'bundle.js',
    publicPath: '/static/',
  },
  plugins: [],
  resolve,
}
