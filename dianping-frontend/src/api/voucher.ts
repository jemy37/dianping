import { api } from './client'

export async function listVouchersByShop(shopId: number) {
  return api.get(`/api/voucher/list/${shopId}`)
}

export async function addVoucher(payload: Record<string, any>) {
  return api.post('/api/voucher', payload)
}

export async function addSeckillVoucher(payload: Record<string, any>) {
  return api.post('/api/voucher/seckill', payload)
}

export async function getSeckillVoucher(id: number) {
  return api.get(`/api/voucher/seckill/${id}`)
}

export async function seckillVoucher(id: number) {
  // Protected route: POST /api/voucher-order/seckill/:id
  return api.post(`/api/voucher-order/seckill/${id}`)
}

