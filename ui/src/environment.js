export const DEBUG = process.env.REACT_APP_DEBUG || false

const defaultApiURL = process.env.NODE_ENV === 'development' ? 'http://localhost:8080' : ''

export const REST_API_URL = process.env.REACT_APP_REST_API_URL || defaultApiURL