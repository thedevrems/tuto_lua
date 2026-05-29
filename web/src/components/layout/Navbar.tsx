import { Link, NavLink } from 'react-router-dom'
import { useAuth } from '../../auth/AuthContext'

const linkClass = ({ isActive }: { isActive: boolean }) =>
  'text-sm font-medium transition-colors ' + (isActive ? 'text-black' : 'text-gray-600 hover:text-black')

/** Top navigation shared by all marketing pages. */
export default function Navbar() {
  const { user, logout } = useAuth()

  return (
    <header className="sticky top-0 z-30 border-b border-gray-200 bg-white/90 backdrop-blur">
      <div className="container-page flex h-16 items-center gap-8">
        <Link to="/" className="flex items-center gap-2.5">
          <span className="grid h-8 w-8 place-items-center rounded-md bg-black font-mono text-sm font-bold text-white">
            L
          </span>
          <span className="text-[15px] font-bold tracking-tight text-black">Lua Academy</span>
        </Link>

        <nav className="hidden items-center gap-6 md:flex">
          <NavLink to="/" className={linkClass} end>
            Accueil
          </NavLink>
          <NavLink to="/pricing" className={linkClass}>
            Tarifs
          </NavLink>
          <NavLink to="/learn" className={linkClass}>
            Cours
          </NavLink>
        </nav>

        <div className="ml-auto flex items-center gap-3">
          {user ? (
            <>
              <span className="hidden text-sm text-gray-600 sm:inline">
                Bonjour, <span className="font-medium text-black">{user.username}</span>
              </span>
              {user.role === 'admin' && (
                <Link to="/admin" className="btn btn-ghost btn-sm">
                  Admin
                </Link>
              )}
              <Link to="/learn" className="btn btn-primary btn-sm">
                Continuer
              </Link>
              <button onClick={logout} className="btn btn-ghost btn-sm">
                Déconnexion
              </button>
            </>
          ) : (
            <>
              <Link to="/login" className="btn btn-ghost btn-sm">
                Connexion
              </Link>
              <Link to="/register" className="btn btn-primary btn-sm">
                Créer un compte
              </Link>
            </>
          )}
        </div>
      </div>
    </header>
  )
}
