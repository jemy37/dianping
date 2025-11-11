import { api } from './client'

export async function createBlog(payload: Record<string, any>) {
  return api.post('/api/blog', payload)
}

export async function likeBlog(id: number) {
  return api.put(`/api/blog/like/${id}`)
}

export async function listHotBlogs() {
  return api.get('/api/blog/hot')
}

export async function listBlogs(page = 1, size = 10) {
  return api.get('/api/blog', { params: { page, size } })
}

export async function myBlogs() {
  return api.get('/api/blog/of/me')
}

export async function getBlog(id: number) {
  return api.get(`/api/blog/${id}`)
}

export async function blogsOfFollow() {
  return api.get('/api/blog/of/follow')
}

export async function blogsOfShop(shopId: number, current = 1, size = 10) {
  return api.get(`/api/blog/of/shop/${shopId}`, { params: { current, size } })
}
