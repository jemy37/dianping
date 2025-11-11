import { useState } from 'react'
import { register, sendCode } from '../api/auth'

export default function Register() {
  const [phone, setPhone] = useState('')
  const [code, setCode] = useState('')
  const [password, setPassword] = useState('')
  const [nickName, setNickName] = useState('')
  const [msg, setMsg] = useState('')

  async function handleSendCode() {
    try {
      setMsg('')
      await sendCode(phone)
      setMsg('验证码已发送（开发环境请看后端日志）。')
    } catch (e: any) {
      setMsg(e?.response?.data?.msg || '发送失败')
    }
  }

  async function handleRegister() {
    try {
      setMsg('')
      await register({ phone, code, password, nickName })
      setMsg('注册成功，现在可以去登录')
    } catch (e: any) {
      setMsg(e?.response?.data?.msg || '注册失败')
    }
  }

  return (
    <div className="card">
      <h2>注册</h2>
      <div className="mt">
        <label>手机号</label>
        <input className="input" value={phone} onChange={e => setPhone(e.target.value)} />
      </div>
      <div className="mt row">
        <input className="input" value={code} onChange={e => setCode(e.target.value)} placeholder="验证码" />
        <button className="btn" onClick={handleSendCode}>发送验证码</button>
      </div>
      <div className="mt">
        <input className="input" value={password} onChange={e => setPassword(e.target.value)} placeholder="密码（6位以上）" type="password" />
      </div>
      <div className="mt">
        <input className="input" value={nickName} onChange={e => setNickName(e.target.value)} placeholder="昵称（可选）" />
      </div>
      <div className="mt">
        <button className="btn" onClick={handleRegister}>注册</button>
      </div>
      {msg && <p className="mt" style={{color:'#2563eb'}}>{msg}</p>}
    </div>
  )
}

