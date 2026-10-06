import { Outlet, useNavigate } from 'react-router-dom'
import { useAuth } from '../store/auth'

export default function DashboardLayout() {
  const user = useAuth((s) => s.user)
  const clearAuth = useAuth((s) => s.clearAuth)
  const navigate = useNavigate()

  const logout = () => {
    clearAuth()
    navigate('/login')
  }

  return (
    <div className="min-h-screen bg-gray-900 text-white flex flex-col">
      <header className="border-b border-gray-800 bg-gray-900 px-6 py-4 flex items-center justify-between">
        <h1 className="text-xl font-bold text-red-500">BooCinS Panel</h1>
        <div className="flex items-center gap-4">
          <span className="text-sm text-gray-300">
            {user?.first_name} ({user?.role.name})
          </span>
          <button
            onClick={logout}
            className="rounded bg-gray-800 px-3 py-1.5 text-xs text-red-400 hover:bg-gray-700"
          >
            Logout
          </button>
        </div>
      </header>
      <main className="flex-1 p-6">
        <Outlet />
      </main>
    </div>
  )
}