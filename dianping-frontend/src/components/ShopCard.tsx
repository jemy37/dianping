import { Link } from 'react-router-dom'

export default function ShopCard({ shop, subtitle }: { shop: any; subtitle?: string }) {
  return (
    <div className="card">
      <Link to={`/shop/${shop.id}`}><strong>{shop.name}</strong></Link>
      <div className="mt" style={{color:'#6b7280'}}>地址：{shop.address}</div>
      {subtitle && <div className="mt" style={{color:'#6b7280'}}>{subtitle}</div>}
    </div>
  )
}

