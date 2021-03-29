export const handler = (f) => (schema, request) => {
  const context = {}
  const self = {}
  request.body = JSON.parse(request.requestBody)
  request.headers = request.requestHeaders
  return f(schema, request, self, context)
}

export default {}
