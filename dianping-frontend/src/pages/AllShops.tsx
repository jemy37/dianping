import { useEffect, useMemo, useState } from 'react'
import { listAllShops, listShopTypes } from '../api/shop'
import ShopCard from '../components/ShopCard'

// 首页：展示所有商铺，附带 name/address/typeId
// 排序逻辑：按照所属类型的 sort 升序排序（sort 越小越靠前）
export default function AllShops() {
  const [shops, setShops] = useState<any[]>([])
  const [types, setTypes] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    async function load() {
      try {
        setLoading(true)
        setError('')
        const [shopRes, typeRes] = await Promise.all([
          listAllShops(1, 1000), // 取较大的 size 以覆盖小型数据集
          listShopTypes()
        ])
        const list = shopRes.data?.data?.list || shopRes.data?.data || []
        const tps = typeRes.data?.data || []
        setShops(list)
        setTypes(tps)
      } catch (e: any) {
        setError(e?.response?.data?.msg || '加载失败')
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [])

  const typeSortMap = useMemo(() => {
    const m = new Map<number, number>()
    for (const t of types) {
      m.set(t.id, t.sort)
    }
    return m
  }, [types])

  const sorted = useMemo(() => {
    const copy = [...shops]
    copy.sort((a, b) => {
      const sa = typeSortMap.get(a.typeId) ?? 1e9
      const sb = typeSortMap.get(b.typeId) ?? 1e9
      if (sa !== sb) return sa - sb
      // 次级排序：按名称
      return String(a.name).localeCompare(String(b.name))
    })
    return copy
  }, [shops, typeSortMap])

  if (loading) return <div className="card">加载中...</div>
  if (error) return <div className="card">{error}</div>

  return (
    <div>
      <h2>全部商铺</h2>
      <div className="grid">
        {sorted.map(s => (
          <ShopCard key={s.id} shop={s} subtitle={`typeId：${s.typeId}（排序：${typeSortMap.get(s.typeId) ?? '-'})`} />
        ))}
      </div>
    </div>
  )
}
