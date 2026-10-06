import { useState, type ChangeEvent, type FormEvent } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import api, { apiMessage } from '../../lib/api'
import { useAuth, homePathFor, type AuthUser } from '../../store/auth'

const inputCls =
  'w-full rounded-lg bg-gray-800 border border-gray-700 px-3 py-2 text-white placeholder:text-gray-500 focus:border-red-500 focus:outline-none'

export default function LoginPage() {
  const navigate = useNavigate()
  const location = useLocation()
  const setAuth = useAuth((s) => s.setAuth)

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const registered = (location.state as { registered?: boolean } | null)?.registered

  const submit = async (e: FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const res = await api.post('/auth/login', { email, password })
      const { user, token } = res.data.data as { user: AuthUser; token: string }
      setAuth(token, user)
      navigate(homePathFor(user.role.name), { replace: true })
    } catch (err) {
      setError(apiMessage(err))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-900 px-4">
      <div className="w-full max-w-md rounded-2xl border border-gray-700 bg-gray-800 p-8">
        <h1 className="text-2xl font-bold text-white">Login</h1>
        <p className="mt-1 text-sm text-gray-400">Masuk ke akun kamu</p>

        <form onSubmit={submit} className="mt-6 space-y-4">
          <div>
            <label htmlFor="email" className="mb-1 block text-sm text-gray-400">Email</label>
            <input
              id="email"
              type="email"
              placeholder="budi@email.com"
              value={email}
              onChange={(e: ChangeEvent<HTMLInputElement>) => setEmail(e.target.value)}
              className={inputCls}
            />
          </div>
          <div>
            <label htmlFor="password" className="mb-1 block text-sm text-gray-400">Password</label>
            <input
              id="password"
              type="password"
              placeholder="Password kamu"
              value={password}
              onChange={(e: ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)}
              className={inputCls}
            />
          </div>

          {registered && !error && (
            <p className="text-sm text-green-400">Registrasi berhasil, silakan login.</p>
          )}
          {error && <p className="text-sm text-red-400">{error}</p>}

          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-lg bg-red-600 py-2.5 font-semibold text-white hover:bg-red-500 disabled:opacity-50"
          >
            {loading ? 'Masuk...' : 'Masuk'}
          </button>
        </form>

        <p className="mt-4 text-center text-sm text-gray-400">
          Belum punya akun?{' '}
          <Link to="/register" className="text-red-400 hover:underline">Daftar</Link>
        </p>
      </div>
    </div>
  )
}