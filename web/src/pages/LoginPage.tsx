import { useState, type FormEvent } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { useAuth } from '../auth/AuthContext'
import { ApiError } from '../lib/api'
import AuthShell from '../components/layout/AuthShell'

export default function LoginPage() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const location = useLocation() as { state?: { from?: string } }

  const [identifier, setIdentifier] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)
    setSubmitting(true)
    try {
      await login(identifier, password)
      navigate(location.state?.from ?? '/learn', { replace: true })
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Connexion impossible. Réessayez.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <AuthShell title="Connexion" subtitle="Heureux de vous revoir sur Lua Academy.">
      <form onSubmit={handleSubmit} className="space-y-5">
        {error && (
          <div className="rounded-md border border-danger-border bg-danger-bg px-4 py-3 text-sm text-danger">
            {error}
          </div>
        )}
        <div>
          <label htmlFor="identifier" className="label">
            E-mail ou nom d'utilisateur
          </label>
          <input
            id="identifier"
            className="input"
            value={identifier}
            onChange={(e) => setIdentifier(e.target.value)}
            autoComplete="username"
            required
          />
        </div>
        <div>
          <label htmlFor="password" className="label">
            Mot de passe
          </label>
          <input
            id="password"
            type="password"
            className="input"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            autoComplete="current-password"
            required
          />
        </div>
        <button type="submit" className="btn btn-primary w-full" disabled={submitting}>
          {submitting ? 'Connexion…' : 'Se connecter'}
        </button>
      </form>

      <p className="mt-6 text-center text-sm text-gray-600">
        Pas encore de compte ?{' '}
        <Link to="/register" className="font-medium text-black underline underline-offset-2">
          Créer un compte
        </Link>
      </p>
    </AuthShell>
  )
}
