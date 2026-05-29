import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api, ApiError } from '../lib/api'
import { useAuth } from '../auth/AuthContext'
import type { LockedCourse } from '../content/useCurriculum'

function formatPrice(cents: number, currency: string): string {
  return (cents / 100).toLocaleString('fr-FR', { style: 'currency', currency: currency.toUpperCase() })
}

/** Shown in the workspace when the user selects a course they haven't unlocked.
 *  The button starts a real Stripe Checkout session and redirects to it. */
export default function Paywall({ course }: { course: LockedCourse }) {
  const { user } = useAuth()
  const navigate = useNavigate()
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const buy = async () => {
    if (!user) {
      navigate('/login', { state: { from: '/learn' } })
      return
    }
    setBusy(true)
    setError(null)
    try {
      const { url } = await api.payments.checkout(course.id)
      window.location.href = url // redirect to Stripe Checkout
    } catch (e) {
      setError(e instanceof ApiError ? e.message : 'Paiement indisponible pour le moment.')
      setBusy(false)
    }
  }

  return (
    <div className="grid h-full place-items-center px-6 text-center">
      <div className="card max-w-md">
        <span className="badge mb-4">Cours verrouillé</span>
        <h2 className="text-2xl font-bold text-black">{course.title}</h2>
        {course.summary && <p className="mt-2 text-gray-600">{course.summary}</p>}
        <div className="mt-6 text-3xl font-black tracking-tight text-black">
          {formatPrice(course.priceCents, course.currency)}
        </div>
        {error && (
          <div className="mt-4 rounded-md border border-danger-border bg-danger-bg px-4 py-2 text-sm text-danger">
            {error}
          </div>
        )}
        <button onClick={buy} disabled={busy} className="btn btn-primary mt-6 w-full">
          {busy ? 'Redirection…' : user ? 'Acheter ce cours' : 'Se connecter pour acheter'}
        </button>
        <p className="mt-3 text-xs text-gray-500">Paiement sécurisé par Stripe. Accès immédiat après l'achat.</p>
      </div>
    </div>
  )
}
