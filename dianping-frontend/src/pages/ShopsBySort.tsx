import { useEffect, useState } from 'react'
import { useSearchParams, Link } from 'react-router-dom'
import { listShopTypes, listShops } from '../api/shop'

export default function ShopsBySort() {
  const [params] = useSearchParams()
  const sort = Number(params.get('sort') || 0)
  const [shops, setShops] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    async function load() {
      setLoading(true)
      setError('')
      try {
        const stRes = await listShopTypes()
        const types: any[] = stRes.data?.data || []
        const matched = types.filter(t => t.sort === sort)
        const results: any[] = []
        for (const t of matched) {
          const sres = await listShops({ typeId: t.id, current: 1 })
          const list = sres.data?.data || []
          for (const s of list) {
            results.push({ ...s, _typeName: t.name })
          }
        }
        setShops(results)
      } catch (e: any) {
        setError(e?.response?.data?.msg || '加载失败')
      } finally {
        setLoading(false)
      }
    }
    if (sort) load()
  }, [sort])

  if (loading) return <div className="card">加载中...</div>
  if (error) return <div className="card">{error}</div>

  return (
    <div>
      <h2>按排序={sort} 的餐厅</h2>
      {shops.length === 0 && <div className="card">暂无数据</div>}
      <div className="grid">
        {shops.map(s => (
          <Link key={s.id} to={`/shop/${s.id}`} className="card">
            <strong>{s.name}</strong>
            <div style={{color:'#6b7280'}}>{s._typeName}</div>
          </Link>
        ))}
      </div>
    </div>
  )
}

