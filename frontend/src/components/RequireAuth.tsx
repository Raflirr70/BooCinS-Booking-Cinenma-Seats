import type { ReactNode } from 'react'
import { Navigate } from 'react-router-dom'
import { useAuth, homePathFor } from '../store/auth'

export default function RequireAuth({ roles, children }: { roles?: string[]; children: ReactNode }) {
  const token = useAuth((s) => s.token)
  const user = useAuth((s) => s.user)
  if (!token || !user) return <Navigate to="/login" replace />
  if (roles && !roles.includes(user.role.name)) {
    return <Navigate to={homePathFor(user.role.name)} replace />
  }
  return <>{children}</>
}