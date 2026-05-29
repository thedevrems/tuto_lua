import type { ReactNode } from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { useAuth } from '../auth/AuthContext'

interface Props {
  children: ReactNode
  /** When set, the user must hold this role (e.g. "admin"). */
  role?: 'admin'
}

/** Guards a route: redirects to /login when unauthenticated, home when the
 *  role requirement is not met. */
export default function ProtectedRoute({ children, role }: Props) {
  const { user, loading } = useAuth()
  const location = useLocation()

  if (loading) {
    return <div className="grid min-h-screen place-items-center text-gray-500">Chargement…</div>
  }
  if (!user) {
    return <Navigate to="/login" state={{ from: location.pathname }} replace />
  }
  if (role && user.role !== role) {
    return <Navigate to="/" replace />
  }
  return <>{children}</>
}
