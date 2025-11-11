import { api } from './client'

export async function listShopTypes() {
  return api.get('/api/shop-type/list')
}

export async function listShops(params: { typeId?: number; current?: number; name?: string }) {
  // /api/shop/of/type?typeId=1&current=1 or /api/shop/of/name?name=xxx&current=1
  if (params.name) {
    return api.get('/api/shop/of/name', { params })
  }
  return api.get('/api/shop/of/type', { params })
}

export async function getShopById(id: number) {
  return api.get(`/api/shop/${id}`)
}

export async function createShop(payload: Record<string, any>) {
  return api.post('/api/shop/createShop', payload)
}

export async function updateShop(payload: Record<string, any>) {
  return api.put('/api/shop/update', payload)
}

export async function listAllShops(page = 1, size = 200) {
  // Backend GetShopList uses query params page/size
  return api.get('/api/shop/list', { params: { page, size } })
}
