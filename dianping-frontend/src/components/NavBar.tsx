import { Link, NavLink, useNavigate } from 'react-router-dom'
import { getToken, clearToken } from '../store/auth'

export default function NavBar() {
  const navigate = useNavigate()
  const token = getToken()
  return (
    <div className="nav">
      <div className="nav-inner">
        <Link to="/" className="brand">Dianping</Link>
        <NavLink to="/" end>首页</NavLink>
        <NavLink to="/blog">博客</NavLink>
        <NavLink to="/types">类型</NavLink>
        <div className="spacer" />
        {token ? (
          <>
            <NavLink to="/me">我的</NavLink>
            <button className="btn secondary" onClick={() => { clearToken(); navigate('/'); }}>退出</button>
          </>
        ) : (
          <>
            <NavLink to="/login">登录</NavLink>
            <NavLink to="/register">注册</NavLink>
          </>
        )}
      </div>
    </div>
  )
}

