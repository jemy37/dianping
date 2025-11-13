import { FormEvent, useState } from 'react'
import { createShop } from '../api/shop'

const initialForm = {
  name: '',
  typeId: '',
  area: '',
  address: '',
  avgPrice: '',
  openHours: ''
}

export default function CreateShop() {
  const [form, setForm] = useState(initialForm)
  const [submitting, setSubmitting] = useState(false)
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null)

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    const { name, value } = e.target
    setForm(prev => ({ ...prev, [name]: value }))
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setMessage(null)

    // 所有字段必填
    for (const [key, value] of Object.entries(form)) {
      if (!value.trim()) {
        setMessage({ type: 'error', text: '请完整填写所有字段' })
        return
      }
    }

    const payload = {
      name: form.name.trim(),
      typeId: Number(form.typeId),
      area: form.area.trim(),
      address: form.address.trim(),
      avgPrice: Number(form.avgPrice),
      openHours: form.openHours.trim(),
      images: '[]',
      typeIcon: '',
      x: 0,
      y: 0
    }

    if (!payload.typeId || payload.typeId <= 0) {
      setMessage({ type: 'error', text: '类型ID必须为正整数' })
      return
    }
    if (!payload.avgPrice || payload.avgPrice <= 0) {
      setMessage({ type: 'error', text: '人均价格必须为正整数' })
      return
    }

    setSubmitting(true)
    try {
      await createShop(payload)
      setMessage({ type: 'success', text: '创建成功！' })
      setForm(initialForm)
    } catch (error: any) {
      const errText = error?.response?.data?.msg || error?.response?.data?.errorMsg || '创建失败，请重试'
      setMessage({ type: 'error', text: errText })
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="card">
      <h2>新增商铺</h2>
      <p className="mt">所有字段均为必填，提交后将调用后端创建接口。</p>
      <form className="mt" onSubmit={handleSubmit}>
        <div className="mt">
          <label>商铺名称</label>
          <input className="input" name="name" value={form.name} onChange={handleChange} placeholder="如：阿光米线" />
        </div>
        <div className="mt">
          <label>类型 ID</label>
          <input className="input" name="typeId" value={form.typeId} onChange={handleChange} placeholder="请输入类型 ID" />
        </div>
        <div className="mt">
          <label>商圈/区域</label>
          <input className="input" name="area" value={form.area} onChange={handleChange} placeholder="如：浦东新区" />
        </div>
        <div className="mt">
          <label>详细地址</label>
          <input className="input" name="address" value={form.address} onChange={handleChange} placeholder="如：世纪大道 88 号" />
        </div>
        <div className="mt">
          <label>人均价格（元）</label>
          <input className="input" name="avgPrice" value={form.avgPrice} onChange={handleChange} placeholder="如：68" />
        </div>
        <div className="mt">
          <label>营业时间</label>
          <input className="input" name="openHours" value={form.openHours} onChange={handleChange} placeholder="如：08:00-22:00" />
        </div>
        <div className="mt">
          <button className="btn" type="submit" disabled={submitting}>{submitting ? '创建中...' : '创建商铺'}</button>
        </div>
      </form>
      {message && (
        <p className="mt" style={{ color: message.type === 'success' ? '#15803d' : '#dc2626' }}>{message.text}</p>
      )}
    </div>
  )
}
