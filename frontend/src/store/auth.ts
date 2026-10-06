import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export interface AuthUser {
  id: number
  first_name: string
  last_name: string
  email: string
  role: { id: number; name: string }
}
interface AuthState {
  token: string
  user: AuthUser | null
  setAuth: (token: string, user: AuthUser) => void
  clearAuth: () => void
}

export const useAuth = create<AuthState>()(
  persist(
    (set) => ({
      token: '',
      user: null,
      setAuth: (token, user) => set({ token, user }),
      clearAuth: () => set({ token: '', user: null }),
    }),
    { name: 'boocins-auth' },
  ),
)

export function homePathFor(role: string): string {
  if (role === 'super_admin') return '/superadmin'
  if (role === 'admin') return '/admin'
  if (role === 'manager') return '/manager'
  if (role === 'staff') return '/staff'
  return '/' // member dan guest ke halaman utama
}