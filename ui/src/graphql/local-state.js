const defaults = {
  someState: 'value',
}

const resolvers = {
  Query: {},
  Mutation: {
    updateSomeState: (_, args, context) => {
      const { value } = args
      const { cache } = context
      cache.writeData({
        someState: {
          __typename: 'SomeState',
          value,
        },
      })
      return null
    },
  },
}

export { defaults, resolvers }
