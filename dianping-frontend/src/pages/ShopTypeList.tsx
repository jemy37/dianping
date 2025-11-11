import { Link } from 'react-router-dom'

// 固定映射：按 shop.typeId 分类，展示名称并按 typeId 升序
const TYPE_MAP: { id: number; name: string }[] = [
  { id: 1, name: '烧烤店' },
  { id: 2, name: '咖啡店' },
  { id: 5, name: '火锅店' }
]

export default function ShopTypeList() {
  const types = TYPE_MAP.sort((a, b) => a.id - b.id)
  return (
    <div>
      <h2>商铺类型</h2>
      <div className="mt" style={{color:'#6b7280'}}>
        返回 <Link to="/">全部商铺</Link>
      </div>
      <div className="grid">
        {types.map(t => (
          <div key={t.id} className="card">
            <Link to={`/shops?typeId=${t.id}`}><strong>{t.name}</strong></Link>
            <div className="mt" style={{color:'#6b7280'}}>typeId：{t.id}</div>
          </div>
        ))}
      </div>
    </div>
  )
}
