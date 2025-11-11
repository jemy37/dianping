import { useEffect, useState } from 'react'
import { me, updateMe } from '../api/auth'

export default function Profile() {
  const [user, setUser] = useState<any>(null)
  const [nickName, setNickName] = useState('')
  const [icon, setIcon] = useState('')
  const [msg, setMsg] = useState('')

  useEffect(() => {
    me().then(res => {
      const u = res.data?.data
      setUser(u)
      setNickName(u?.nickName || '')
      setIcon(u?.icon || '')
    }).catch((e:any) => setMsg(e?.response?.data?.msg || '请先登录'))
  }, [])

  async function save() {
    try {
      setMsg('')
      await updateMe({ nickName, icon })
      setMsg('保存成功')
    } catch (e: any) {
      setMsg(e?.response?.data?.msg || '保存失败')
    }
  }

  if (!user) return <div className="card">加载中...</div>

  return (
    <div className="card">
      <h2>我的资料</h2>
      <div>手机号：{user.phone}</div>
      <div className="mt">
        <input className="input" value={nickName} onChange={e => setNickName(e.target.value)} placeholder="昵称" />
      </div>
      <div className="mt">
        <input className="input" value={icon} onChange={e => setIcon(e.target.value)} placeholder="头像URL" />
      </div>
      <div className="mt">
        <button className="btn" onClick={save}>保存</button>
      </div>
      {msg && <p className="mt" style={{color:'#2563eb'}}>{msg}</p>}
    </div>
  )
}

