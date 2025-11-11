import { useState } from 'react'
import { createBlog } from '../api/blog'

export default function CreateBlog() {
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [msg, setMsg] = useState('')

  async function handleSubmit() {
    try {
      setMsg('')
      await createBlog({ title, content })
      setMsg('发布成功')
      setTitle('')
      setContent('')
    } catch (e: any) {
      setMsg(e?.response?.data?.msg || '发布失败（请先登录）')
    }
  }

  return (
    <div className="card">
      <h2>写博客</h2>
      <div className="mt">
        <input className="input" placeholder="标题" value={title} onChange={e => setTitle(e.target.value)} />
      </div>
      <div className="mt">
        <textarea className="input" placeholder="内容" rows={6} value={content} onChange={e => setContent(e.target.value)} />
      </div>
      <div className="mt">
        <button className="btn" onClick={handleSubmit}>发布</button>
      </div>
      {msg && <p className="mt" style={{color:'#2563eb'}}>{msg}</p>}
    </div>
  )
}

