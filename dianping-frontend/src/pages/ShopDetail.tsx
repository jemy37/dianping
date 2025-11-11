import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { getShopById } from '../api/shop'
import { getSeckillVoucher, listVouchersByShop, seckillVoucher } from '../api/voucher'
import { blogsOfShop, createBlog } from '../api/blog'
import { getToken } from '../store/auth'

export default function ShopDetail() {
  const { id } = useParams()
  const shopId = Number(id)
  const [shop, setShop] = useState<any>(null)
  const [vouchers, setVouchers] = useState<any[]>([])
  const [seckill, setSeckill] = useState<any>(null)
  const [msg, setMsg] = useState('')
  const [blogs, setBlogs] = useState<any[]>([])
  const [newTitle, setNewTitle] = useState('')
  const [newContent, setNewContent] = useState('')

  useEffect(() => {
    if (!shopId) return
    getShopById(shopId).then(res => setShop(res.data?.data)).catch(() => {})
    listVouchersByShop(shopId).then(res => setVouchers(res.data?.data || [])).catch(() => {})
    blogsOfShop(shopId).then(res => setBlogs(res.data?.data?.list || [])).catch(() => {})
  }, [shopId])

  async function loadSeckill(id: number) {
    try {
      const res = await getSeckillVoucher(id)
      setSeckill(res.data?.data)
    } catch {}
  }

  async function doSeckill(id: number) {
    try {
      setMsg('')
      const res = await seckillVoucher(id)
      setMsg(res.data?.msg || '下单成功')
    } catch (e: any) {
      setMsg(e?.response?.data?.msg || '下单失败')
    }
  }

  if (!shop) return <div className="card">加载中...</div>

  return (
    <div>
      <h2>{shop.name}</h2>
      <div className="card">
        <div>地址：{shop.address}</div>
        <div>电话：{shop.phone}</div>
      </div>

      <h3>优惠券</h3>
      {vouchers.length === 0 && <div className="card">暂无优惠券</div>}
      {vouchers.map(v => (
        <div key={v.id} className="card">
          <div className="row" style={{justifyContent:'space-between'}}>
            <div>
              <strong>{v.title || '优惠券'}</strong>
              <div style={{color:'#6b7280'}}>满减：{v.payValue ?? v.discount ?? ''}</div>
            </div>
            {v.type === 1 ? (
              <button className="btn" onClick={() => loadSeckill(v.id)}>查看秒杀</button>
            ) : null}
          </div>
          {seckill && seckill.voucherId === v.id && (
            <div className="mt row" style={{justifyContent:'space-between'}}>
              <div>库存：{seckill.stock}，开始：{new Date(seckill.beginTime).toLocaleString()}</div>
              <button className="btn" onClick={() => doSeckill(v.id)}>立即秒杀</button>
            </div>
          )}
        </div>
      ))}

      {msg && <p className="mt" style={{color:'#2563eb'}}>{msg}</p>}

      <h3 className="mt">相关博客</h3>
      {blogs.length === 0 && <div className="card">还没有关于本店的博客</div>}
      {blogs.map(b => (
        <div key={b.id} className="card">
          <strong>{b.title || '无标题'}</strong>
          <div className="mt">{b.content}</div>
          <div className="mt" style={{color:'#6b7280'}}>点赞：{b.liked || 0}</div>
        </div>
      ))}

      {getToken() && (
        <div className="card mt">
          <h4>写一篇关于本店的博客</h4>
          <div className="mt">
            <input className="input" placeholder="标题" value={newTitle} onChange={e => setNewTitle(e.target.value)} />
          </div>
          <div className="mt">
            <textarea className="input" rows={5} placeholder="内容" value={newContent} onChange={e => setNewContent(e.target.value)} />
          </div>
          <div className="mt">
            <button className="btn" onClick={async () => {
              try {
                setMsg('')
                await createBlog({ title: newTitle, content: newContent, shopId: shopId })
                setNewTitle(''); setNewContent('')
                const res = await blogsOfShop(shopId)
                setBlogs(res.data?.data?.list || [])
                setMsg('发布成功')
              } catch (e: any) {
                setMsg(e?.response?.data?.msg || '发布失败')
              }
            }}>发布</button>
          </div>
        </div>
      )}
    </div>
  )
}
