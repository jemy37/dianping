import { useEffect, useState } from 'react'
import { useSearchParams, Link } from 'react-router-dom'
import { listShops } from '../api/shop'

export default function ShopList() {
  const [params] = useSearchParams()
  const typeId = Number(params.get('typeId') || 0)
  const [shops, setShops] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    setLoading(true)
    listShops({ typeId, current: 1 })
      .then(res => setShops(res.data?.data?.list || res.data?.data || []))
      .catch(e => setError(e?.response?.data?.msg || '加载失败'))
      .finally(() => setLoading(false))
  }, [typeId])

  if (loading) return <div className="card">加载中...</div>
  if (error) return <div className="card">{error}</div>

  return (
    <div>
      <h2>商铺列表</h2>
      <div className="grid">
        {shops.map(s => (
          <Link key={s.id} to={`/shop/${s.id}`} className="card">
            <strong>{s.name}</strong>
            <div style={{color:'#6b7280'}}>{s.address}</div>
          </Link>
        ))}
      </div>
    </div>
  )
}
