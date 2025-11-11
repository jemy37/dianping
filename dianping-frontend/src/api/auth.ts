import { api } from './client'

export async function sendCode(phone: string) {
  // Backend route is POST /api/user/code and reads phone from query string
  return api.post('/api/user/code', null, { params: { phone } })
}

export async function login(phone: string, code: string) {
  // POST /api/user/login { phone, code }
  return api.post('/api/user/login', { phone, code })
}

export async function loginWithPassword(payload: { phone?: string; nickName?: string; password: string }) {
  return api.post('/api/user/login/password', payload)
}

export async function register(payload: { phone: string; code: string; password: string; nickName?: string }) {
  return api.post('/api/user/register', payload)
}

export async function me() {
  return api.get('/api/user/me')
}

export async function updateMe(payload: { nickName?: string; icon?: string }) {
  return api.put('/api/user/update', payload)
}
