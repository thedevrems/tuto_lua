import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../auth/AuthContext'
import { ApiError } from '../lib/api'
import AuthShell from '../components/layout/AuthShell'

export default function RegisterPage() {
  const { register } = useAuth()
  const navigate = useNavigate()

  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirm, setConfirm] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)
    if (password !== confirm) {
      setError('Les mots de passe ne correspondent pas.')
      return
    }
    setSubmitting(true)
    try {
      await register(username, email, password)
      navigate('/learn', { replace: true })
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Inscription impossible. Réessayez.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <AuthShell title="Créer un compte" subtitle="Commencez à apprendre le Lua dès aujourd'hui.">
      <form onSubmit={handleSubmit} className="space-y-5">
        {error && (
          <div className="rounded-md border border-danger-border bg-danger-bg px-4 py-3 text-sm text-danger">
            {error}
          </div>
        )}
        <div>
          <label htmlFor="username" className="label">
            Nom d'utilisateur
          </label>
          <input
            id="username"
            className="input"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            autoComplete="username"
            required
          />
        </div>
        <div>
          <label htmlFor="email" className="label">
            Adresse e-mail
          </label>
          <input
            id="email"
            type="email"
            className="input"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            autoComplete="email"
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
            autoComplete="new-password"
            required
          />
          <p className="mt-1.5 text-xs text-gray-500">
            8 caractères minimum, avec une majuscule, une minuscule et un chiffre.
          </p>
        </div>
        <div>
          <label htmlFor="confirm" className="label">
            Confirmer le mot de passe
          </label>
          <input
            id="confirm"
            type="password"
            className="input"
            value={confirm}
            onChange={(e) => setConfirm(e.target.value)}
            autoComplete="new-password"
            required
          />
        </div>
        <button type="submit" className="btn btn-primary w-full" disabled={submitting}>
          {submitting ? 'Création…' : 'Créer mon compte'}
        </button>
      </form>

      <p className="mt-6 text-center text-sm text-gray-600">
        Déjà un compte ?{' '}
        <Link to="/login" className="font-medium text-black underline underline-offset-2">
          Se connecter
        </Link>
      </p>
    </AuthShell>
  )
}
