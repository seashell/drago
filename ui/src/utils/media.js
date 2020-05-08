import { css } from 'styled-components'

export default props => {
  console.log('yay')
  const queries = props.theme.mediaQueries
  const x = Object.keys(queries).reduce((acc, label) => {
    acc[label] = (...args) => css`
      ${queries[label]} {
        ${css(...args)};
      }
    `
    return acc
  }, {})
  console.log(x)
}
