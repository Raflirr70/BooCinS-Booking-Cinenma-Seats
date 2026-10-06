import { Link, useLocation } from 'react-router-dom'

const menus = [
  { label: 'Dashboard', path: '/admin' },
  { label: 'Kelola Film', path: '/admin/films' },
  { label: 'Kelola Room', path: '/admin/rooms' },
  { label: 'Kelola Schedule', path: '/admin/schedules' },
  { label: 'Kelola Promo', path: '/admin/promos' },
]

export default function Sidebar() {
  const location = useLocation()

  return (
    <aside className="w-64 border-r border-gray-800 bg-gray-900 p-4 min-h-screen">
      <h2 className="text-lg font-bold text-white mb-6 px-3">Admin Menu</h2>
      <nav className="space-y-1">
        {menus.map((item) => {
          const active = location.pathname === item.path
          return (
            <Link
              key={item.path}
              to={item.path}
              className={`block px-3 py-2 rounded-lg text-sm font-medium ${
                active ? 'bg-red-600 text-white' : 'text-gray-400 hover:bg-gray-800 hover:text-white'
              }`}
            >
              {item.label}
            </Link>
          )
        })}
      </nav>
    </aside>
  )
}