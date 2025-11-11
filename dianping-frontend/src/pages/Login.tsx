import { useState } from 'react'
import { login, sendCode, loginWithPassword } from '../api/auth'
import { setToken } from '../store/auth'
import { useNavigate } from 'react-router-dom'

export default function Login() {
  const [mode, setMode] = useState<'code' | 'password'>('code')
  const [phone, setPhone] = useState('')
  const [code, setCode] = useState('')
  const [password, setPassword] = useState('')
  const [useNick, setUseNick] = useState(false)
  const [nickName, setNickName] = useState('')
  const [loading, setLoading] = useState(false)
  const [msg, setMsg] = useState('')
  const navigate = useNavigate()

  async function handleSendCode() {
    try {
      setMsg('')
      await sendCode(phone)
      setMsg('验证码已发送（开发环境请看后端日志）。')
    } catch (e: any) {
      setMsg(e?.response?.data?.msg || '发送失败')
    }
  }

  async function handleLogin() {
    setLoading(true)
    setMsg('')
    try {
      let token: string | null = null
      if (mode === 'code') {
        // 验证码登录
        const res = await login(phone, code)
        token = res.data?.data?.token
      } else {
        // 密码登录：支持手机号或昵称
        const payload: any = { password }
        if (useNick) payload.nickName = nickName
        else payload.phone = phone
        const res = await loginWithPassword(payload)
        token = res.data?.data?.token
      }
      if (!token) throw new Error('no-token')
      setToken(token)
      navigate('/')
    } catch (e: any) {
      setMsg(e?.response?.data?.msg || '登录失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="card">
      <div className="row">
        <button className="btn" onClick={() => setMode('code')} disabled={mode==='code'}>验证码登录</button>
        <button className="btn secondary" onClick={() => setMode('password')} disabled={mode==='password'}>密码登录</button>
      </div>

      {mode === 'code' ? (
        <>
          <h2 className="mt">验证码登录</h2>
          <div className="mt">
            <label>手机号</label>
            <input className="input" value={phone} onChange={e => setPhone(e.target.value)} placeholder="例如：13800000000" />
          </div>
          <div className="mt row">
            <input className="input" value={code} onChange={e => setCode(e.target.value)} placeholder="验证码" />
            <button className="btn" onClick={handleSendCode}>发送验证码</button>
          </div>
        </>
      ) : (
        <>
          <h2 className="mt">密码登录</h2>
          <div className="mt row">
            <label><input type="checkbox" checked={useNick} onChange={e => setUseNick(e.target.checked)} /> 使用昵称登录</label>
          </div>
          {useNick ? (
            <div className="mt">
              <input className="input" value={nickName} onChange={e => setNickName(e.target.value)} placeholder="昵称" />
            </div>
          ) : (
            <div className="mt">
              <input className="input" value={phone} onChange={e => setPhone(e.target.value)} placeholder="手机号" />
            </div>
          )}
          <div className="mt">
            <input className="input" type="password" value={password} onChange={e => setPassword(e.target.value)} placeholder="密码" />
          </div>
        </>
      )}

      <div className="mt">
        <button className="btn" disabled={loading} onClick={handleLogin}>{loading ? '登录中...' : '登录'}</button>
      </div>
      {msg && <p className="mt" style={{color:'#2563eb'}}>{msg}</p>}
    </div>
  )
}
