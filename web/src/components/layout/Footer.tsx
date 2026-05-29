import { Link } from 'react-router-dom'

/** Site footer shared by all marketing pages. */
export default function Footer() {
  return (
    <footer className="border-t border-gray-200 bg-white-soft">
      <div className="container-page flex flex-col items-center justify-between gap-4 py-10 sm:flex-row">
        <div className="flex items-center gap-2.5">
          <span className="grid h-7 w-7 place-items-center rounded-md bg-black font-mono text-xs font-bold text-white">
            L
          </span>
          <span className="text-sm font-semibold text-black">Lua Academy</span>
        </div>
        <nav className="flex items-center gap-6 text-sm text-gray-600">
          <Link to="/" className="hover:text-black">
            Accueil
          </Link>
          <Link to="/pricing" className="hover:text-black">
            Tarifs
          </Link>
          <Link to="/learn" className="hover:text-black">
            Cours
          </Link>
        </nav>
        <p className="text-xs text-gray-500">© {new Date().getFullYear()} Lua Academy. Tous droits réservés.</p>
      </div>
    </footer>
  )
}
