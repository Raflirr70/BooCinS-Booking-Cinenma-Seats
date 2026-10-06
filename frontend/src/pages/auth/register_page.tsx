import { useState, type ChangeEvent, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import api, { apiMessage } from '../../lib/api'

type Form = { first_name: string; last_name: string; email: string; password: string }
const empty: Form = { first_name: '', last_name: '', email: '', password: '' }

const inputCls =
  'w-full rounded-lg bg-gray-800 border border-gray-700 px-3 py-2 text-white placeholder:text-gray-500 focus:border-red-500 focus:outline-none'

export default function RegisterPage() {
  const navigate = useNavigate()
  const [form, setForm] = useState<Form>(empty)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const field = (name: keyof Form, label: string, type: string, placeholder: string) => (
    <div>
      <label htmlFor={name} className="mb-1 block text-sm text-gray-400">{label}</label>
      <input
        id={name}
        type={type}
        placeholder={placeholder}
        value={form[name]}
        onChange={(e: ChangeEvent<HTMLInputElement>) => setForm({ ...form, [name]: e.target.value })}
        className={inputCls}
      />
    </div>
  )

  const submit = async (e: FormEvent) => {
    e.preventDefault()
    setError('')
    if (!form.email.includes('@')) { setError('Email tidak valid'); return }
    if (form.password.length < 6) { setError('Password minimal 6 karakter'); return }

    setLoading(true)
    try {
      await api.post('/auth/register', form)
      navigate('/login', { state: { registered: true } })
    } catch (err) {
      setError(apiMessage(err))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-900 px-4">
      <div className="w-full max-w-md rounded-2xl border border-gray-700 bg-gray-800 p-8">
        <h1 className="text-2xl font-bold text-white">Buat Akun</h1>
        <p className="mt-1 text-sm text-gray-400">Daftar untuk mulai menonton</p>

        <form onSubmit={submit} className="mt-6 space-y-4">
          {field('first_name', 'Nama Depan', 'text', 'Budi')}
          {field('last_name', 'Nama Belakang', 'text', 'Santoso')}
          {field('email', 'Email', 'email', 'budi@email.com')}
          {field('password', 'Password', 'password', 'Minimal 6 karakter')}

          {error && <p className="text-sm text-red-400">{error}</p>}

          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-lg bg-red-600 py-2.5 font-semibold text-white hover:bg-red-500 disabled:opacity-50"
          >
            {loading ? 'Mendaftar...' : 'Daftar'}
          </button>
        </form>

        <p className="mt-4 text-center text-sm text-gray-400">
          Sudah punya akun?{' '}
          <Link to="/login" className="text-red-400 hover:underline">Login</Link>
        </p>
      </div>
    </div>
  )
}