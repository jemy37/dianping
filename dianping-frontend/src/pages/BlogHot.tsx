import { useEffect, useState } from 'react'
import { likeBlog, listHotBlogs } from '../api/blog'

export default function BlogHot() {
  const [blogs, setBlogs] = useState<any[]>([])
  const [msg, setMsg] = useState('')

  useEffect(() => {
    listHotBlogs().then(res => setBlogs(res.data?.data || [])).catch(() => {})
  }, [])

  async function doLike(id: number) {
    try {
      setMsg('')
      await likeBlog(id)
      setMsg('点赞成功')
    } catch (e: any) {
      setMsg(e?.response?.data?.msg || '点赞失败')
    }
  }

  return (
    <div>
      <h2>热门博客</h2>
      {blogs.map(b => (
        <div key={b.id} className="card">
          <strong>{b.title || '无标题'}</strong>
          <div className="mt">{b.content}</div>
          <div className="mt row" style={{justifyContent:'space-between'}}>
            <span>点赞：{b.liked || 0}</span>
            <button className="btn" onClick={() => doLike(b.id)}>点赞</button>
          </div>
        </div>
      ))}
      {msg && <p className="mt" style={{color:'#2563eb'}}>{msg}</p>}
    </div>
  )
}

