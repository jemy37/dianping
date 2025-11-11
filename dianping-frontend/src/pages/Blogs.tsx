import { useEffect, useState } from 'react'
import { listBlogs, listHotBlogs } from '../api/blog'
import BlogCard from '../components/BlogCard'

export default function Blogs() {
  const [mode, setMode] = useState<'all' | 'hot'>('all')
  const [blogs, setBlogs] = useState<any[]>([])
  const [error, setError] = useState('')

  async function load() {
    try {
      setError('')
      if (mode === 'all') {
        const r = await listBlogs(1, 20)
        const list = r.data?.data?.list || r.data?.data || []
        setBlogs(list)
      } else {
        const r = await listHotBlogs()
        const list = r.data?.data?.list || r.data?.data || []
        setBlogs(list)
      }
    } catch (e: any) {
      setError(e?.response?.data?.msg || '加载失败')
    }
  }

  useEffect(() => { load() }, [mode])

  return (
    <div>
      <h2>博客</h2>
      <div className="row mt">
        <button className="btn" onClick={() => setMode('all')} disabled={mode==='all'}>全部</button>
        <button className="btn secondary" onClick={() => setMode('hot')} disabled={mode==='hot'}>热门</button>
      </div>
      {error && <div className="card mt">{error}</div>}
      {blogs.map(b => (
        <BlogCard key={b.id} blog={b} />
      ))}
    </div>
  )
}

