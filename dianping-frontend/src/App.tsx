import { Link, NavLink, Route, Routes, useNavigate } from 'react-router-dom'
import Login from './pages/Login'
import Register from './pages/Register'
import ShopTypeList from './pages/ShopTypeList'
import AllShops from './pages/AllShops'
import ShopList from './pages/ShopList'
import ShopDetail from './pages/ShopDetail'
import Blogs from './pages/Blogs'
import BlogDetail from './pages/BlogDetail'
import Profile from './pages/Profile'
import ShopsBySort from './pages/ShopsBySort'
import { getToken, clearToken } from './store/auth'

function NavBar() {
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

export default function App() {
  return (
    <>
      <NavBar />
      <div className="container">
        <Routes>
          <Route path="/" element={<AllShops />} />
          <Route path="/types" element={<ShopTypeList />} />
          <Route path="/shops" element={<ShopList />} />
          <Route path="/shops/by-sort" element={<ShopsBySort />} />
          <Route path="/shop/:id" element={<ShopDetail />} />
          <Route path="/blog" element={<Blogs />} />
          <Route path="/blog/:id" element={<BlogDetail />} />
          <Route path="/me" element={<Profile />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
        </Routes>
      </div>
    </>
  )
}
