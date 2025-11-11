import axios from 'axios'
import { getToken } from '../store/auth'

// Create a pre-configured axios instance.
// Base URL is left empty so relative paths ("/api/...") go through Vite proxy in dev.
export const api = axios.create({
  // Optionally set baseURL: '/api' to always prefix. We keep it explicit in calls for clarity.
  timeout: 10000
})

// Attach Authorization header for protected routes.
api.interceptors.request.use(config => {
  const token = getToken()
  if (token) {
    // Backend expects Bearer token in Authorization header.
    config.headers = config.headers || {}
    config.headers['Authorization'] = `Bearer ${token}`
  }
  return config
})

// Unwrap common result structure if needed.
// Assuming backend returns { code, msg, data } via utils.Response.
api.interceptors.response.use(
  (resp) => resp,
  (error) => Promise.reject(error)
)

