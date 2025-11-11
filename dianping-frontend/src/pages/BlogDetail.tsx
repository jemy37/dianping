import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { getBlog, likeBlog } from '../api/blog'
import BlogCard from '../components/BlogCard'

export default function BlogDetail() {
  const { id } = useParams()
  const blogId = Number(id)
  const [blog, setBlog] = useState<any>(null)
  const [msg, setMsg] = useState('')

  useEffect(() => {
    if (!blogId) return
    getBlog(blogId).then(r => setBlog(r.data?.data)).catch(() => setMsg('加载失败'))
  }, [blogId])

  async function like() {
    try {
      setMsg('')
      await likeBlog(blogId)
      const r = await getBlog(blogId)
      setBlog(r.data?.data)
    } catch (e: any) {
      setMsg(e?.response?.data?.msg || '操作失败')
    }
  }

  if (!blog) return <div className="card">加载中...</div>
  return (
    <div>
      <BlogCard blog={blog} onLike={like} />
      {msg && <div className="mt" style={{color:'#2563eb'}}>{msg}</div>}
    </div>
  )
}

