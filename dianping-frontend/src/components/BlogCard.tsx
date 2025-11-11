import { Link } from 'react-router-dom'

type Props = {
  blog: any
  onLike?: () => void
}

export default function BlogCard({ blog, onLike }: Props) {
  return (
    <div className="card">
      <Link to={`/blog/${blog.id}`}><strong>{blog.title || '无标题'}</strong></Link>
      <div className="mt">{blog.content}</div>
      <div className="mt row" style={{justifyContent:'space-between'}}>
        <span>点赞：{blog.liked || 0}</span>
        {onLike && <button className="btn" onClick={onLike}>点赞/取消</button>}
      </div>
    </div>
  )
}

